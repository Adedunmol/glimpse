import logging

logger = logging.getLogger(__name__)

async def process_image(image_id: str, event_id: str, s3_key: str, redis_client):
    logger.info("processing image")
    pass