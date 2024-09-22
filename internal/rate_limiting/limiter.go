package rate_limiting

import (
	"sync"
	"time"
)

// Credits per window
const MaxCreditsPerSecond = 2000
const CreditsWindow = time.Millisecond * 250
const MaxCredits = MaxCreditsPerSecond / int(time.Second/CreditsWindow)

// Rate limiter using a fixed window counter
type RateLimiter struct {
	lock      sync.Mutex
	window    time.Duration
	timestamp time.Time
	credits   int
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		window:    CreditsWindow,
		timestamp: time.Now(),
		credits:   MaxCredits,
	}
}

func (c *RateLimiter) AllowedAfter(requestCredits int, retryWait ...time.Duration) time.Duration {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.allowed(requestCredits) {
		var retryAfter time.Duration

		if len(retryWait) > 0 {
			retryAfter = retryWait[0]
		} else {
			retryAfter = c.window + 50*time.Millisecond
		}

		if retryAfter == 0 {
			return time.Until(c.timestamp.Add(c.window))
		} else {
			// TODO: implement sliding window
			// For now sleeps with the lock to avoid starvation
			time.Sleep(retryAfter)

			if !c.allowed(requestCredits) {
				return time.Until(c.timestamp.Add(c.window))
			}
		}
	}

	return 0
}

func (c *RateLimiter) allowed(requestCredits int) bool {
	switch {
	case c.credits-requestCredits > 0:
		c.credits -= requestCredits
	case time.Since(c.timestamp) > c.window:
		c.timestamp = time.Now()
		c.credits = MaxCredits - requestCredits
	default:
		return false
	}

	return true
}
