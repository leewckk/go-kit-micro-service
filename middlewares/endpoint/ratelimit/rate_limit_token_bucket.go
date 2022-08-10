package ratelimit

import (
	"fmt"
	"time"

	juju "github.com/juju/ratelimit"
)

/// juju/ratelimit的token桶算法
/// https://github.com/juju/ratelimit
type RateLimiterTokenBucketJuju struct {
	bucket *juju.Bucket
}

func (limiter *RateLimiterTokenBucketJuju) Allow() bool {
	/// 只要有可用的token则返回true
	if limiter.bucket.TakeAvailable(1) == 0 {
		return false
	}
	return true
}

func (limiter *RateLimiterTokenBucketJuju) Debug() {
	for {
		available := limiter.bucket.Available()
		fmt.Printf("rate limiter : %p available: %v \n", limiter, available)
		time.Sleep(time.Second * 1)
	}
}

func NewRateLimiterTokenBucketJuju(interval time.Duration, capacity int64) *RateLimiterTokenBucketJuju {
	limiter := &RateLimiterTokenBucketJuju{}
	limiter.bucket = juju.NewBucketWithQuantum(interval, capacity, capacity)
	return limiter
}
