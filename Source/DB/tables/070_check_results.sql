CREATE TABLE IF NOT EXISTS check_results (
    id           bigint      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    check_id     bigint      NOT NULL REFERENCES checks (id)   ON DELETE CASCADE,
    network_id   bigint      NOT NULL REFERENCES networks (id) ON DELETE CASCADE,
    checked_at   timestamptz NOT NULL,
    success      boolean     NOT NULL,
    latency_ms   real,
    status_code  int,
    error        text,
    details      jsonb
);

CREATE INDEX IF NOT EXISTS check_results_check_checked_idx
    ON check_results (check_id, checked_at DESC);

CREATE INDEX IF NOT EXISTS check_results_network_checked_idx
    ON check_results (network_id, checked_at);

-- Cheap index for time-range scans and retention deletes
CREATE INDEX IF NOT EXISTS check_results_checked_brin
    ON check_results USING brin (checked_at);