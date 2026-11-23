CREATE TABLE IF NOT EXISTS network_sessions (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    network_id BIGINT NOT NULL REFERENCES networks(id) ON DELETE CASCADE,
    connected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    disconnected_at TIMESTAMPTZ,
    local_ip INET,
    interface_name TEXT,
    connnection_type TEXT,

    CONSTRAINT network_sessions_time_order_ck
        CHECK (disconnected_at IS NULL OR disconnected_at >= connected_at)
);

-- Only one open (current) session at a time
CREATE UNIQUE INDEX IF NOT EXISTS network_sessions_one_open_uq
    ON network_sessions ((true))
    WHERE disconnected_at IS NULL;

CREATE INDEX IF NOT EXISTS network_sessions_network_connected_idx
    ON network_sessions (network_id, connected_at DESC);