CREATE TABLE IF NOT EXISTS incidents (
    id               bigint      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    check_id         bigint      NOT NULL REFERENCES checks (id)   ON DELETE CASCADE,
    network_id       bigint      NOT NULL REFERENCES networks (id) ON DELETE CASCADE,
    started_at       timestamptz NOT NULL,
    resolved_at      timestamptz,
    failure_count    int         NOT NULL DEFAULT 0,
    first_error      text,
    last_error       text,
    acknowledged_at  timestamptz,
    notes            text,

    CONSTRAINT incidents_time_order_ck
        CHECK (resolved_at IS NULL OR resolved_at >= started_at)
);

-- One open incident per check
CREATE UNIQUE INDEX IF NOT EXISTS incidents_one_open_per_check_uq
    ON incidents (check_id)
    WHERE resolved_at IS NULL;

CREATE INDEX IF NOT EXISTS incidents_started_idx
    ON incidents (started_at DESC);

CREATE INDEX IF NOT EXISTS incidents_network_started_idx
    ON incidents (network_id, started_at DESC);