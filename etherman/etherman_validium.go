package etherman

import (
	"crypto/ecdsa"
	"fmt"

	dataCommitteeClient "github.com/0xPolygon/cdk-data-availability/client"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/dataavailability"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/dataavailability/datacommittee"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/smartcontracts/dataavailabilityprotocol"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/smartcontracts/etrogvalidiumpolygonzkevm"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/jsonrpcclient"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type validiumContractBind = *etrogvalidiumpolygonzkevm.Etrogvalidiumpolygonzkevm
type dataAvailabilityProtocolContractBind = *dataavailabilityprotocol.Dataavailabilityprotocol

type EthermanValidium struct {
	Cfg                              Config
	ZkEVMValidiumContract            validiumContractBind
	DataAvailabilityProtocolAddress  common.Address
	DataAvailabilityProtocolContract dataAvailabilityProtocolContractBind
	DataAvailabilityClient           dataavailability.BatchDataProvider
}

func NewEthermanValidium(cfg Config, ethClient bind.ContractBackend) (*EthermanValidium, error) {
	zkevmValidum, err := newZkevmValidiumContractBind(cfg.Contracts.ZkEVMAddr, ethClient)
	if err != nil {
		return nil, err
	}
	DAProtocolAddr, err := getDAProtocolAddr(zkevmValidum)
	if err != nil {
		return nil, err
	}
	daContract, err := newDAProtocolContractBind(DAProtocolAddr, ethClient)
	if err != nil {
		return nil, err
	}

	res := &EthermanValidium{
		Cfg:                              cfg,
		ZkEVMValidiumContract:            zkevmValidum,
		DataAvailabilityProtocolContract: daContract,
		DataAvailabilityProtocolAddress:  DAProtocolAddr,
	}
	da, err := res.newDataAvailabilityClient()
	if err != nil {
		return nil, err
	}
	res.DataAvailabilityClient = da
	return res, nil
}

func newZkevmValidiumContractBind(addr common.Address, ethClient bind.ContractBackend) (*etrogvalidiumpolygonzkevm.Etrogvalidiumpolygonzkevm, error) {
	zkevmValidum, err := etrogvalidiumpolygonzkevm.NewEtrogvalidiumpolygonzkevm(addr, ethClient)
	if err != nil {
		err = fmt.Errorf("error binding zkevm Validium contract (%s)  %w", addr.String(), err)
		log.Errorf(err.Error())
		return nil, err
	}
	return zkevmValidum, nil
}

func getDAProtocolAddr(zkevmValidum validiumContractBind) (common.Address, error) {
	addr, err := zkevmValidum.DataAvailabilityProtocol(&bind.CallOpts{Pending: false})
	if err != nil {
		return common.Address{}, fmt.Errorf("error getting DataAvailabilityProtocol address: %w", err)
	}
	return addr, nil
}

func newDAProtocolContractBind(dapAddr common.Address, ethClient bind.ContractBackend) (*dataavailabilityprotocol.Dataavailabilityprotocol, error) {
	dap, err := dataavailabilityprotocol.NewDataavailabilityprotocol(dapAddr, ethClient)
	if err != nil {
		return nil, err
	}
	return dap, nil
}

func (ev *EthermanValidium) GetTrustedSequencerURL() (string, error) {
	if ev.Cfg.Validium.TrustedSequencerURL == "" {
		url, err := ev.ZkEVMValidiumContract.TrustedSequencerURL(&bind.CallOpts{Pending: false})
		if err != nil {
			return "", fmt.Errorf("error getting trusted sequencer URL: %w", err)
		}
		return url, nil
	}
	return ev.Cfg.Validium.TrustedSequencerURL, nil
}

// GetDAProtocolName returns the name of the data availability protocol
func (ev *EthermanValidium) GetDAProtocolName() (string, error) {
	return ev.DataAvailabilityProtocolContract.GetProcotolName(&bind.CallOpts{Pending: false})
}

func (ev *EthermanValidium) newDataAvailabilityClient() (*dataavailability.DataAvailability, error) {
	var (
		dataSourcePriority []dataavailability.DataSourcePriority
	)
	trustedURL, err := ev.GetTrustedSequencerURL()
	if err != nil {
		return nil, fmt.Errorf("error getting trusted sequencer URL: %w", err)
	}
	log.Debugf("Creating Trusted Sequencer Client with URL: %s", trustedURL)
	trustedRPCClient := jsonrpcclient.NewClient(trustedURL)

	// TODO: Configurable data source priority
	dataSourcePriority = dataavailability.DefaultPriority

	// Backend specific config
	daProtocolName, err := ev.GetDAProtocolName()
	if err != nil {
		return nil, fmt.Errorf("error getting data availability protocol name: %w", err)
	}
	log.Debugf("Data Availability Protocol: %s", daProtocolName)
	var daBackend dataavailability.DABackender
	switch daProtocolName {
	case string(dataavailability.DataAvailabilityCommittee):
		var (
			pk  *ecdsa.PrivateKey
			err error
		)

		dacAddr := ev.DataAvailabilityProtocolAddress

		daBackend, err = datacommittee.New(
			ev.Cfg.L1URL,
			dacAddr,
			pk,
			dataCommitteeClient.NewFactory(),
		)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unexpected / unsupported DA protocol: %s", daProtocolName)
	}

	return dataavailability.New(
		false,
		daBackend,
		trustedRPCClient,
		dataSourcePriority,
	)
}
