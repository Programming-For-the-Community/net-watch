CREATE TABLE IF NOT EXISTS networks (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    fingerprint TEXT NOT NULL,
    name TEXT NOT NULL,
    gateway_ip INET,
    gateway_mac MACADDR,
    subnet CIDR,
    ssid TEXT,
    dns_servers INET[],
    notes TEXT,
    first_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT networks_fingerprint_uq UNIQUE (fingerprint)
);

-- Keep updated_at current (function defined in functions/000_set_updated_at.sql)
CREATE OR REPLACE TRIGGER networks_set_updated_at
    BEFORE UPDATE ON networks
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
