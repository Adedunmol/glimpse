from sqlalchemy import text
from db.connection import get_db
from collections import defaultdict
import json

async def save_face_embeddings(image_id: str, upload_id: str, embeddings: list):
    async with get_db() as conn:
        await conn.execute(
            text("""
                INSERT INTO faces (image_id, upload_id, embedding, bbox)
                VALUES (:image_id, :upload_id, (:embedding)::vector, :bbox)
            """),
            [
                {
                    "image_id": image_id,
                    "upload_id": upload_id,
                    "embedding": str(e["embedding"]),  # pgvector expects string format
                    "bbox": json.dumps(e["bbox"]),
                }
                for e in embeddings
            ],
        )


async def mark_image_processed(image_id: str):
    async with get_db() as conn:
        await conn.execute(
            text("UPDATE photos SET is_embedded = TRUE WHERE id = :image_id"),
            {"image_id": image_id},
        )


async def fetch_embeddings_for_upload(upload_id: str):
    async with get_db() as conn:
        result = await conn.execute(
            text("SELECT id, image_id, embedding FROM faces WHERE upload_id = :upload_id"),
            {"upload_id": upload_id},
        )
        return result.mappings().all()  # returns list of dict-like rows


async def save_cluster_labels(upload_id: str, face_ids: list, image_ids: list, labels: list):
    async with get_db() as conn:
        await conn.execute(
            text("UPDATE faces SET cluster_id = :cluster_id WHERE id = :face_id"),
            [{"cluster_id": int(l), "face_id": fid} for fid, l in zip(face_ids, labels)],
        )

        groups: dict[int, set[str]] = defaultdict(set)
        for image_id, label in zip(image_ids, labels):
            if label == -1:
                continue
            groups[label].add(image_id)

        for label, photo_ids in groups.items():
            cluster_id = await conn.scalar(
                text("INSERT INTO clusters (upload_id) VALUES (:upload_id) RETURNING id"),
                {"upload_id": upload_id},
            )
            await conn.execute(
                text("""
                    INSERT INTO cluster_photos (cluster_id, photo_id)
                    VALUES (:cluster_id, :photo_id)
                    ON CONFLICT DO NOTHING
                """),
                [{"cluster_id": cluster_id, "photo_id": pid} for pid in photo_ids],
            )