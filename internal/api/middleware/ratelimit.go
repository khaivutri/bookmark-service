package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/internal/repository/ratelimit"
	"github.com/khaivutri/bookmark-service/pkg/requestutils"
	"github.com/rs/zerolog/log"
)

const (
	rateLimitInterval  = 1 * time.Minute
	rateLimitCount     = 10
	rateLimitKeyFormat = "rate_limit:%s"
)

type RateLimiter interface {
	RateLimit() gin.HandlerFunc
	LimitByUser(limit int, interval time.Duration) gin.HandlerFunc
	LimitByIP(limit int, interval time.Duration) gin.HandlerFunc
}

type rateLimiter struct {
	repo ratelimit.RateLimiter
}

// NewRateLimiter returns a new rate limiter instance using the given repository.
// The repository is used to store and retrieve the rate limit data.
// The returned rate limiter implements the RateLimiter interface and can be used to limit the number of requests from a user or IP address.
// The default rate limit configuration is 10 requests per minute, but this can be customized using the LimitByUser and LimitByIP methods.
func NewRateLimiter(repo ratelimit.RateLimiter) RateLimiter {
	return &rateLimiter{repo: repo}
}

// RateLimit returns a middleware that limits the number of requests from a user
// to a given endpoint path using the default rate limit configuration (10 requests per minute).
// The rate limit is specified by the limit parameter and the interval parameter specifies the time duration
// in which the limit is applied.
func (r *rateLimiter) RateLimit() gin.HandlerFunc {
	return r.LimitByUser(rateLimitCount, rateLimitInterval)
}

// LimitByUser returns a middleware that limits the number of requests from a user
// to a given endpoint path. The rate limit is specified by the limit parameter and the interval
// parameter specifies the time duration in which the limit is applied.
func (r *rateLimiter) LimitByUser(limit int, interval time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := requestutils.GetUserIDFromRequest(c)
		if err != nil {
			return
		}

		r.limit(c, userID, limit, interval)
	}
}


// LimitByIP returns a middleware that limits the number of requests from a client IP
// to a given endpoint path. The rate limit is specified by the limit parameter and the interval
// parameter specifies the time duration in which the limit is applied.
func (r *rateLimiter) LimitByIP(limit int, interval time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		rateLimitKey := fmt.Sprintf("%s:%s", c.ClientIP(), path)
		r.limit(c, rateLimitKey, limit, interval)
	}
}

// limit checks the current rate limit for the given key and aborts the request if it exceeds the limit.
// Otherwise, it increases the rate limit by 1 and continues the request.
func (r *rateLimiter) limit(c *gin.Context, key string, limit int, interval time.Duration) {
	rateLimitKey := fmt.Sprintf(rateLimitKeyFormat, key)

	curRate, err := r.repo.GetCurrentRateLimit(c, rateLimitKey)
	if err != nil {
		log.Error().Err(err).Msg("fail to get current rate limit")
	}

	if curRate >= limit {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"message": "rate limit exceeded",
		})
		return
	}

	r.repo.IncreaseRateLimit(c, rateLimitKey, interval)
	c.Next()
}
