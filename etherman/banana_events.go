package etherman

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (etherMan *Client) processBananaEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) (bool, error) {
	switch vLog.Topics[0] {
	case rollbackBatchesSignatureHash:
		return true, etherMan.rollbackBatchesManagerEvent(ctx, vLog, blocks, blocksOrder)
	case updateL1InfoTreeV2SignatureHash:
		return true, etherMan.updateL1InfoTreeV2Event(ctx, vLog, blocks, blocksOrder)
	}
	return false, nil
}

func (etherMan *Client) rollbackBatchesManagerEvent(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {
	/*
			   event RollbackBatches(
		        uint32 indexed rollupID,
		        uint64 indexed targetBatch,
		        bytes32 accInputHashToRollback
		    );
	*/
	eventData, err := etherMan.BananaZkEVM.ParseRollbackBatches(vLog)
	if err != nil {
		log.Warnf("error parsing RollbackBatches event: %v", err)
		return err
	}

	block, err := addNewBlockToResult(ctx, etherMan, vLog, blocks, blocksOrder)
	if err != nil {
		log.Warnf("error addNewBlockToResult RollbackBatches event: %v", err)
		return err
	}
	rollbackBatchesData := RollbackBatchesData{
		TargetBatch:            eventData.TargetBatch,
		AccInputHashToRollback: eventData.AccInputHashToRollback,
	}
	block.RollbackBatches = append(block.RollbackBatches, rollbackBatchesData)
	order := Order{
		Name: RollbackBatchesOrder,
		Pos:  len(block.RollbackBatches) - 1,
	}
	addNewOrder(&order, block.BlockHash, blocksOrder)
	return nil
}

func (etherMan *Client) updateL1InfoTreeV2Event(ctx context.Context, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) error {

	/* https://github.com/0xPolygonHermez/zkevm-contracts/blob/949b0b96c10056fa7be9632bcc2f26202a9c3a9c/contracts/v2/PolygonZkEVMGlobalExitRootV2.sol#L39C1-L44C7

		    event UpdateL1InfoTreeV2(
	        bytes32 currentL1InfoRoot,
	        uint32 indexed leafCount,
	        uint256 blockhash,
	        uint64 minTimestamp
	    );
	*/
	eventData, err := etherMan.GlobalExitRootManager.ParseUpdateL1InfoTreeV2(vLog)
	if err != nil {
		return err
	}
	block, err := addNewBlockToResult(ctx, etherMan, vLog, blocks, blocksOrder)
	if err != nil {
		return err
	}
	L1InfoTreeV2Data := L1InfoTreeV2Data{
		CurrentL1InfoRoot: eventData.CurrentL1InfoRoot,
		LeafCount:         eventData.LeafCount,
		// TODO: Fix this type
		//BlockHash:         eventData.Blockhash,
		MinTimestamp: eventData.MinTimestamp,
	}
	block.L1InfoTreeV2 = append(block.L1InfoTreeV2, L1InfoTreeV2Data)
	order := Order{
		Name: UpdateL1InfoTreeV2Order,
		Pos:  len(block.L1InfoTreeV2) - 1,
	}
	addNewOrder(&order, block.BlockHash, blocksOrder)
	return nil
}
