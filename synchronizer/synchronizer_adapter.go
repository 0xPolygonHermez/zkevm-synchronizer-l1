package synchronizer

type SynchronizerAdapter struct {
	*SyncrhronizerQueries
	*SynchronizerImpl
}

func NewSynchronizerAdapter(queries *SyncrhronizerQueries, sync *SynchronizerImpl) *SynchronizerAdapter {
	return &SynchronizerAdapter{
		SyncrhronizerQueries: queries,
		SynchronizerImpl:     sync,
	}
}
