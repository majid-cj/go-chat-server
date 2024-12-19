package middleware

import (
	"sync"
	"time"

	"github.com/kataras/iris/v12"
	"golang.org/x/time/rate"
)

// Create a custom visitor struct which holds the rate limiter for each
// visitor and the last time that the visitor was seen.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Change the the map to hold values of the type visitor.
var context = sync.Map{}

// Run a background goroutine to remove old entries from the visitors map.
func init() {
	go cleanupVisitors()
}

func getVisitor(ip string) *rate.Limiter {
	v, exists := context.Load(ip)
	value, _ := v.(*visitor)
	if !exists {
		limiter := rate.NewLimiter(5, 10)
		// Include the current time when creating a new visitor.
		context.Store(ip, &visitor{limiter, time.Now()})
		return limiter
	}

	// Update the last seen time for the visitor.
	value.lastSeen = time.Now()
	return value.limiter
}

// Every minute check the map for visitors that haven't been seen for
// more than 3 minutes and delete the entries.
func cleanupVisitors() {
	for {
		context.Range(func(key, value interface{}) bool {
			data, ok := value.(*visitor)
			if ok {
				if time.Since(data.lastSeen) > time.Minute {
					context.Delete(key)
				}
			}
			return true
		})
	}
}

// RateLimit ...
func RateLimit(c iris.Context) {
	ip := c.Request().RemoteAddr

	limiter := getVisitor(ip)

	if !limiter.Allow() {
		c.StatusCode(429)
		return
	}

	c.Next()
}
