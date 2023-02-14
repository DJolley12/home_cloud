package payload

import (
	"bytes"
	"crypto/ed25519"
	"fmt"
	"io"
	"math/rand"

	"filippo.io/age"
	"github.com/DJolley12/home_cloud/server/persist"
)

func decrypt(privK, data []byte) (*bytes.Buffer, error) {
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

func Encrypt(userId int64, data []byte) ([]byte, error) {
	key, err := persist.GetPublicKey(userId)
	if err != nil {
		return nil, err
	}
	return encrypt(key, data)
}

func encrypt(pubK, data []byte) ([]byte, error) {
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
	return out.Bytes(), nil
}

func GenKeyPairForUser(userId int64, userKey []byte) (ed25519.PublicKey, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}
	if err = persist.InsertKeySet(userId, pub, priv, userKey); err != nil {
		return nil, err
	}
	return pub, nil
}

func VerifyPassphrase(encr, actPass string, userId int64) (bool, error) {
	key, err := persist.GetPrivateKey(userId)
	if err != nil {
		return false, err
	}
	pass, err := decrypt(key, []byte(encr))
	if err != nil {
		return false, err
	}
	if string(pass.Bytes()) == actPass {
		return true, nil
	}

	return false, nil
}

func VerifyRefreshToken(userId int64, encr string) error {
	privK, err := persist.GetPrivateKey(userId)
	if err != nil {
		return err
	}
	tk, err := decrypt(privK, []byte(encr))
	if err != nil {
		return err
	}

	dbToken, err := persist.GetRefreshToken(userId)
	if err != nil {
		return err
	}

	// match
	if dbToken == string(tk.Bytes()) {
		return nil
	}

	return fmt.Errorf("token does not match")
}

func MakeRefreshToken(userId int64) (string, error) {
	token := GenerateToken(50)
	key, err := persist.GetPublicKey(userId)
	if err != nil {
		return "", err
	}
	buf, err := encrypt(key, []byte(token))
	if err != nil {
		return "", err
	}
	err = persist.InsertRefreshToken(userId, string(buf.Bytes()))
	if err != nil {
		return "", err
	}
	return string(buf.Bytes()), nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-!@#$%^&*()-=+[]")

func GenerateToken(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
