-- +migrate Up
CREATE TABLE IF NOT EXISTS sync.rollback_batches_log (
    id VARCHAR(66) PRIMARY KEY,
	block_num  BIGINT REFERENCES sync.block (block_num) ON DELETE CASCADE,
    last_batch_number BIGINT NOT NULL,
    last_batch_acc_input_hash VARCHAR(66) NOT NULL,
    description VARCHAR,
    l1_event_at TIMESTAMP WITH TIME ZONE NOT NULL,
    received_at TIMESTAMP WITH TIME ZONE NOT NULL,
    undo_first_block_num BIGINT NOT NULL,
    sequences_deleted JSONB, 
    sync_version VARCHAR(128)
);
comment on column sync.rollback_batches_log.id is 'Unique identifier of the rollback event (hash block, batch, accinputhash)';
comment on column sync.rollback_batches_log.l1_event_at is 'is the L1 timestamp of this event'; 
comment on column sync.rollback_batches_log.received_at is 'received_at is the execution time of local sync'; 
comment on column sync.rollback_batches_log.undo_first_block_num is 'First block number to sync to undo the rollback'; 
comment on column sync.rollback_batches_log.last_batch_number is 'Corresponding to targetBatch of L1 event';
comment on column sync.rollback_batches_log.last_batch_acc_input_hash is 'Corresponding to accInputHashToRollback of L1 event';
comment on column sync.rollback_batches_log.sync_version is 'Version of the library that make this rollback';


-- +migrate Down

DROP TABLE IF EXISTS sync.rollback_batches_log;