package syncinterfaces

type SynchronizerIsTrustedSequencer interface {
	IsTrustedSequencer() bool
}

type SynchronizerCleanTrustedState interface {
	CleanTrustedState()
}

type SynchronizerFullInterface interface {
	SynchronizerIsTrustedSequencer
	SynchronizerCleanTrustedState
}
