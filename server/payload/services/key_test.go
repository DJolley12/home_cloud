package services

import (
	"crypto/ed25519"
	"testing"

	"filippo.io/age"
	"github.com/stretchr/testify/assert"
)

func TestKeyEncrypt(t *testing.T) {
	t.Run("test encrypt and sign, decrypt and verify", func(t *testing.T) {
		identity, err := age.GenerateX25519Identity()
		assert.Nil(t, err)
		recipient := identity.Recipient()

		pub, priv, _ := ed25519.GenerateKey(nil)
		data := []byte("test test testing")

		tokenSig, err := encryptAndSign(data, []byte(recipient.String()), priv)
		assert.NotEqual(t, string(data), string(tokenSig.Token))

		c, err := DecryptAndVerify(tokenSig.Token, []byte(identity.String()), tokenSig.Signature, pub)
		assert.Equal(t, "test test testing", string(c))
	})
}
