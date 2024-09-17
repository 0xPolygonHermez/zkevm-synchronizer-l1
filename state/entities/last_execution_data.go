package entities

import (
	"fmt"
	"time"
)

type LastExecutionData struct {
	SyncVersion   string
	LastStart     time.Time
	Configuration string
}

func (s *LastExecutionData) String() string {
	return fmt.Sprintf("{SyncVersion: %s, LastStart: %s}", s.SyncVersion, s.LastStart.String())
}
