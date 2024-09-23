-- +migrate Up
CREATE TABLE IF NOT EXISTS kv (
    key VARCHAR(256) PRIMARY KEY,
    value VARCHAR, 
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    sync_version VARCHAR(128) -- 'Version of the library that make the last change'
);


-- +migrate Down

DROP TABLE IF EXISTS kv;
