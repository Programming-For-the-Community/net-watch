CREATE TABLE IF NOT EXISTS devices (
    id            bigint      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    network_id    bigint      REFERENCES networks (id) ON DELETE CASCADE,  -- NULL = global
    kind          text        NOT NULL DEFAULT 'other' REFERENCES device_kinds (code) ON UPDATE CASCADE,
    name          text        NOT NULL,
    address       text        NOT NULL,
    mac_address   macaddr,
    system_role   text,
    notes         text,
    paused_until  timestamptz,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT devices_network_name_uq UNIQUE NULLS NOT DISTINCT (network_id, name)
);

-- One gateway per network; one global internet/dns baseline
CREATE UNIQUE INDEX IF NOT EXISTS devices_network_system_role_uq
    ON devices (network_id, system_role) NULLS NOT DISTINCT
    WHERE system_role IS NOT NULL;

CREATE INDEX IF NOT EXISTS devices_mac_address_idx
    ON devices (mac_address);

-- Keep updated_at current (function defined in functions/000_set_updated_at.sql)
CREATE OR REPLACE TRIGGER devices_set_updated_at
    BEFORE UPDATE ON devices
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
