CREATE EXTENSION IF NOT EXISTS vector; -- if not already enabled

CREATE TABLE faces (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    image_id   UUID NOT NULL REFERENCES photos(id) ON DELETE CASCADE,
    upload_id   UUID NOT NULL REFERENCES uploads(id) ON DELETE CASCADE,
    embedding  VECTOR(512) NOT NULL,
    bbox       JSONB NOT NULL,
    cluster_id INTEGER
);

CREATE TRIGGER set_updated_at_faces
    BEFORE UPDATE ON faces
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

CREATE INDEX ON faces (upload_id);
CREATE INDEX ON faces (image_id);

ALTER TABLE photos
ADD COLUMN is_embedded BOOLEAN NOT NULL DEFAULT FALSE;