package entities

import "fmt"

type StorageContentsBoundData struct {
	RollupID  uint64
	L1ChainID uint64
}

func (s *StorageContentsBoundData) String() string {
	return fmt.Sprintf("{RollupID: %d, L1ChainID: %d}", s.RollupID, s.L1ChainID)
}
