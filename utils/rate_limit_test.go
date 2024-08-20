package utils_test

import (
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/utils"
	"github.com/stretchr/testify/require"
)

func TestRateLimit(t *testing.T) {
	mock := &utils.MockTimerProvider{}
	mock.SetNow(time.Now())
	sut := utils.NewRateLimit(utils.NewRateLimitConfig(2, time.Second), mock)
	require.Nil(t, sut.Call("test", false))
	require.Nil(t, sut.Call("test", false))
	sleepTime := sut.Call("test", false)
	require.NotNil(t, sleepTime)
	require.Equal(t, time.Second, *sleepTime)

	mock.SetNow(mock.GetNow().Add(time.Second * 2))
	require.Nil(t, sut.Call("test", false))
	require.Nil(t, sut.Call("test", false))
}

func TestRateLimitSleepTime(t *testing.T) {
	mock := &utils.MockTimerProvider{}
	mock.SetNow(time.Now())
	sut := utils.NewRateLimit(utils.NewRateLimitConfig(2, time.Minute), mock)
	require.Nil(t, sut.Call("test", false))
	mock.SetNow(mock.GetNow().Add(time.Second * 55))
	require.Nil(t, sut.Call("test", false))
	sleepTime := sut.Call("test", false)
	require.NotNil(t, sleepTime)
	require.Equal(t, time.Second*5, *sleepTime)
	mock.SetNow(mock.GetNow().Add(time.Second * 4))
	// It sleeps 1 second
	sut.Call("test", true)
}

func TestRateLimitDisabled(t *testing.T) {
	mock := &utils.MockTimerProvider{}
	mock.SetNow(time.Now())
	sut := utils.NewRateLimit(utils.NewRateLimitConfig(0, time.Minute), mock)
	require.Nil(t, sut.Call("test", false))
	for i := 1; i <= 1000; i++ {
		require.Nil(t, sut.Call("test", false))
	}
}
