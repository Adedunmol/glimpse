import logging
import json
import asyncio
import numpy as np
from sklearn.cluster import DBSCAN
from db.repository import fetch_embeddings_for_upload, save_cluster_labels
from executor import cpu_executor

logger = logging.getLogger(__name__)

async def cluster_upload(upload_id: str):
    logger.info(f'clustering embeddings for {upload_id}')
    rows = await fetch_embeddings_for_upload(upload_id)

    if not rows:
        logger.info("No embeddings for upload %s — skipping clustering", upload_id)
        return

    face_ids   = [row["id"] for row in rows]
    image_ids  = [row["image_id"] for row in rows]
    embeddings = np.array([json.loads(row["embedding"]) for row in rows])

    db = DBSCAN(eps=0.4, min_samples=1, metric="cosine")

    loop = asyncio.get_running_loop()
    labels = await loop.run_in_executor(cpu_executor, db.fit_predict, embeddings)

    cluster_count = len(set(labels)) - (1 if -1 in labels else 0)
    logger.info(
        "Upload %s: %d faces → %d clusters (%d outliers)",
        upload_id,
        len(face_ids),
        cluster_count,
        list(labels).count(-1),
    )

    await save_cluster_labels(upload_id, face_ids, image_ids, labels.tolist())