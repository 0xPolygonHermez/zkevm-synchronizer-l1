-- +migrate Up
CREATE TABLE IF NOT EXISTS sync.sequenced_batches
(
    block_num BIGINT NOT NULL,
    from_batch_num int8 NOT NULL PRIMARY KEY,
    to_batch_num int8 NOT NULL,
    "timestamp" timestamptz NOT NULL,
    l1_info_root VARCHAR NULL,
    CONSTRAINT sequenced_batches_block_num_fkey FOREIGN KEY (block_num) REFERENCES sync.block(block_num) ON DELETE CASCADE

);

-- +migrate Down
DROP TABLE IF EXISTS sync.sequenced_batches;
