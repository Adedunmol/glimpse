import logging

logger = logging.getLogger(__name__)


async def handle_job(job: dict):
    """
    Entry point for processing a single job payload.
    Expected shape (adjust to match your producer):
    {
        "event_id": "...",
        "image_key": "s3://bucket/path.jpg",
        "type": "face_detect" | "embedding" | "cluster"
    }
    """
    logger.info("Processing job: %s", job)

    job_type = job.get("type")

    if job_type == "face_detect":
        await _handle_face_detect(job)
    elif job_type == "embedding":
        await _handle_embedding(job)
    elif job_type == "cluster":
        await _handle_cluster(job)
    else:
        logger.warning("Unknown job type: %s", job_type)


async def _handle_face_detect(job: dict):
    # TODO: download image, run face detection
    pass


async def _handle_embedding(job: dict):
    # TODO: generate face/body embeddings
    pass


async def _handle_cluster(job: dict):
    # TODO: trigger clustering for an event
    pass