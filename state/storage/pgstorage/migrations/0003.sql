-- +migrate Up


CREATE TABLE IF NOT EXISTS sync.virtual_batch (

    batch_num BIGINT PRIMARY KEY,
    fork_id BIGINT NOT NULL,
    raw_txs_data     BYTEA,
    vlog_tx_hash   VARCHAR(66),
    coinbase  VARCHAR(42),
    sequence_from_batch_num BIGINT NOT NULL REFERENCES sync.sequenced_batches (from_batch_num) ON DELETE CASCADE,
    block_num BIGINT NOT NULL REFERENCES sync.block (block_num) ON DELETE CASCADE,
	sequencer_addr VARCHAR(42) NOT NULL,
    received_at TIMESTAMP WITH TIME ZONE NOT NULL,
    l1_info_root VARCHAR(66) NULL,
    extra_info VARCHAR NULL,
    batch_timestamp TIMESTAMP WITH TIME ZONE NULL, -- node: timestamp_batch_etrog
    sync_version VARCHAR(128)
);

comment on column sync.virtual_batch.vlog_tx_hash is 'hash of Tx inside L1 block with that vlog';
comment on table sync.virtual_batch is 'This is a sequenced batch from L1';

CREATE TABLE IF NOT EXISTS sync.reorg_log (
    "timestamp" timestamptz DEFAULT now() NOT NULL,
    batch_num BIGINT NULL,
    block_num BIGINT NULL,
    reason VARCHAR NOT NULL,
    extra_info VARCHAR NULL,
    CONSTRAINT trusted_reorg_pkey PRIMARY KEY ("timestamp")
);
comment on table sync.reorg_log is 'Logs of the reorg on DB';

-- +migrate Down
DROP TABLE IF EXISTS sync.virtual_batch;
DROP TABLE IF EXISTS sync.reorg_log;