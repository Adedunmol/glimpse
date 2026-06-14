import logging
from config import STREAM_NAME

logger = logging.getLogger(__name__)

async def check_and_trigger_clustering(client, event_id: str):
    logger.info("check and confirm whether to start clustering")
    
    total_key = f"event:{event_id}:total"
    processed_key = f"event:{event_id}:processed"
    lock_key = f"event:{event_id}:cluster_triggered"

    processed = await client.incr(processed_key)
    total = await client.get(total_key)

    if total is None:
        logger.warning("No total count found for event %s", event_id)
        return

    logger.info("Event %s: %s/%s images processed", event_id, processed, total)

    if processed >= int(total):
        acquired = await client.set(lock_key, 1, nx=True, ex=60)
        if acquired:
            await client.xadd(STREAM_NAME, {"type": "cluster_event", "event_id": event_id})
            logger.info("Pushed cluster_event for event %s", event_id)
            await client.delete(total_key, processed_key)