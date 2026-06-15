import logging
import numpy as np
from sklearn.cluster import DBSCAN
from db.repository import fetch_embeddings_for_event, save_cluster_labels


logger = logging.getLogger(__name__)

async def cluster_event(event_id: str):
    logger.info(f'clustering embeddings for {event_id}')
    rows = await fetch_embeddings_for_event(event_id)

    if not rows:
        logger.info("No embeddings for event %s — skipping clustering", event_id)
        return

    face_ids  = [row["id"] for row in rows]
    embeddings = np.array([row["embedding"] for row in rows])

    # DBSCAN groups embeddings by cosine distance.
    # eps=0.4 is a reasonable starting threshold for InsightFace buffalo_l
    # embeddings — faces within 0.4 cosine distance are the same person.
    # min_samples=1 means a single face still gets its own cluster rather
    # than being marked as noise (-1).
    db = DBSCAN(eps=0.4, min_samples=1, metric="cosine")
    labels = db.fit_predict(embeddings)

    cluster_count = len(set(labels)) - (1 if -1 in labels else 0)
    logger.info(
        "Event %s: %d faces → %d clusters (%d outliers)",
        event_id,
        len(face_ids),
        cluster_count,
        list(labels).count(-1),
    )

    await save_cluster_labels(event_id, face_ids, labels.tolist())