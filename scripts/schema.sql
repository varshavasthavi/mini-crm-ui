-- ClickHouse Event Table
CREATE TABLE IF NOT EXISTS events (
    event_id UUID,
    player_id String,
    event_type String,
    amount Float64,
    timestamp DateTime64(3)
) ENGINE = MergeTree()
ORDER BY timestamp;