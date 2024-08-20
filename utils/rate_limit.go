package utils

import (
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/config/types"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
)

type RateLimitConfig struct {
	NumRequests int
	Interval    types.Duration
}

func NewRateLimitConfig(numRequests int, period time.Duration) RateLimitConfig {
	return RateLimitConfig{
		NumRequests: numRequests,
		Interval:    types.Duration{Duration: period},
	}
}

func (r RateLimitConfig) String() string {
	return fmt.Sprintf("RateLimitConfig{NumRequests: %d, Period: %s}", r.NumRequests, r.Interval)
}

func (r RateLimitConfig) Enabled() bool {
	return r.NumRequests > 0 && r.Interval.Duration > 0
}

type RateLimit struct {
	cfg RateLimitConfig

	timeProvider TimeProvider
	// Calls realized in the current period
	bucket []time.Time
}

func NewRateLimit(cfg RateLimitConfig, timeProvider TimeProvider) RateLimit {
	return RateLimit{
		cfg:          cfg,
		timeProvider: timeProvider,
	}
}

// This is a call
func (r *RateLimit) Call(msg string, allowToSleep bool) *time.Duration {
	if !r.cfg.Enabled() {
		return nil
	}
	var returnSleepTime *time.Duration
	now := r.timeProvider.Now()
	r.cleanOutdatedCalls(now)
	if len(r.bucket) >= r.cfg.NumRequests {
		sleepTime := r.cfg.Interval.Duration - r.timeProvider.Now().Sub(r.bucket[0])
		if allowToSleep {
			if msg != "" {
				log.Debugf("Rate limit reached, sleeping for %s for %s", sleepTime, msg)
			}
			time.Sleep(sleepTime)
		} else {
			// If no sleep, ignore the call
			return &sleepTime
		}
		returnSleepTime = &sleepTime
	}
	r.bucket = append(r.bucket, now)
	return returnSleepTime
}

func (r *RateLimit) cleanOutdatedCalls(now time.Time) {
	for i, call := range r.bucket {
		diff := now.Sub(call)
		if diff < r.cfg.Interval.Duration {
			r.bucket = r.bucket[i:]
			return
		}
	}
	r.bucket = []time.Time{}
}
