-- +migrate Up
CREATE TABLE IF NOT EXISTS sync.sequenced_batches
(
    block_num BIGINT NOT NULL,
    from_batch_num BIGINT NOT NULL PRIMARY KEY,
    to_batch_num BIGINT NOT NULL,
    fork_id BIGINT NOT NULL,
    "timestamp" TIMESTAMP WITH TIME ZONE NOT NULL,
    received_at TIMESTAMP WITH TIME ZONE NOT NULL,
    l1_info_root VARCHAR(66) NULL,
    source VARCHAR(128) NULL,
    CONSTRAINT sequenced_batches_block_num_fkey FOREIGN KEY (block_num) REFERENCES sync.block(block_num) ON DELETE CASCADE
);

comment on column sync.sequenced_batches.source is 'it store the origin of this sequence';


-- +migrate Down
DROP TABLE IF EXISTS sync.sequenced_batches;
