-- +migrate Up
CREATE TABLE IF NOT EXISTS sync.kv (
    key VARCHAR(256) PRIMARY KEY,
    value VARCHAR, 
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    sync_version VARCHAR(128)
);

comment on column sync.kv.sync_version is 'Version of the library that make the last change';

-- +migrate Down

DROP TABLE IF EXISTS sync.kv;