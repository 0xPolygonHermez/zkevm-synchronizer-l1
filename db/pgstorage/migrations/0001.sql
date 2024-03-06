-- +migrate Down
DROP SCHEMA IF EXISTS sync CASCADE;

-- +migrate Up
CREATE SCHEMA sync;

-- History






CREATE TABLE sync.block
(
    block_num   BIGINT PRIMARY KEY,
    block_hash  VARCHAR NOT NULL,
    parent_hash VARCHAR,
    
    received_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE sync.exit_root (
	id serial4 NOT NULL,
	block_num int8 NOT NULL,
	"timestamp" timestamptz NOT NULL,
	mainnet_exit_root VARCHAR NULL,
	rollup_exit_root VARCHAR NULL,
	global_exit_root VARCHAR NULL,
	prev_block_hash VARCHAR NULL,
	l1_info_root VARCHAR NULL,
	l1_info_tree_index int8 NULL,
	CONSTRAINT exit_root_l1_info_tree_index_key UNIQUE (l1_info_tree_index),
	CONSTRAINT exit_root_pkey PRIMARY KEY (id)
);

CREATE INDEX idx_exit_root_l1_info_tree_index ON sync.exit_root USING btree (l1_info_tree_index);
CREATE INDEX idx_exit_root_global_exit_root_index ON sync.exit_root USING btree (global_exit_root);


-- CREATE TABLE sync.exit_root
-- (
--     id                      SERIAL,
--     block_id                BIGINT REFERENCES sync.block (id) ON DELETE CASCADE,
--     global_exit_root        VARCHAR,
--     exit_roots              BYTEA[],
--     PRIMARY KEY (id),
--     CONSTRAINT UC UNIQUE (block_id, global_exit_root)
-- );

-- CREATE TABLE sync.batch
-- (
--     batch_num            BIGINT PRIMARY KEY,
--     sequencer            VARCHAR,
--     raw_tx_data          VARCHAR, 
--     global_exit_root     VARCHAR,
--     timestamp            TIMESTAMP WITH TIME ZONE
-- );

-- CREATE TABLE IF NOT EXISTS sync.rollup_exit
-- (
-- 	id        BIGSERIAL PRIMARY KEY,
--     leaf      BYTEA,
--     rollup_id BIGINT,
-- 	root      BYTEA,
-- 	block_id BIGINT NOT NULL REFERENCES sync.block (id) ON DELETE CASCADE
-- );
