-- Represents a user
CREATE TABLE users (
    clerk_user_id TEXT PRIMARY KEY,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER set_updated_at_users
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- Represents a device associated with a user for push notifications. A user can have multiple devices (e.g. phone and tablet).
CREATE TABLE user_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    user_id TEXT NOT NULL,

    push_token TEXT NOT NULL UNIQUE,
    platform TEXT NOT NULL
);

CREATE TRIGGER set_updated_at_user_devices
    BEFORE UPDATE ON user_devices
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- Represents one upload session created by a host.
CREATE TABLE uploads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    name TEXT NOT NULL,
    host_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending'
           CHECK (status IN ('pending', 'processing', 'done', 'failed')),
    expires_at TIMESTAMPTZ NOT NULL DEFAULT (now() + INTERVAL '30 days')
);

CREATE TRIGGER set_updated_at_uploads
    BEFORE UPDATE ON uploads
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- One row per uploaded photo.
CREATE TABLE photos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    upload_id    UUID NOT NULL REFERENCES uploads(id) ON DELETE CASCADE,
    storage_key  TEXT NOT NULL UNIQUE,
    status       TEXT NOT NULL DEFAULT 'pending'
                 CHECK (status IN ('pending', 'uploaded'))
);

CREATE TRIGGER set_updated_at_photos
    BEFORE UPDATE ON photos
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- One row per unique face identity detected across an upload session.
CREATE TABLE clusters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        
    upload_id          UUID NOT NULL REFERENCES uploads(id) ON DELETE CASCADE,
    label              TEXT,
    thumbnail_photo_id UUID REFERENCES photos(id) ON DELETE SET NULL
);

CREATE TRIGGER set_updated_at_clusters
    BEFORE UPDATE ON clusters
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- Join table to represent the many-to-many relationship between clusters and photos.
CREATE TABLE cluster_photos (
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    cluster_id  UUID NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
    photo_id    UUID NOT NULL REFERENCES photos(id) ON DELETE CASCADE,
    PRIMARY KEY (cluster_id, photo_id)
);

CREATE TRIGGER set_updated_at_cluster_photos
    BEFORE UPDATE ON cluster_photos
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- Represents a shareable link to a cluster, which can optionally be password-protected and have an expiration date.
CREATE TABLE links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    cluster_id UUID NOT NULL REFERENCES clusters(id) ON DELETE CASCADE, -- the cluster this link points to
    token TEXT NOT NULL UNIQUE, -- random token included in the shareable URL for lookup
    is_password_protected BOOLEAN NOT NULL DEFAULT false,
    password_hash TEXT,
    expires_at TIMESTAMPTZ NOT NULL DEFAULT (now() + INTERVAL '30 days'),
    is_active BOOLEAN NOT NULL DEFAULT true,

    CONSTRAINT password_hash_required
        CHECK (is_password_protected = false OR password_hash IS NOT NULL)
);

CREATE TRIGGER set_updated_at_links
    BEFORE UPDATE ON links
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

-- Frequent lookup: all photos for an upload
CREATE INDEX idx_photos_upload_id ON photos(upload_id);

-- Frequent lookup: all clusters for an upload
CREATE INDEX idx_clusters_upload_id ON clusters(upload_id);

-- Frequent lookup: all photos in a cluster (recipient view)
CREATE INDEX idx_cluster_photos_cluster_id ON cluster_photos(cluster_id);

-- Token lookup on every share link access (must be fast)
CREATE INDEX idx_links_token ON links(token);

-- Expiry cleanup job
CREATE INDEX idx_links_expires_at ON links(expires_at) WHERE is_active = true;
CREATE INDEX idx_uploads_expires_at ON uploads(expires_at);