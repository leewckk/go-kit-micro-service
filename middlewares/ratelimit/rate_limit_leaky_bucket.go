package ratelimit

import uber "go.uber.org/ratelimit"

/// 基于uber ratelimit的漏桶算法
/// https://github.com/uber-go/ratelimit
type RateLimiterLeakyBucketUber struct {
	limiter uber.Limiter
}

func (limiter *RateLimiterLeakyBucketUber) Allow() bool {
	limiter.limiter.Take() //// 漏桶算法自带延时，延时到期后返回，所以返回的都是可以allow的
	return true
}

func NewRateLimiterLeakyBucketUber(rps int) *RateLimiterLeakyBucketUber {
	limiter := &RateLimiterLeakyBucketUber{}
	limiter.limiter = uber.New(rps)
	return limiter
}
