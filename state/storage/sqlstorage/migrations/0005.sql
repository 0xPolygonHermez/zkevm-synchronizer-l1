-- +migrate Up
CREATE TABLE IF NOT EXISTS rollback_batches_log (
    id VARCHAR(66) PRIMARY KEY, -- Unique identifier of the rollback event (hash block, batch, accinputhash)
	block_num  BIGINT REFERENCES block (block_num) ON DELETE CASCADE,
    last_batch_number BIGINT NOT NULL, -- Corresponding to targetBatch of L1 event
    last_batch_acc_input_hash VARCHAR(66) NOT NULL, -- Corresponding to accInputHashToRollback of L1 event
    description VARCHAR,
    l1_event_at TIMESTAMP NOT NULL, -- is the L1 timestamp of this event
    received_at TIMESTAMP NOT NULL, -- received_at is the execution time of local sync
    undo_first_block_num BIGINT NOT NULL, -- First block number to sync to undo the rollback
    sequences_deleted JSONB, 
    sync_version VARCHAR(128) -- Version of the library that make this rollback
);

-- +migrate Down

DROP TABLE IF EXISTS sync.rollback_batches_log;