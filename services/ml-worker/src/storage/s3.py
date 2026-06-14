import logging

logger = logging.getLogger(__name__)

async def download_image(s3_key: str):
    logger.info(f'downloading image for  {s3_key}')
    pass