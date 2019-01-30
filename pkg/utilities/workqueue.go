package utilities

import (
	"k8s.io/client-go/util/workqueue"
	"time"
)

func NewNamedRateLimitingQueue(rateLimiter workqueue.RateLimiter, name string) *DelayWithRateLimitQueue {
	return &DelayWithRateLimitQueue{
		DelayingInterface: workqueue.NewNamedDelayingQueue(name),
		rateLimiter:       rateLimiter,
	}
}

// rateLimitingType wraps an Interface and provides rateLimited re-enquing
type DelayWithRateLimitQueue struct {
	workqueue.DelayingInterface

	rateLimiter workqueue.RateLimiter
}

// AddRateLimited AddAfter's the item based on the time when the rate limiter says it's ok
func (q *DelayWithRateLimitQueue) AddRateLimited(item interface{}) {
	q.DelayingInterface.AddAfter(item, q.rateLimiter.When(item))
}

// AddRateLimited AddAfter's the item based on the time when the rate limiter says it's ok
func (q *DelayWithRateLimitQueue) AddDelayDefined(item interface{}, duration time.Duration) {
	q.DelayingInterface.AddAfter(item, duration)
}

func (q *DelayWithRateLimitQueue) NumRequeues(item interface{}) int {
	return q.rateLimiter.NumRequeues(item)
}

func (q *DelayWithRateLimitQueue) Forget(item interface{}) {
	q.rateLimiter.Forget(item)
}

