package payload

import (
	"context"
	"time"

	"github.com/DJolley12/home_cloud/shared/encryption"
)

type tokenCache struct {
	cache       map[int64]userTokenInfo
	sizeLimit   int
	expiryLimit time.Duration
}

type userTokenInfo struct {
	expiry    time.Time
	userId    int64
	sig       []byte
	cryptoKey []byte
	token     string
	sigKey    []byte
}

func newTokenCache(sizeLimit int, expiryLimit time.Duration) tokenCache {
	return tokenCache{
		cache:       make(map[int64]userTokenInfo),
		sizeLimit:   sizeLimit,
		expiryLimit: expiryLimit,
	}
}

func (c *tokenCache) add(userId int64, sig []byte, cryptoKey []byte, token string, sigKey []byte) {
	c.cache[userId] = userTokenInfo{
		expiry:    time.Now(),
		userId:    userId,
		sig:       sig,
		cryptoKey: cryptoKey,
		token:     token,
		sigKey:    sigKey,
	}
}

func (c *tokenCache) tokenIsValid(ctx context.Context) bool {
	// m, ok := metadata.FromIncomingContext(ctx)
	token, ok := ctx.Value("access-token").(string)
	if !ok {
		panic("not ok")
		return false
	}
	userId, ok := ctx.Value("user-id").(int64)
	if !ok {
		panic("not ok")
		return false
	}

	val, ok := c.cache[userId]
	if !ok {
		panic("not ok cache")
		return false
	}

	if val.expiry.Before(time.Now().Add(-c.expiryLimit * time.Minute)) {
		return false
	}

	t, err := encryption.DecryptAndVerify([]byte(token), val.cryptoKey, val.sig, val.sigKey)
	if err != nil {
		return false
	}

	if string(t) == val.token {
		return true
	}

	return false
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
