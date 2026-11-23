CREATE TABLE IF NOT EXISTS device_kinds (
    code TEXT PRIMARY KEY,
    display_name text NOT NULL,
    icon TEXT,
    sort_order INT NOT NULL DEFAULT 0,
    is_builtin BOOLEAN NOT NULL DEFAULT FALSE
);

-- Seed table with defaults -> Do nothing if row already exists
INSERT INTO device_kinds (code, display_name, icon, sort_order, is_builtin) VALUES
    ('router',  'Router',           'router',  10, true),
    ('server',  'Server',           'server',  20, true),
    ('nas',     'NAS',              'storage', 30, true),
    ('sbc',     'Single-board PC',  'chip',    40, true),
    ('iot',     'Smart device',     'plug',    50, true),
    ('service', 'Service',          'globe',   60, true),
    ('other',   'Other',            'device',  99, true)
ON CONFLICT (code) DO NOTHING;