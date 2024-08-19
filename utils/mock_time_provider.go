package utils

import "time"

// MockTimerProvider is a mock implementation of the TimerProvider interface that return the internal variable
type MockTimerProvider struct {
	now time.Time
}

func (m *MockTimerProvider) SetNow(now time.Time) {
	m.now = now
}
func (m *MockTimerProvider) GetNow() time.Time {
	return m.now
}

// Now in the implementation of TimeProvider.Now()
func (m *MockTimerProvider) Now() time.Time {
	return m.now
}
