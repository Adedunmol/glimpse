import asyncio
import logging
import signal

from redis_queue.consumer import QueueConsumer
from processing.handler import handle_job

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
)

logger = logging.getLogger(__name__)


async def main():
    logger.info("Starting ML processing service")

    consumer = QueueConsumer()
    await consumer.connect()

    loop = asyncio.get_running_loop()
    stop_event = asyncio.Event()

    def _shutdown():
        logger.info("Shutdown signal received")
        stop_event.set()

    for sig in (signal.SIGINT, signal.SIGTERM):
        loop.add_signal_handler(sig, _shutdown)

    listen_task = asyncio.create_task(consumer.listen(handle_job))

    await stop_event.wait()

    listen_task.cancel()
    await asyncio.gather(listen_task, return_exceptions=True)
    await consumer.close()
    logger.info("Service stopped cleanly")


if __name__ == "__main__":
    asyncio.run(main())