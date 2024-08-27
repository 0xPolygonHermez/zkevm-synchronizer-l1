package etherman

import "github.com/ethereum/go-ethereum/crypto"

var (
	// Events RollupManager
	verifyBatchesTrustedAggregatorSignatureHash = crypto.Keccak256Hash([]byte("VerifyBatchesTrustedAggregator(uint32,uint64,bytes32,bytes32,address)"))
	updateRollupSignatureHash                   = crypto.Keccak256Hash([]byte("UpdateRollup(uint32,uint32,uint64)"))
	addExistingRollupSignatureHash              = crypto.Keccak256Hash([]byte("AddExistingRollup(uint32,uint64,address,uint64,uint8,uint64)"))
	createNewRollupSignatureHash                = crypto.Keccak256Hash([]byte("CreateNewRollup(uint32,uint32,address,uint64,address)"))
	// https://github.com/0xPolygonHermez/zkevm-contracts/blob/949b0b96c10056fa7be9632bcc2f26202a9c3a9c/contracts/v2/PolygonRollupManager.sol#L328
	//rollbackBatchesManagerSignatureHash = crypto.Keccak256Hash([]byte("RollbackBatches(uint32,uint64,bytes32)"))
	// Events new ZkEvm/RollupBase
	sequenceForceBatchesSignatureHash   = crypto.Keccak256Hash([]byte("SequenceForceBatches(uint64)"))             // Used in oldZkEvm as well
	forceBatchSignatureHash             = crypto.Keccak256Hash([]byte("ForceBatch(uint64,bytes32,address,bytes)")) // Used in oldZkEvm as well
	sequenceBatchesSignatureHash        = crypto.Keccak256Hash([]byte("SequenceBatches(uint64,bytes32)"))          // Used in oldZkEvm as well
	initialSequenceBatchesSignatureHash = crypto.Keccak256Hash([]byte("InitialSequenceBatches(bytes,bytes32,address)"))
	updateEtrogSequenceSignatureHash    = crypto.Keccak256Hash([]byte("UpdateEtrogSequence(uint64,bytes,bytes32,address)"))
	// https://github.com/0xPolygonHermez/zkevm-contracts/blob/949b0b96c10056fa7be9632bcc2f26202a9c3a9c/contracts/v2/lib/PolygonRollupBaseEtrog.sol#L589
	rollbackBatchesSignatureHash = crypto.Keccak256Hash([]byte("RollbackBatches(uint64,bytes32)")) // targetBatch, accInputHashToRollback

	// New GER event
	updateL1InfoTreeSignatureHash   = crypto.Keccak256Hash([]byte("UpdateL1InfoTree(bytes32,bytes32)"))
	updateL1InfoTreeV2SignatureHash = crypto.Keccak256Hash([]byte("UpdateL1InfoTreeV2(bytes32,uint32,uint256,uint64)"))
	//initL1InfoRootMapSignatureHash  = crypto.Keccak256Hash([]byte("InitL1InfoRootMap(uint32,bytes32)"))

	// PreLxLy events
	updateGlobalExitRootSignatureHash              = crypto.Keccak256Hash([]byte("UpdateGlobalExitRoot(bytes32,bytes32)"))
	oldVerifyBatchesTrustedAggregatorSignatureHash = crypto.Keccak256Hash([]byte("VerifyBatchesTrustedAggregator(uint64,bytes32,address)"))
	updateZkEVMVersionSignatureHash                = crypto.Keccak256Hash([]byte("UpdateZkEVMVersion(uint64,uint64,string)"))
	sequenceBatchesPreEtrogSignatureHash           = crypto.Keccak256Hash([]byte("SequenceBatches(uint64)"))

	// methodIDSequenceBatchesEtrog: MethodID for sequenceBatches in Etrog
	methodIDSequenceBatchesEtrog = []byte{0xec, 0xef, 0x3f, 0x99} // 0xecef3f99
	// methodIDSequenceBatchesElderberry: MethodID for sequenceBatches in Elderberry
	// https://github.com/0xPolygonHermez/zkevm-contracts/blob/73758334f8568b74e9493fcc530b442bd73325dc/contracts/v2/lib/PolygonRollupBaseEtrog.sol#L32
	methodIDSequenceBatchesElderberry = []byte{0xde, 0xf5, 0x7e, 0x54} // 0xdef57e54 sequenceBatches((bytes,bytes32,uint64,bytes32)[],uint64,uint64,address)

	signatures = []string{
		"SetBatchFee(uint256)",
		"SetTrustedAggregator(address)",
		"SetVerifyBatchTimeTarget(uint64)",
		"SetMultiplierBatchFee(uint16)",
		"SetPendingStateTimeout(uint64)",
		"SetTrustedAggregatorTimeout(uint64)",
		"OverridePendingState(uint32,uint64,bytes32,bytes32,address)",
		"ProveNonDeterministicPendingState(bytes32,bytes32)",
		"ConsolidatePendingState(uint32,uint64,bytes32,bytes32,uint64)",
		"VerifyBatchesTrustedAggregator(uint32,uint64,bytes32,bytes32,address)",
		"VerifyBatches(uint32,uint64,bytes32,bytes32,address)",
		"OnSequenceBatches(uint32,uint64)",
		"UpdateRollup(uint32,uint32,uint64)",
		"AddExistingRollup(uint32,uint64,address,uint64,uint8,uint64)",
		"CreateNewRollup(uint32,uint32,address,uint64,address)",
		"ObsoleteRollupType(uint32)",
		"AddNewRollupType(uint32,address,address,uint64,uint8,bytes32,string)",
		"AcceptAdminRole(address)",
		"TransferAdminRole(address)",
		"SetForceBatchAddress(address)",
		"SetForceBatchTimeout(uint64)",
		"SetTrustedSequencerURL(string)",
		"SetTrustedSequencer(address)",
		"VerifyBatches(uint64,bytes32,address)",
		"SequenceForceBatches(uint64)",
		"ForceBatch(uint64,bytes32,address,bytes)",
		"SequenceBatches(uint64,bytes32)",
		"InitialSequenceBatches(bytes,bytes32,address)",
		"UpdateEtrogSequence(uint64,bytes,bytes32,address)",
		"Initialized(uint64)",
		"RoleAdminChanged(bytes32,bytes32,bytes32)",
		"RoleGranted(bytes32,address,address)",
		"RoleRevoked(bytes32,address,address)",
		"EmergencyStateActivated()",
		"EmergencyStateDeactivated()",
		"UpdateL1InfoTree(bytes32,bytes32)",
		"UpdateGlobalExitRoot(bytes32,bytes32)",
		"VerifyBatchesTrustedAggregator(uint64,bytes32,address)",
		"OwnershipTransferred(address,address)",
		"UpdateZkEVMVersion(uint64,uint64,string)",
		"ConsolidatePendingState(uint64,bytes32,uint64)",
		"OverridePendingState(uint64,bytes32,address)",
		"SequenceBatches(uint64)",
		"Initialized(uint8)",
		"AdminChanged(address,address)",
		"BeaconUpgraded(address)",
		"Upgraded(address)",
		"RollbackBatches(uint64,bytes32)",
		"SetDataAvailabilityProtocol(address)",
		"UpdateL1InfoTreeV2(bytes32,uint32,uint256,uint64)",
		"InitL1InfoRootMap(uint32,bytes32)",
		"RollbackBatches(uint32,uint64,bytes32)",
	}
)
