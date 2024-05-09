-- +migrate Up
CREATE TABLE IF NOT EXISTS sync.kv (
    key VARCHAR(256) PRIMARY KEY,
    value VARCHAR
);


-- +migrate Down

DROP TABLE IF EXISTS sync.kv;