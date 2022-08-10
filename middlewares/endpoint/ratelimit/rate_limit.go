package ratelimit

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
)

const (
	DEFAULT_TOKEN_BUCKET_CAPACITY = 1e5
)

//// 不限流
type RateLimiterUnlimited struct {
}

func (limiter *RateLimiterUnlimited) Allow() bool {
	return true
}

func NewRateLimiterTokenBucketUnmilited() *RateLimiterUnlimited {
	return &RateLimiterUnlimited{}
}

//// limiter 工厂函数
func NewDefaultRateLimiter() ratelimit.Allower {
	/// 默认返回unlimit
	return NewRateLimiterTokenBucketUnmilited()
	// return NewRateLimiterLeakyBucketUber(DEFAULT_TOKEN_BUCKET_CAPACITY)
	// return NewRateLimiterTokenBucketJuju(100*time.Millisecond, DEFAULT_TOKEN_BUCKET_CAPACITY/10)
}

func NewLimitMiddleware(limiter ratelimit.Allower) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			if limiter.Allow() != true {
				return nil, ratelimit.ErrLimited
			}
			return next(ctx, request)
		}
	}
}
