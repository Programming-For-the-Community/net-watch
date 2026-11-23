CREATE TABLE IF NOT EXISTS audit_log (
    id           bigint      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    occurred_at  timestamptz NOT NULL DEFAULT now(),
    actor        text        NOT NULL,
    action       text        NOT NULL,
    entity_type  text        NOT NULL,
    entity_id    bigint,          -- no FK, so history survives deletes
    before       jsonb,
    after        jsonb
);

CREATE INDEX IF NOT EXISTS audit_log_occurred_idx
    ON audit_log (occurred_at DESC);

CREATE INDEX IF NOT EXISTS audit_log_entity_idx
    ON audit_log (entity_type, entity_id, occurred_at DESC);