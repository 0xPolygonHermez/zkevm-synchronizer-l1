-- +migrate Down


-- +migrate Up


-- History






CREATE TABLE block
(
    block_num   BIGINT PRIMARY KEY,
    block_hash  VARCHAR(66) NOT NULL,
    parent_hash VARCHAR(66),
    received_at TIMESTAMP  NOT NULL, -- It is the creation time of the block (not the received time)
	checked BOOLEAN NOT NULL DEFAULT FALSE, -- If it is true this block is not going to be reorg
	has_events BOOLEAN NOT NULL DEFAULT FALSE,
	sync_version VARCHAR(128)
);


CREATE TABLE exit_root (
	block_num  BIGINT REFERENCES block (block_num) ON DELETE CASCADE,
	"timestamp" TIMESTAMP NOT NULL,
	mainnet_exit_root VARCHAR(66) NULL,
	rollup_exit_root VARCHAR(66) NULL,
	global_exit_root VARCHAR(66) NULL,
	prev_block_hash VARCHAR(66) NULL,
	l1_info_root VARCHAR(66) NULL,
	l1_info_tree_index int8 NULL,
	CONSTRAINT exit_root_pkey PRIMARY KEY (l1_info_tree_index)
);

CREATE INDEX idx_exit_root_l1_info_tree_index ON exit_root(l1_info_tree_index);
CREATE INDEX idx_exit_root_global_exit_root_index ON exit_root(global_exit_root);


CREATE TABLE fork_id (
	fork_id int8 NOT NULL,
	from_batch_num numeric NOT NULL,
	to_batch_num numeric NOT NULL,
	"version" varchar(128) NULL,
	block_num BIGINT REFERENCES block (block_num) ON DELETE CASCADE,
	CONSTRAINT fork_id_pkey PRIMARY KEY (fork_id)
);
