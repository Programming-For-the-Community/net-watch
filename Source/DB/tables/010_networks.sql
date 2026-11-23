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