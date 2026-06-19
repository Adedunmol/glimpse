import logging
import aiobotocore.session
from config import AWS_ACCESS_KEY, AWS_REGION, AWS_SECRET_KEY, S3_BUCKET, ENDPOINT_URL

logger = logging.getLogger(__name__)

async def download_image(s3_key: str) -> bytes:
    logger.info(f'downloading image for: {s3_key}')
    logger.info(f'endpoint url: {ENDPOINT_URL}')

    session  = aiobotocore.session.get_session()
    async with session.create_client(
        "s3",
        region_name=AWS_REGION,
        aws_access_key_id=AWS_ACCESS_KEY,
        aws_secret_access_key=AWS_SECRET_KEY,
        endpoint_url=ENDPOINT_URL
    ) as client:
        response = await client.get_object(Bucket=S3_BUCKET, Key=s3_key)
        async with response["Body"] as stream:
            return await stream.read()