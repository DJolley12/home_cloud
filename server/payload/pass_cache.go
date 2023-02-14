package payload

import (
	"fmt"
	"time"
)

type passCache struct {
	cache       map[int64]userPassInfo
	sizeLimit   int
	expiryLimit time.Duration
}

type userPassInfo struct {
	expiry     time.Time
	passphrase string
}

func newPassCache(sizeLimit int, expiryLimit time.Duration) passCache {
	return passCache{
		cache:       make(map[int64]userPassInfo),
		sizeLimit:   sizeLimit,
		expiryLimit: expiryLimit,
	}
}

func (c *passCache) passIsValid(userId int64) (string, error) {
	// m, ok := metadata.FromIncomingContext(ctx)
	val, ok := c.cache[userId]
	if !ok {
		return "", fmt.Errorf("no password found for user %#v")
	}

	if val.expiry.Before(time.Now().Add(-c.expiryLimit * time.Minute)) {
		return "", fmt.Errorf("passphrase has expired")
	}

	return val.passphrase, nil
}

func (c *passCache) collectPass() {
	if len(c.cache) < c.sizeLimit {
		return
	}

	for k, v := range c.cache {
		if v.expiry.Before(time.Now().Add(-10 * time.Minute)) {
			delete(c.cache, k)
		}
	}
}
