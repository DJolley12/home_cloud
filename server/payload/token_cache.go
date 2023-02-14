package payload

import (
	"context"
	"time"
)

type tokenCache struct {
	cache       map[string]userTokenInfo
	sizeLimit   int
	expiryLimit time.Duration
}

type userTokenInfo struct {
	expiry time.Time
	userId int64
}

func newTokenCache(sizeLimit int, expiryLimit time.Duration) tokenCache {
	return tokenCache{
		cache: make(map[string]userTokenInfo),
		sizeLimit: sizeLimit,
		expiryLimit: expiryLimit,
	}
}

func (c *tokenCache) add(token string, userId int64) {
	c.cache[token] = userTokenInfo{
		expiry: time.Now(),
		userId: userId,
	}
}

func (c *tokenCache) tokenIsValid(ctx context.Context) bool {
	// m, ok := metadata.FromIncomingContext(ctx)
	token, ok := ctx.Value("access-token").(string)
	if !ok {
		panic("not ok")
		return false
	}

	val, ok := c.cache[token]
	if !ok {
		panic("not ok cache")
		return false
	}

	if val.expiry.Before(time.Now().Add(-c.expiryLimit * time.Minute)) {
		return false
	}

	return true
}

func (c *tokenCache) collectTokens() {
	if len(c.cache) < c.sizeLimit {
		return
	}

	for k, v := range c.cache {
		if v.expiry.Before(time.Now().Add(-10 * time.Minute)) {
			delete(c.cache, k)
		}
	}
}
