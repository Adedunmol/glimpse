import logging
import io
import asyncio
import numpy as np
from sqlalchemy import text
from db.connection import get_db
from insightface.app import FaceAnalysis
from storage.s3 import download_image
from processing.trigger import check_and_trigger_clustering
from db.repository import save_face_embeddings, mark_image_processed
from PIL import Image
from .executor import cpu_executor

logger = logging.getLogger(__name__)

_app = FaceAnalysis(name="buffalo_l", providers=["CPUExecutionProvider"])
_app.prepare(ctx_id=0, det_size=(640, 640))

async def process_image(image_id: str, event_id: str, s3_key: str, redis_client):
    logger.info("processing image")

    async with get_db() as conn:
        already_processed = await conn.scalar(
            text("SELECT is_embedded FROM photos WHERE id = :image_id"),
            {"image_id": image_id},
        )

    if already_processed:
        logger.info("Image %s already embedded, skipping", image_id)
        return

    raw = await download_image(s3_key)

    # Convert bytes → numpy array → BGR image for InsightFace
    # image_array = np.frombuffer(raw, dtype=np.uint8)
    # image = cv2.imdecode(image_array, cv2.IMREAD_COLOR)

    image = Image.open(io.BytesIO(raw)).convert("RGB")
    image_array = np.array(image)[:, :, ::-1]  # RGB -> BGR for InsightFace


    loop = asyncio.get_running_loop()
    faces = await loop.run_in_executor(cpu_executor, _app.get, image_array)

    if not faces:
        logger.info("No faces detected in image %s", image_id)
    else:
        logger.info("Detected %d face(s) in image %s", len(faces), image_id)
        embeddings = [
            {
                "embedding": face.embedding.tolist(),  # 512-dim vector for pgvector
                "bbox": face.bbox.tolist(),            # [x1, y1, x2, y2]
            }
            for face in faces
        ]
        await save_face_embeddings(image_id, event_id, embeddings)


    await mark_image_processed(image_id)

    await check_and_trigger_clustering(redis_client, event_id)
