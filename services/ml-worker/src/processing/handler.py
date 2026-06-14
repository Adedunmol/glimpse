import logging

logger = logging.getLogger(__name__)


async def handle_job(job: dict):
    """
    Entry point for processing a single job payload.
    Expected shape (adjust to match your producer):
    {
        "event_id": "...",
        "image_key": "s3://bucket/path.jpg",
        "type": "process_image" | "cluster_event"
        "callback_url": "http://app:8080/uploads/callback"
    }
    """
    logger.info("Processing job: %s", job)

    job_type = job.get("type")

    if job_type == "process_image":
        await _handle_process_image(job)
    elif job_type == "cluster_event":
        await _handle_cluster(job)
    else:
        logger.warning("Unknown job type: %s", job_type)


async def _handle_process_image(job: dict):
    # TODO: download image, run face detection
    pass


async def _handle_cluster(job: dict):
    # TODO: trigger clustering for an event
    pass