CREATE INDEX IF NOT EXISTS idx_uploads_user_id ON uploads(host_id);--user_id -> host_id
CREATE INDEX IF NOT EXISTS idx_links_cluster_id ON links(cluster_id);