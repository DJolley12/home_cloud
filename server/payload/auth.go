package payload

import (
	"bytes"
	"crypto/ed25519"
	"io"

	"filippo.io/age"
	"github.com/DJolley12/home_cloud/server/persist"
)

func Decrypt(privK, data []byte) (*bytes.Buffer, error) {
	rec, err := age.ParseX25519Identity(string(privK))
	if err != nil {
		return nil, err
	}

	out := &bytes.Buffer{}

	r, err := age.Decrypt(out, rec)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(out, r); err != nil {
		return nil, err
	}
	return out, nil
}

func Encrypt(pubK, data []byte) (*bytes.Buffer, error) {
	rec, err := age.ParseX25519Recipient(string(pubK))
	if err != nil {
		return nil, err
	}

	out := &bytes.Buffer{}

	w, err := age.Encrypt(out, rec)
	if err != nil {
		return nil, err
	}
	defer w.Close()
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	return out, nil
}

func GenKeyPairForUser(userId int64, userKey []byte) error {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
	  return err
	}

	return persist.InsertKeySet(userId, pub, priv, userKey)
}
