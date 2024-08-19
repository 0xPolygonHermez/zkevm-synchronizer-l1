package datacommittee

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygon/cdk-data-availability/client"
	mock_dataavailability "github.com/0xPolygonHermez/zkevm-synchronizer-l1/dataavailability/mocks"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/smartcontracts/polygondatacommittee"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockFactory struct {
	MockClient map[string]*mock_dataavailability.DACClientMock
}

func (m *MockFactory) New(url string) client.Client {
	return m.MockClient[url]
}

type testDataCommitteeBackend struct {
	mockTimeProvider *utils.MockTimerProvider
	mockFactory      *MockFactory
	ctx              context.Context
	skipOnErrorTime  time.Duration
	hashData         []common.Hash
	batchData        [][]byte
	sut              *DataCommitteeBackend
}

func newTestDataCommiteeBackend(t *testing.T) *testDataCommitteeBackend {
	mockTimeProvider := &utils.MockTimerProvider{}
	skipOnErrorTime := time.Second
	mock := &MockFactory{
		MockClient: map[string]*mock_dataavailability.DACClientMock{
			"url1": mock_dataavailability.NewDACClientMock(t),
			"url2": mock_dataavailability.NewDACClientMock(t),
		},
	}
	sut := DataCommitteeBackend{
		dataCommitteeContract:      nil,
		dataCommitteeClientFactory: mock,
		committeeMembers: []DataCommitteeMember{
			{
				URL:  "url1",
				Addr: common.HexToAddress("0x1"),
			},
			{
				URL:  "url2",
				Addr: common.HexToAddress("0x2"),
			},
		},

		NoReloadCommitteMembersOnError: true,
		committeeMemberControl:         NewDataCommitteeMemberControl(mockTimeProvider, skipOnErrorTime, utils.RateLimitConfig{}),
	}
	ctx := context.TODO()
	batchData := [][]byte{{0, 1, 2, 3}, {4, 5, 6, 7}}
	hashes := []common.Hash{
		crypto.Keccak256Hash(batchData[0]),
		crypto.Keccak256Hash(batchData[1]),
	}
	return &testDataCommitteeBackend{
		mockTimeProvider: mockTimeProvider,
		mockFactory:      mock,
		ctx:              ctx,
		skipOnErrorTime:  skipOnErrorTime,
		hashData:         hashes,
		batchData:        batchData,
		sut:              &sut,
	}
}

func TestDACServerNoError(t *testing.T) {
	testData := newTestDataCommiteeBackend(t)
	testData.mockFactory.MockClient["url1"].EXPECT().GetOffChainData(testData.ctx, testData.hashData[0]).Return(testData.batchData[0], nil)
	testData.mockFactory.MockClient["url1"].EXPECT().GetOffChainData(testData.ctx, testData.hashData[1]).Return(testData.batchData[1], nil)
	_, err := testData.sut.GetSequence(testData.ctx, testData.hashData, []byte{0, 1, 2, 3})
	require.NoError(t, err)
}

func TestDACServerErrorServer1ButServer2ReturnsOk(t *testing.T) {
	testData := newTestDataCommiteeBackend(t)
	testData.mockFactory.MockClient["url1"].EXPECT().GetOffChainData(testData.ctx, testData.hashData[0]).Return(testData.batchData[0], assert.AnError)
	testData.mockFactory.MockClient["url2"].EXPECT().GetOffChainData(testData.ctx, testData.hashData[0]).Return(testData.batchData[0], nil)
	testData.mockFactory.MockClient["url2"].EXPECT().GetOffChainData(testData.ctx, testData.hashData[1]).Return(testData.batchData[1], nil)
	_, err := testData.sut.GetSequence(testData.ctx, testData.hashData, []byte{0, 1, 2, 3})
	require.NoError(t, err)
}

func TestDACServerErrorServer1AndServer2(t *testing.T) {
	testData := newTestDataCommiteeBackend(t)
	testData.mockFactory.MockClient["url1"].EXPECT().GetOffChainData(testData.ctx, testData.hashData[0]).Return(testData.batchData[0], assert.AnError)
	testData.mockFactory.MockClient["url2"].EXPECT().GetOffChainData(testData.ctx, testData.hashData[0]).Return(testData.batchData[0], assert.AnError)
	_, err := testData.sut.GetSequence(testData.ctx, testData.hashData, []byte{0, 1, 2, 3})
	require.Error(t, err)
}

