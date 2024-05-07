-- +migrate Down
DROP SCHEMA IF EXISTS sync CASCADE;

-- +migrate Up
CREATE SCHEMA sync;

-- History






CREATE TABLE sync.block
(
    block_num   BIGINT PRIMARY KEY,
    block_hash  VARCHAR(66) NOT NULL,
    parent_hash VARCHAR(66),
    received_at TIMESTAMP WITH TIME ZONE NOT NULL,
	checked BOOLEAN NOT NULL DEFAULT FALSE,
	has_events BOOLEAN NOT NULL DEFAULT FALSE,
	sync_version VARCHAR(128)
);

comment on column sync.block.checked is 'If it is true this block is not going to be reorg';
comment on column sync.block.received_at is 'It is the creation time of the block (not the received time)';


CREATE TABLE sync.exit_root (
	id serial4 NOT NULL,
	block_num  BIGINT REFERENCES sync.block (block_num) ON DELETE CASCADE,
	"timestamp" timestamptz NOT NULL,
	mainnet_exit_root VARCHAR(66) NULL,
	rollup_exit_root VARCHAR(66) NULL,
	global_exit_root VARCHAR(66) NULL,
	prev_block_hash VARCHAR(66) NULL,
	l1_info_root VARCHAR(66) NULL,
	l1_info_tree_index int8 NULL,
	CONSTRAINT exit_root_l1_info_tree_index_key UNIQUE (l1_info_tree_index),
	CONSTRAINT exit_root_pkey PRIMARY KEY (id)
);

CREATE INDEX idx_exit_root_l1_info_tree_index ON sync.exit_root USING btree (l1_info_tree_index);
CREATE INDEX idx_exit_root_global_exit_root_index ON sync.exit_root USING btree (global_exit_root);


CREATE TABLE sync.fork_id (
	fork_id int8 NOT NULL,
	from_batch_num numeric NOT NULL,
	to_batch_num numeric NOT NULL,
	"version" varchar(128) NULL,
	block_num BIGINT REFERENCES sync.block (block_num) ON DELETE CASCADE,
	CONSTRAINT fork_id_pkey PRIMARY KEY (fork_id)
);
