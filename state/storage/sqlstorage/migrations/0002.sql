-- +migrate Up
CREATE TABLE IF NOT EXISTS sequenced_batches
(
    block_num BIGINT NOT NULL,
    from_batch_num BIGINT NOT NULL PRIMARY KEY,
    to_batch_num BIGINT NOT NULL,
    fork_id BIGINT NOT NULL,
    "timestamp" TIMESTAMP  NOT NULL,
    received_at TIMESTAMP  NOT NULL,
    l1_info_root VARCHAR(66) NULL,
    source VARCHAR(128) NULL,   -- it store the origin of this sequence
    CONSTRAINT sequenced_batches_block_num_fkey FOREIGN KEY (block_num) REFERENCES block(block_num) ON DELETE CASCADE
);


-- +migrate Down
DROP TABLE IF EXISTS sync.sequenced_batches;
