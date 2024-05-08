package l1_check_block

//go:generate bash -c "rm -Rf mocks"
//go:generate mockery --all --case snake --dir . --output ./mocks --outpkg mock_l1_check_block  --disable-version-string --with-expecter

import "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"

type L1Block = entities.L1Block
type stateTxType = entities.Tx
