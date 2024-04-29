package state

import (
	"context"
	"testing"

	mock_state "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/mocks"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/storage/pgstorage"
	"github.com/stretchr/testify/assert"
)

func TestGetL1InfoTreeLeaves(t *testing.T) {
	// Create a mock storage implementation
	mockStorage := mock_state.NewStorageL1InfoTreeInterface(t)

	// Create a new instance of L1InfoTreeState
	state := NewL1InfoTreeManager(mockStorage)

	// Define the expected result
	expectedResult := map[uint32]L1InfoTreeLeaf{
		10: L1InfoTreeLeaf{
			// Set the fields of the expected leaf
		},
		20: L1InfoTreeLeaf{
			// Set the fields of the expected leaf
		},
		// Add more expected leaves if needed
	}
	mockStorage.EXPECT().GetL1InfoLeafPerIndex(context.Background(), uint32(10), nil).Return(&pgstorage.L1InfoTreeLeaf{}, nil)
	mockStorage.EXPECT().GetL1InfoLeafPerIndex(context.Background(), uint32(20), nil).Return(&pgstorage.L1InfoTreeLeaf{}, nil)
	// Call the GetL1InfoTreeLeaves function
	result, err := state.GetL1InfoTreeLeaves(context.Background(), []uint32{10, 20}, nil)

	// Check for errors
	assert.NoError(t, err)

	// Compare the result with the expected result
	assert.Equal(t, expectedResult, result)
}
