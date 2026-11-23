CREATE TABLE IF NOT EXISTS checks (
    id                  bigint      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    device_id           bigint      NOT NULL REFERENCES devices (id) ON DELETE CASCADE,
    type                text        NOT NULL,
    name                text        NOT NULL,
    enabled             boolean     NOT NULL DEFAULT true,
    interval_seconds    int         NOT NULL DEFAULT 60,
    timeout_ms          int         NOT NULL DEFAULT 2000,
    failure_threshold   int         NOT NULL DEFAULT 3,
    recovery_threshold  int         NOT NULL DEFAULT 2,
    config              jsonb       NOT NULL DEFAULT '{}'::jsonb,
    paused_until        timestamptz,
    created_at          timestamptz NOT NULL DEFAULT now(),
    updated_at          timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT checks_device_name_uq          UNIQUE (device_id, name),
    CONSTRAINT checks_interval_ck             CHECK (interval_seconds >= 5),
    CONSTRAINT checks_timeout_ck              CHECK (timeout_ms > 0),
    CONSTRAINT checks_failure_threshold_ck    CHECK (failure_threshold >= 1),
    CONSTRAINT checks_recovery_threshold_ck   CHECK (recovery_threshold >= 1)
);

CREATE INDEX IF NOT EXISTS checks_device_idx
    ON checks (device_id);