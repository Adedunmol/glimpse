from sqlalchemy import text
from db.connection import get_db


async def save_face_embeddings(image_id: str, event_id: str, embeddings: list):
    async with get_db() as conn:
        await conn.execute(
            text("""
                INSERT INTO faces (image_id, event_id, embedding, bbox)
                VALUES (:image_id, :event_id, :embedding::vector, :bbox)
            """),
            [
                {
                    "image_id": image_id,
                    "event_id": event_id,
                    "embedding": str(e["embedding"]),  # pgvector expects string format
                    "bbox": e["bbox"],
                }
                for e in embeddings
            ],
        )


async def mark_image_processed(image_id: str):
    async with get_db() as conn:
        await conn.execute(
            text("UPDATE images SET is_embedded = TRUE WHERE id = :image_id"),
            {"image_id": image_id},
        )


async def fetch_embeddings_for_event(event_id: str):
    async with get_db() as conn:
        result = await conn.execute(
            text("SELECT id, embedding FROM faces WHERE event_id = :event_id"),
            {"event_id": event_id},
        )
        return result.mappings().all()  # returns list of dict-like rows


async def save_cluster_labels(event_id: str, face_ids: list, labels: list):
    async with get_db() as conn:
        await conn.execute(
            text("""
                UPDATE faces SET cluster_id = :cluster_id
                WHERE id = :face_id
            """),
            [
                {"cluster_id": int(label), "face_id": face_id}
                for face_id, label in zip(face_ids, labels)
            ],
        )