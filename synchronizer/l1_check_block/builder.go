package l1_check_block

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	syncconfig "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/config"
)

// NewL1CheckBlockFeature creates the feature, basically is an instance of L1BlockCheckerIntegration
func NewL1CheckBlockFeature(cfg syncconfig.L1BlockCheckConfig, ethMan interface{}, state interface{}, sync interface{}) *L1BlockCheckerIntegration {
	if !cfg.Enable {
		return nil
	}
	log.Infof("L1BlockChecker enabled: %s", cfg.String())
	mainChecker := NewAsyncMainL1CheckBlock(cfg, ethMan, state)

	var preCheckAsync AsyncL1BlockChecker
	if cfg.PreCheckEnable {
		log.Infof("L1BlockChecker enabled precheck from: %s/%d to: %s/%d",
			cfg.L1SafeBlockPoint, cfg.L1SafeBlockOffset,
			cfg.L1PreSafeBlockPoint, cfg.L1PreSafeBlockOffset)

		preCheckAsync = NewAsyncPreL1CheckBlock(cfg, ethMan, state)
	}

	return NewL1BlockCheckerIntegration(
		mainChecker,
		preCheckAsync,
		state.(StateForL1BlockCheckerIntegration),
		sync.(SyncCheckReorger),
		cfg.ForceCheckBeforeStart,
		time.Second)
}

// NewAsyncMainL1CheckBlock create main ASsync Cheker(from xxxx(...)xxxxxxxx[safe]-------------[preSafe]--------->)
func NewAsyncMainL1CheckBlock(cfg syncconfig.L1BlockCheckConfig, ethMan interface{}, state interface{}) AsyncL1BlockChecker {
	blockFetcher := NewSafeL1BlockNumberFetch(StringToL1BlockPoint(cfg.L1SafeBlockPoint), cfg.L1SafeBlockOffset)
	l1BlockChecker := NewCheckL1BlockHash(ethMan.(L1Requester), state.(StateInterfacer), blockFetcher)
	return NewAsyncCheck(l1BlockChecker)
}

// NewAsyncPreL1CheckBlock create pre ASsync PreChecker (from ---------[safe]xxxxxxxxxxxx[preSafe]--------->)
func NewAsyncPreL1CheckBlock(cfg syncconfig.L1BlockCheckConfig, ethMan interface{}, state interface{}) AsyncL1BlockChecker {
	initialFetcher := NewSafeL1BlockNumberFetch(StringToL1BlockPoint(cfg.L1SafeBlockPoint), cfg.L1SafeBlockOffset)
	endFetcher := NewSafeL1BlockNumberFetch(StringToL1BlockPoint(cfg.L1PreSafeBlockPoint), cfg.L1PreSafeBlockOffset)
	l1BlockChecker := NewPreCheckL1BlockHash(ethMan.(L1Requester), state.(StatePreCheckInterfacer), initialFetcher, endFetcher)
	return NewAsyncCheck(l1BlockChecker)
}
