import logging
from processing.embedding import process_image
from processing.clustering import cluster_upload


logger = logging.getLogger(__name__)


async def handle_job(data: dict, redis_client):
    logger.info("Processing job: %s", data)

    job_type = data.get("type")

    if job_type == "process_image":
        await _handle_process_image(data, redis_client)
    elif job_type == "cluster_event":
        await _handle_cluster(data)
    else:
        logger.warning("Unknown job type: %s", job_type)


async def _handle_process_image(data: dict, redis_client):
    await process_image(
        image_id=data["image_id"],
        event_id=data["upload_id"],
        s3_key=data["s3_key"],
        redis_client=redis_client,
    )


async def _handle_cluster(data: dict):
    await cluster_upload(upload_id=data["upload_id"])