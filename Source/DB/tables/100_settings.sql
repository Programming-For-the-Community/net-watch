CREATE TABLE IF NOT EXISTS settings (
    key         text        PRIMARY KEY,
    value       jsonb       NOT NULL,
    updated_at  timestamptz NOT NULL DEFAULT now()
);

-- Seed default settings
INSERT INTO settings (key, value) VALUES
    ('retention.raw_days',           '30'::jsonb),
    ('retention.rollup_days',        '365'::jsonb),
    ('defaults.interval_seconds',    '60'::jsonb),
    ('defaults.failure_threshold',   '3'::jsonb),
    ('defaults.recovery_threshold',  '2'::jsonb)
ON CONFLICT (key) DO NOTHING;

-- Keep updated_at current (function defined in functions/000_set_updated_at.sql)
CREATE OR REPLACE TRIGGER settings_set_updated_at
    BEFORE UPDATE ON settings
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
