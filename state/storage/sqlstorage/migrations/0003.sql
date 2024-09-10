-- +migrate Up


CREATE TABLE IF NOT EXISTS virtual_batch (
    -- This is a sequenced batch from L1
    batch_num BIGINT PRIMARY KEY,
    fork_id BIGINT NOT NULL,
    raw_txs_data     BYTEA,
    vlog_tx_hash   VARCHAR(66), -- hash of Tx inside L1 block with that vlog
    coinbase  VARCHAR(42),
    sequence_from_batch_num BIGINT NOT NULL REFERENCES sequenced_batches (from_batch_num) ON DELETE CASCADE,
    block_num BIGINT NOT NULL REFERENCES block (block_num) ON DELETE CASCADE,
	sequencer_addr VARCHAR(42) NOT NULL,
    received_at TIMESTAMP  NOT NULL,
    l1_info_root VARCHAR(66) NULL,
    extra_info VARCHAR NULL,
    batch_timestamp TIMESTAMP  NULL, -- node: timestamp_batch_etrog
    sync_version VARCHAR(128)
);


CREATE TABLE IF NOT EXISTS reorg_log (
    -- Logs of the reorg on DB
    timestamp timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    batch_num BIGINT NULL,
    block_num BIGINT NULL,
    reason VARCHAR NOT NULL,
    extra_info VARCHAR NULL,
    CONSTRAINT trusted_reorg_pkey PRIMARY KEY ("timestamp")
);


-- +migrate Down
DROP TABLE IF EXISTS virtual_batch;
DROP TABLE IF EXISTS reorg_log;