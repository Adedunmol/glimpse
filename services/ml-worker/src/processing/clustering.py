import logging

logger = logging.getLogger(__name__)

async def cluster_event(event_id: str):
    logger.info(f'clustering embeddings for {event_id}')
    pass