func TestDACServerErrorServer1AndServer2NoFastRetries(t *testing.T) {
	testData := newTestDataCommiteeBackend(t)
	testData.mockTimeProvider.SetNow(time.Now())

	testData.mockFactory.MockClient["url1"].EXPECT().GetOffChainData(testData.ctx, testData.hashData[0]).Return(testData.batchData[0], assert.AnError).Once()
	testData.mockFactory.MockClient["url2"].EXPECT().GetOffChainData(testData.ctx, testData.hashData[0]).Return(testData.batchData[0], assert.AnError).Once()
	_, err := testData.sut.GetSequence(testData.ctx, testData.hashData, []byte{0, 1, 2, 3})
	require.Error(t, err)
	// This doesnt call to clients because < SkipOnErrorTime
	_, err = testData.sut.GetSequence(testData.ctx, testData.hashData, []byte{0, 1, 2, 3})
	require.Error(t, err)

	// We pass the SkipOnErrorTime skip period and must request again to the DAC server
	testData.mockTimeProvider.SetNow(time.Now().Add(testData.skipOnErrorTime))
	testData.mockFactory.MockClient["url1"].EXPECT().GetOffChainData(testData.ctx, testData.hashData[0]).Return(testData.batchData[0], nil).Once()
	testData.mockFactory.MockClient["url1"].EXPECT().GetOffChainData(testData.ctx, testData.hashData[1]).Return(testData.batchData[1], nil).Once()
	_, err = testData.sut.GetSequence(testData.ctx, testData.hashData, []byte{0, 1, 2, 3})
	require.NoError(t, err)

}

func TestUpdateDataCommitteeEvent(t *testing.T) {
	// Set up testing environment
	dac, ethBackend, auth, da := newTestingEnv(t)

	// Update the committee
	requiredAmountOfSignatures := big.NewInt(2)
	URLs := []string{"1", "2", "3"}
	addrs := []common.Address{
		common.HexToAddress("0x1"),
		common.HexToAddress("0x2"),
		common.HexToAddress("0x3"),
	}
	addrsBytes := []byte{}
	for _, addr := range addrs {
		addrsBytes = append(addrsBytes, addr.Bytes()...)
	}
	_, err := da.SetupCommittee(auth, requiredAmountOfSignatures, URLs, addrsBytes)
	require.NoError(t, err)
	ethBackend.Commit()

	// Assert the committee update
	actualSetup, err := dac.getCurrentDataCommittee()
	require.NoError(t, err)
	expectedMembers := []DataCommitteeMember{}
	expectedSetup := DataCommittee{
		RequiredSignatures: uint64(len(URLs) - 1),
		AddressesHash:      crypto.Keccak256Hash(addrsBytes),
	}
	for i, url := range URLs {
		expectedMembers = append(expectedMembers, DataCommitteeMember{
			URL:  url,
			Addr: addrs[i],
		})
	}
	expectedSetup.Members = expectedMembers
	assert.Equal(t, expectedSetup, *actualSetup)
}

func init() {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stderr"},
	})
}

// This function prepare the blockchain, the wallet with funds and deploy the smc
func newTestingEnv(t *testing.T) (
	dac *DataCommitteeBackend,
	ethBackend *simulated.Backend,
	auth *bind.TransactOpts,
	da *polygondatacommittee.Polygondatacommittee,
) {
	t.Helper()
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	auth, err = bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		log.Fatal(err)
	}
	dac, ethBackend, da, err = newSimulatedDacman(t, auth)
	if err != nil {
		log.Fatal(err)
	}
	return dac, ethBackend, auth, da
}

// NewSimulatedEtherman creates an etherman that uses a simulated blockchain. It's important to notice that the ChainID of the auth
// must be 1337. The address that holds the auth will have an initial balance of 10 ETH
func newSimulatedDacman(t *testing.T, auth *bind.TransactOpts) (
	dacman *DataCommitteeBackend,
	ethBackend *simulated.Backend,
	da *polygondatacommittee.Polygondatacommittee,
	err error,
) {
	t.Helper()
	if auth == nil {
		// read only client
		return &DataCommitteeBackend{}, nil, nil, nil
	}
	// 10000000 ETH in wei
	balance, _ := new(big.Int).SetString("10000000000000000000000000", 10) //nolint:gomnd
	address := auth.From
	genesisAlloc := map[common.Address]types.Account{
		address: {
			Balance: balance,
		},
	}
	blockGasLimit := uint64(999999999999999999) //nolint:gomnd
	client := simulated.NewBackend(genesisAlloc, simulated.WithBlockGasLimit(blockGasLimit))

	// DAC Setup
	_, _, da, err = polygondatacommittee.DeployPolygondatacommittee(auth, client.Client())
	if err != nil {
		return &DataCommitteeBackend{}, nil, nil, err
	}
	client.Commit()
	_, err = da.Initialize(auth)
	if err != nil {
		return &DataCommitteeBackend{}, nil, nil, err
	}
	client.Commit()
	_, err = da.SetupCommittee(auth, big.NewInt(0), []string{}, []byte{})
	if err != nil {
		return &DataCommitteeBackend{}, nil, nil, err
	}
	client.Commit()

	c := &DataCommitteeBackend{
		dataCommitteeContract: da,
	}
	return c, client, da, nil
}
