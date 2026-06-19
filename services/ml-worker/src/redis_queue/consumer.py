import socket
import asyncio
import logging
import redis.asyncio as redis
from config import REDIS_URL, STREAM_NAME

GROUP_NAME = "ml-workers"
CONSUMER_NAME = f"worker-{socket.gethostname()}"  # unique per instance when scaling

logger = logging.getLogger(__name__)


class QueueConsumer:
    def __init__(
        self,
        redis_url: str = REDIS_URL,
        stream_name: str = STREAM_NAME,
        group_name: str = GROUP_NAME,
        consumer_name: str = CONSUMER_NAME,
    ):
        self.redis_url = redis_url
        self.stream_name = stream_name
        self.group_name = group_name
        self.consumer_name = consumer_name
        self._client: redis.Redis | None = None

    async def connect(self):
        self._client = redis.from_url(
            self.redis_url,
            decode_responses=True,
            socket_timeout=10,
            socket_connect_timeout=5,
        )
        await self._client.ping()

        try:
            await self._client.xgroup_create(
                self.stream_name,
                self.group_name,
                id="$",
                mkstream=True,
            )
            logger.info("Created consumer group '%s'", self.group_name)
        except redis.ResponseError as e:
            if "BUSYGROUP" in str(e):
                logger.info("Consumer group '%s' already exists, continuing", self.group_name)
            else:
                raise

        logger.info("Connected to Redis at %s", self.redis_url)

    async def listen(self, handler):
        if self._client is None:
            await self.connect()

        logger.info(
            "Listening on stream '%s' as '%s' in group '%s'",
            self.stream_name,
            self.consumer_name,
            self.group_name,
        )

        # Start at "0" to reclaim any messages that were in-flight when we
        # last died. Once the PEL is drained we switch to ">" for new ones.
        pending_id = "0"

        while True:
            try:
                messages = await self._client.xreadgroup(
                    groupname=self.group_name,
                    consumername=self.consumer_name,
                    streams={self.stream_name: pending_id},
                    count=1,
                    block=5000,
                )

                if not messages:
                    if pending_id == "0":
                        logger.info("PEL drained, switching to new messages")
                        pending_id = ">"
                    continue

                for _, stream_messages in messages:
                    for message_id, data in stream_messages:
                        success = await self._process_message(message_id, data, handler)
                        if success:
                            await self._client.xack(
                                self.stream_name, self.group_name, message_id
                            )
                        else:
                            logger.warning(
                                "Message %s left in PEL for retry", message_id
                            )

            except asyncio.CancelledError:
                logger.info("Consumer task cancelled, shutting down")
                break
            except Exception:
                logger.exception("Unexpected error in listen loop")
                await asyncio.sleep(1)

    async def _process_message(self, message_id: str, data: dict, handler) -> bool:
        try:
            logger.info("Processing message %s: %s", message_id, data)
            await handler(data, self._client)
            return True
        except Exception:
            logger.exception(
                "Failed to process message %s — leaving in PEL", message_id
            )
            return False

    async def close(self):
        if self._client:
            await self._client.aclose()