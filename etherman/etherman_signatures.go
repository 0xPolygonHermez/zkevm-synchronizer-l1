package etherman

import "github.com/ethereum/go-ethereum/crypto"

var (
	// Events RollupManager
	setBatchFeeSignatureHash                       = crypto.Keccak256Hash([]byte("SetBatchFee(uint256)"))
	setTrustedAggregatorSignatureHash              = crypto.Keccak256Hash([]byte("SetTrustedAggregator(address)"))       // Used in oldZkEvm as well
	setVerifyBatchTimeTargetSignatureHash          = crypto.Keccak256Hash([]byte("SetVerifyBatchTimeTarget(uint64)"))    // Used in oldZkEvm as well
	setMultiplierBatchFeeSignatureHash             = crypto.Keccak256Hash([]byte("SetMultiplierBatchFee(uint16)"))       // Used in oldZkEvm as well
	setPendingStateTimeoutSignatureHash            = crypto.Keccak256Hash([]byte("SetPendingStateTimeout(uint64)"))      // Used in oldZkEvm as well
	setTrustedAggregatorTimeoutSignatureHash       = crypto.Keccak256Hash([]byte("SetTrustedAggregatorTimeout(uint64)")) // Used in oldZkEvm as well
	overridePendingStateSignatureHash              = crypto.Keccak256Hash([]byte("OverridePendingState(uint32,uint64,bytes32,bytes32,address)"))
	proveNonDeterministicPendingStateSignatureHash = crypto.Keccak256Hash([]byte("ProveNonDeterministicPendingState(bytes32,bytes32)")) // Used in oldZkEvm as well
	consolidatePendingStateSignatureHash           = crypto.Keccak256Hash([]byte("ConsolidatePendingState(uint32,uint64,bytes32,bytes32,uint64)"))
	verifyBatchesTrustedAggregatorSignatureHash    = crypto.Keccak256Hash([]byte("VerifyBatchesTrustedAggregator(uint32,uint64,bytes32,bytes32,address)"))
	rollupManagerVerifyBatchesSignatureHash        = crypto.Keccak256Hash([]byte("VerifyBatches(uint32,uint64,bytes32,bytes32,address)"))
	onSequenceBatchesSignatureHash                 = crypto.Keccak256Hash([]byte("OnSequenceBatches(uint32,uint64)"))
	updateRollupSignatureHash                      = crypto.Keccak256Hash([]byte("UpdateRollup(uint32,uint32,uint64)"))
	addExistingRollupSignatureHash                 = crypto.Keccak256Hash([]byte("AddExistingRollup(uint32,uint64,address,uint64,uint8,uint64)"))
	createNewRollupSignatureHash                   = crypto.Keccak256Hash([]byte("CreateNewRollup(uint32,uint32,address,uint64,address)"))
	obsoleteRollupTypeSignatureHash                = crypto.Keccak256Hash([]byte("ObsoleteRollupType(uint32)"))
	addNewRollupTypeSignatureHash                  = crypto.Keccak256Hash([]byte("AddNewRollupType(uint32,address,address,uint64,uint8,bytes32,string)"))

	// Events new ZkEvm/RollupBase
	acceptAdminRoleSignatureHash        = crypto.Keccak256Hash([]byte("AcceptAdminRole(address)"))                 // Used in oldZkEvm as well
	transferAdminRoleSignatureHash      = crypto.Keccak256Hash([]byte("TransferAdminRole(address)"))               // Used in oldZkEvm as well
	setForceBatchAddressSignatureHash   = crypto.Keccak256Hash([]byte("SetForceBatchAddress(address)"))            // Used in oldZkEvm as well
	setForceBatchTimeoutSignatureHash   = crypto.Keccak256Hash([]byte("SetForceBatchTimeout(uint64)"))             // Used in oldZkEvm as well
	setTrustedSequencerURLSignatureHash = crypto.Keccak256Hash([]byte("SetTrustedSequencerURL(string)"))           // Used in oldZkEvm as well
	setTrustedSequencerSignatureHash    = crypto.Keccak256Hash([]byte("SetTrustedSequencer(address)"))             // Used in oldZkEvm as well
	verifyBatchesSignatureHash          = crypto.Keccak256Hash([]byte("VerifyBatches(uint64,bytes32,address)"))    // Used in oldZkEvm as well
	sequenceForceBatchesSignatureHash   = crypto.Keccak256Hash([]byte("SequenceForceBatches(uint64)"))             // Used in oldZkEvm as well
	forceBatchSignatureHash             = crypto.Keccak256Hash([]byte("ForceBatch(uint64,bytes32,address,bytes)")) // Used in oldZkEvm as well
	sequenceBatchesSignatureHash        = crypto.Keccak256Hash([]byte("SequenceBatches(uint64,bytes32)"))          // Used in oldZkEvm as well
	initialSequenceBatchesSignatureHash = crypto.Keccak256Hash([]byte("InitialSequenceBatches(bytes,bytes32,address)"))
	updateEtrogSequenceSignatureHash    = crypto.Keccak256Hash([]byte("UpdateEtrogSequence(uint64,bytes,bytes32,address)"))
	rollbackBatchesSignatureHash        = crypto.Keccak256Hash([]byte("RollbackBatches(uint64,bytes32)"))

	// Extra RollupValidiumEtrog
	setDataAvailabilitySignatureHash = crypto.Keccak256Hash([]byte("SetDataAvailabilityProtocol(address)"))

	// Extra RollupManager
	initializedSignatureHash               = crypto.Keccak256Hash([]byte("Initialized(uint64)"))                       // Initializable. Used in RollupBase as well
	roleAdminChangedSignatureHash          = crypto.Keccak256Hash([]byte("RoleAdminChanged(bytes32,bytes32,bytes32)")) // IAccessControlUpgradeable
	roleGrantedSignatureHash               = crypto.Keccak256Hash([]byte("RoleGranted(bytes32,address,address)"))      // IAccessControlUpgradeable
	roleRevokedSignatureHash               = crypto.Keccak256Hash([]byte("RoleRevoked(bytes32,address,address)"))      // IAccessControlUpgradeable
	emergencyStateActivatedSignatureHash   = crypto.Keccak256Hash([]byte("EmergencyStateActivated()"))                 // EmergencyManager. Used in oldZkEvm as well
	emergencyStateDeactivatedSignatureHash = crypto.Keccak256Hash([]byte("EmergencyStateDeactivated()"))               // EmergencyManager. Used in oldZkEvm as well

	// New GER event
	updateL1InfoTreeSignatureHash   = crypto.Keccak256Hash([]byte("UpdateL1InfoTree(bytes32,bytes32)"))
	updateL1InfoTreeV2SignatureHash = crypto.Keccak256Hash([]byte("UpdateL1InfoTreeV2(bytes32,uint32,uint256,uint64)"))
	initL1InfoRootMapSignatureHash  = crypto.Keccak256Hash([]byte("InitL1InfoRootMap(uint32,bytes32)"))

	// PreLxLy events
	updateGlobalExitRootSignatureHash              = crypto.Keccak256Hash([]byte("UpdateGlobalExitRoot(bytes32,bytes32)"))
	oldVerifyBatchesTrustedAggregatorSignatureHash = crypto.Keccak256Hash([]byte("VerifyBatchesTrustedAggregator(uint64,bytes32,address)"))
	transferOwnershipSignatureHash                 = crypto.Keccak256Hash([]byte("OwnershipTransferred(address,address)"))
	updateZkEVMVersionSignatureHash                = crypto.Keccak256Hash([]byte("UpdateZkEVMVersion(uint64,uint64,string)"))
	oldConsolidatePendingStateSignatureHash        = crypto.Keccak256Hash([]byte("ConsolidatePendingState(uint64,bytes32,uint64)"))
	oldOverridePendingStateSignatureHash           = crypto.Keccak256Hash([]byte("OverridePendingState(uint64,bytes32,address)"))
	sequenceBatchesPreEtrogSignatureHash           = crypto.Keccak256Hash([]byte("SequenceBatches(uint64)"))

	// Proxy events
	initializedProxySignatureHash = crypto.Keccak256Hash([]byte("Initialized(uint8)"))
	adminChangedSignatureHash     = crypto.Keccak256Hash([]byte("AdminChanged(address,address)"))
	beaconUpgradedSignatureHash   = crypto.Keccak256Hash([]byte("BeaconUpgraded(address)"))
	upgradedSignatureHash         = crypto.Keccak256Hash([]byte("Upgraded(address)"))

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
	}
)
