CREATE TABLE IF NOT EXISTS discovered_devices (
    id             bigint      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    network_id     bigint      NOT NULL REFERENCES networks (id) ON DELETE CASCADE,
    mac_address    macaddr     NOT NULL,
    ip_address     inet        NOT NULL,
    hostname       text,
    vendor         text,
    first_seen_at  timestamptz NOT NULL DEFAULT now(),
    last_seen_at   timestamptz NOT NULL DEFAULT now(),
    device_id      bigint      REFERENCES devices (id) ON DELETE SET NULL,
    ignored        boolean     NOT NULL DEFAULT false,

    CONSTRAINT discovered_devices_network_mac_uq UNIQUE (network_id, mac_address)
);