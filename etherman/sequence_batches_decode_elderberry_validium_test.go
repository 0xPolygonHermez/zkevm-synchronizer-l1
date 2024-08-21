package etherman

import (
	"encoding/hex"
	"testing"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/dataavailability"
	mock_dataavailability "github.com/0xPolygonHermez/zkevm-synchronizer-l1/dataavailability/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	txDataElderberryValidiumHex = "db5b0ed700000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000066708c910000000000000000000000000000000000000000000000000000000000000001000000000000000000000000353800524721e11b453f73f523dd8840c215a21300000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000001923c579cccd3bf22f74f4ba058c0338462434d3d2b9995e292477a8fd3c7bc4d00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000691a4f83394abd1764eb1cab1482950567909146e9d7b2a4a87ee14c5af97c73b970a7c5f8d33de3a347f6c629d7226a0db7b700d9b34213d81561ca1e67980fef1bcc543f5a2052edf584216093a0547c4acd84b80bef5f06e5c0493601829dacfa23f2fe30303b01660000000000000000000000000000000000000000000000"
	expectedElderberryValidium  = `BatchNumber: 2
L1InfoRoot: 0xa59ece237550cae22559a6654efccbd8fc33fb7b56ed11983c97e211ddd17346
SequencerAddr: 0x353800524721e11B453f73f523dD8840c215a213
TxHash: 0x2cd4d0310ae93fe9e068fa3cc1d7bcb1b6ea66374f48efb2b389bb5a64105a15
Nonce: 1
Coinbase: 0x353800524721e11B453f73f523dD8840c215a213
PolygonZkEVMBatchData: nil
___EtrogSequenceData:ForcedTimestamp: 0
___EtrogSequenceData:ForcedGlobalExitRoot: 0000000000000000000000000000000000000000000000000000000000000000
___EtrogSequenceData:ForcedBlockHashL1: 0000000000000000000000000000000000000000000000000000000000000000
___EtrogSequenceData:Transactions: 01020304
___SequencedBatchElderberryData:MaxSequenceTimestamp 1718652049
___SequencedBatchElderberryData:InitSequencedBatchNumber 1
Metadata: SourceBatchData: DA/External RollupFlavor: Validium CallFunctionName: sequenceBatchesElderberryValidium ForkName: elderberry
`
)

type testDataElderberryValidium struct {
	mockDA *mock_dataavailability.BatchDataProvider
	sut    *SequenceBatchesDecodeElderberryValidium
}

func TestSequencedBatchElderberryValidiumDecode(t *testing.T) {
	data := newTestDataElderberryValidium(t)
	txData, err := hex.DecodeString(txDataElderberryValidiumHex)
	require.NoError(t, err)
	batchL2Data := dataavailability.BatchL2Data{
		Data:   []byte{0x01, 0x02, 0x03, 0x04},
		Source: dataavailability.External,
	}

	data.mockDA.EXPECT().GetBatchL2Data([]uint64{2}, mock.Anything, mock.Anything).Return([]dataavailability.BatchL2Data{batchL2Data}, nil)

	res, err := data.sut.DecodeSequenceBatches(txData, uint64(2),
		common.HexToAddress("0x353800524721e11B453f73f523dD8840c215a213"),
		common.HexToHash("0x2cd4d0310ae93fe9e068fa3cc1d7bcb1b6ea66374f48efb2b389bb5a64105a15"),
		uint64(1),
		common.HexToHash("0xa59ece237550cae22559a6654efccbd8fc33fb7b56ed11983c97e211ddd17346"))

	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, 1, len(res))
	res0 := res[0].String()
	require.Equal(t, expectedElderberryValidium, res0)

}

func newTestDataElderberryValidium(t *testing.T) testDataElderberryValidium {
	mockDA := mock_dataavailability.NewBatchDataProvider(t)
	sut, err := NewDecodeSequenceBatchesElderberryValidium(mockDA)
	require.NoError(t, err)
	require.NotNil(t, sut)

	return testDataElderberryValidium{
		mockDA: mockDA,
		sut:    sut,
	}
}
