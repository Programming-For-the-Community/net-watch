CREATE TABLE IF NOT EXISTS check_rollups_hourly (
    check_id       bigint      NOT NULL REFERENCES checks (id)   ON DELETE CASCADE,
    network_id     bigint      NOT NULL REFERENCES networks (id) ON DELETE CASCADE,
    bucket_start   timestamptz NOT NULL,
    total_count    int         NOT NULL,
    success_count  int         NOT NULL,
    latency_min    real,
    latency_avg    real,
    latency_max    real,
    latency_p95    real,

    PRIMARY KEY (check_id, network_id, bucket_start),
    CONSTRAINT check_rollups_counts_ck
        CHECK (success_count >= 0 AND success_count <= total_count)
);

CREATE INDEX IF NOT EXISTS check_rollups_network_bucket_idx
    ON check_rollups_hourly (network_id, bucket_start);