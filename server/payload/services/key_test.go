package services

import (
	"crypto/ed25519"
	"testing"

	"filippo.io/age"
	"github.com/stretchr/testify/assert"
)

func TestKeyEncrypt(t *testing.T) {
	t.Run("test encryptAndSign, decryptAndVerify",
		func(t *testing.T) {
			ident, err := age.GenerateX25519Identity()
			assert.Nil(t, err)

			pubSign, privSign, err := ed25519.GenerateKey(nil)

			tokenSig, err := encryptAndSign([]byte("testicles, 1, 2"), []byte(ident.Recipient().String()), privSign)
			assert.Nil(t, err)

			b, err := DecryptAndVerify([]byte(ident.String()), tokenSig.Token, tokenSig.Signature, pubSign)
			assert.Nil(t, err)

			assert.Equal(t, "testicles, 1, 2", string(b))
		})
}
