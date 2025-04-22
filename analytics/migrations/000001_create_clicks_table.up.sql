CREATE TABLE IF NOT EXISTS xlink.clicks
(
    link_owner UUID,
    short_link String,
    clicked_at DateTime,
    referrer String,
    region LowCardinality(String),
    browser LowCardinality(String),
    device_type LowCardinality(String),
    os LowCardinality(String),
    ip_address IPv4,
    is_unique UInt8
) ENGINE = MergeTree()
ORDER BY (short_link, clicked_at)
TTL clicked_at + INTERVAL 1 MONTH DELETE
SETTINGS index_granularity = 8192;
