CREATE TABLE IF NOT EXISTS check_status (
    check_id               bigint      PRIMARY KEY REFERENCES checks (id) ON DELETE CASCADE,
    state                  text        NOT NULL DEFAULT 'unknown',
    consecutive_failures   int         NOT NULL DEFAULT 0,
    consecutive_successes  int         NOT NULL DEFAULT 0,
    last_checked_at        timestamptz,
    last_latency_ms        real,
    last_error             text,
    last_network_id        bigint      REFERENCES networks (id)  ON DELETE SET NULL,
    state_changed_at       timestamptz NOT NULL DEFAULT now(),
    current_incident_id    bigint      REFERENCES incidents (id) ON DELETE SET NULL,
    updated_at             timestamptz NOT NULL DEFAULT now()
);