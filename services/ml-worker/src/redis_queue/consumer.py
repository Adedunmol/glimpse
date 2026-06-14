import asyncio
import logging

import redis.asyncio as redis

from config import REDIS_URL, STREAM_NAME

logger = logging.getLogger(__name__)


class QueueConsumer:
    def __init__(
        self,
        redis_url: str = REDIS_URL,
        stream_name: str = STREAM_NAME,
    ):
        self.redis_url = redis_url
        self.stream_name = stream_name
        self._client: redis.Redis | None = None
        self._last_id = "$"  # only read new messages from startup

    async def connect(self):
        self._client = redis.from_url(
            self.redis_url,
            decode_responses=True,
            socket_timeout=10,
            socket_connect_timeout=5,
        )
        await self._client.ping()
        logger.info("Connected to Redis at %s", self.redis_url)

    async def listen(self, handler):
        if self._client is None:
            await self.connect()

        logger.info("Listening on stream '%s'", self.stream_name)

        while True:
            try:
                messages = await self._client.xread(
                    streams={self.stream_name: self._last_id},
                    count=1,
                    block=5000,
                )

                if not messages:
                    continue

                for _, stream_messages in messages:
                    for message_id, data in stream_messages:
                        self._last_id = message_id  # advance cursor
                        await self._process_message(message_id, data, handler)

            except asyncio.CancelledError:
                logger.info("Consumer task cancelled, shutting down")
                break
            except Exception:
                logger.exception("Unexpected error in listen loop")
                await asyncio.sleep(1)

    async def _process_message(self, message_id: str, data: dict, handler):
        try:
            logger.info("Received message %s: %s", message_id, data)
            await handler(data)
        except Exception:
            logger.exception("Failed to process message %s", message_id)

    async def close(self):
        if self._client:
            await self._client.aclose()