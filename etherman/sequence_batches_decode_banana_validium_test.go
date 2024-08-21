package etherman

import (
	"encoding/hex"
	"testing"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/dataavailability"
	mock_dataavailability "github.com/0xPolygonHermez/zkevm-synchronizer-l1/dataavailability/mocks"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	txDataBananaValidiumHex   = "165e8a8d00000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000066c43b8e443448cabd1b7808f64eb3f2737180b7ef27a0a28e19afebb2981714a6492c9b0000000000000000000000005b06837a43bdc3dd9f114558daf4b26ed49842ed00000000000000000000000000000000000000000000000000000000000002600000000000000000000000000000000000000000000000000000000000000003d09d1238668f9e4dd74d503f15a90c8fc230069b5ace4a3327896b50bf809db5000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000df729b193d68d2b455c59f7040687d4dc2d147557442aeef4c10cc5769d47ade00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000014df0077a11ccf22bd14d1ec3e209ba5d8e065ce73556fdc8fb9c19b83af4ec30000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000055f10249e353edebbaa84755e6216d6bf44a1be5e5b140ed6cd94f56dcd7611ad2202708675673a014fcea40c1519ae920ed7652e23e43bf837e11fc5c403539fd1b5951f5b2604c9b42e478d5e2b2437f44073ef9a60000000000000000000000"
	expectedSeqBananaValidium = `BatchNumber: 2
L1InfoRoot: 0xa59ece237550cae22559a6654efccbd8fc33fb7b56ed11983c97e211ddd17346
SequencerAddr: 0x353800524721e11B453f73f523dD8840c215a213
TxHash: 0x2cd4d0310ae93fe9e068fa3cc1d7bcb1b6ea66374f48efb2b389bb5a64105a15
Nonce: 1
Coinbase: 0x5b06837A43bdC3dD9F114558DAf4B26ed49842Ed
PolygonZkEVMBatchData: nil
___EtrogSequenceData:ForcedTimestamp: 0
___EtrogSequenceData:ForcedGlobalExitRoot: 0000000000000000000000000000000000000000000000000000000000000000
___EtrogSequenceData:ForcedBlockHashL1: 0000000000000000000000000000000000000000000000000000000000000000
___EtrogSequenceData:Transactions: 0b0000008b000000000b00000006000000000b0000000600000000
SequencedBatchElderberryData: nil
BananaData: CounterL1InfoRoot: 1 MaxSequenceTimestamp: 1724136334 ExpectedFinalAccInputHash: 0x443448cabd1b7808f64eb3f2737180b7ef27a0a28e19afebb2981714a6492c9b DataAvailabilityMsg(85): f10249e353edebbaa84755e6216d6bf44a1be5e5b140ed6cd94f56dcd7611ad2202708675673a014fcea40c1519ae920ed7652e23e43bf837e11fc5c403539fd1b5951f5b2604c9b42e478d5e2b2437f44073ef9a6
Metadata: SourceBatchData: DA/External RollupFlavor: Validium CallFunctionName: sequenceBatchesBananaValidium ForkName: banana
`
)

var (
	batchDataBananaValidiumHex = [...]string{
		"0b0000008b000000000b00000006000000000b0000000600000000",
		"0b00000000000000000b0000000600000000",
		"0b00000006000000000b00000006000000000b0000000600000000",
	}
)

func TestSequencedBatchBananaValidiumDecode(t *testing.T) {
	mockDA := mock_dataavailability.NewBatchDataProvider(t)
	sut, err := NewSequenceBatchesDecoderBananaValidium(mockDA)
	require.NoError(t, err)
	require.NotNil(t, sut)

	txData, err := hex.DecodeString(txDataBananaValidiumHex)
	require.NoError(t, err)
	batchData := make([][]byte, len(batchDataBananaValidiumHex))
	for i, v := range batchDataBananaValidiumHex {
		batchData[i], err = hex.DecodeString(v)
		require.NoError(t, err)
	}

	resultGetBatchL2Data := []dataavailability.BatchL2Data{
		{
			Data:   batchData[0],
			Source: dataavailability.External,
		},
		{
			Data:   batchData[0],
			Source: dataavailability.External,
		},
		{
			Data:   batchData[0],
			Source: dataavailability.External,
		},
	}
	mockDA.EXPECT().GetBatchL2Data([]uint64{2, 3, 4}, mock.Anything, mock.Anything).Return(resultGetBatchL2Data, nil)

	require.NoError(t, err)
	res, err := sut.DecodeSequenceBatches(txData, uint64(4),
		common.HexToAddress("0x353800524721e11B453f73f523dD8840c215a213"),
		common.HexToHash("0x2cd4d0310ae93fe9e068fa3cc1d7bcb1b6ea66374f48efb2b389bb5a64105a15"),
		uint64(1),
		common.HexToHash("0xa59ece237550cae22559a6654efccbd8fc33fb7b56ed11983c97e211ddd17346"))

	require.NoError(t, err)
	require.NotNil(t, res)
	res0 := res[0].String()
	log.Debug(res0)
	require.Equal(t, expectedSeqBananaValidium, res0)
}
