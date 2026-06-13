import asyncio
import json
import logging

import redis.asyncio as redis

from config import REDIS_URL, QUEUE_NAME

logger = logging.getLogger(__name__)


class QueueConsumer:
    def __init__(self, redis_url: str = REDIS_URL, queue_name: str = QUEUE_NAME):
        self.redis_url = redis_url
        self.queue_name = queue_name
        self._client: redis.Redis | None = None

    async def connect(self):
        self._client = redis.from_url(self.redis_url, decode_responses=True)
        await self._client.ping()
        logger.info("Connected to Redis at %s", self.redis_url)

    async def listen(self, handler):
        """Blocking pop loop. handler is an async callable that receives the parsed job dict."""
        if self._client is None:
            await self.connect()

        logger.info("Listening on queue: %s", self.queue_name)
        while True:
            try:
                # BLPOP blocks until an item is available or timeout elapses
                result = await self._client.blpop(self.queue_name, timeout=5)
                if result is None:
                    continue  # timeout, loop again

                _, raw_job = result
                try:
                    job = json.loads(raw_job)
                except json.JSONDecodeError:
                    logger.error("Failed to decode job payload: %s", raw_job)
                    continue

                await handler(job)

            except asyncio.CancelledError:
                logger.info("Consumer task cancelled, shutting down")
                break
            except Exception:
                logger.exception("Error while processing queue item")
                await asyncio.sleep(1)  # backoff before retrying

    async def close(self):
        if self._client:
            await self._client.close()