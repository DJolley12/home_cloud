package payload

import (
	"bytes"
	"crypto/ed25519"
	"fmt"
	"io"
	"math/rand"
	"strings"

	"filippo.io/age"
	"github.com/DJolley12/home_cloud/server/persist"
)

func decrypt(privK, data []byte) (*bytes.Buffer, error) {
	rec, err := age.ParseX25519Identity(string(privK))
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(data)
	b := make([]byte, 0)
	out := bytes.NewBuffer(b)

	println(buf.String)
	r, err := age.Decrypt(buf, rec)
	if err != nil {
		return nil, err
	}
	_, _ = r.Read(b)
	println(string(b))
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

type keySet struct {
	privSign   []byte
	pubSignKey []byte
	pubEncrKey string
}

// returns server's generated public encryption key, public sign key, error
func GenKeyPairForUser(userId int64, userEncrKey, userSignKey []byte) (*keySet, error) {
	privEncrId, err := age.GenerateX25519Identity()
	if err != nil {
		return nil, err
	}
	pubEncryR := privEncrId.Recipient()

	pubSign, privSign, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	pubEncrStr := pubEncryR.String()
	k := persist.Keys{
		PubEncrKey:  []byte(pubEncrStr),
		PrivEncrKey: []byte(privEncrId.String()),
		PubSignKey:  pubSign,
		PrivSignKey: privSign,
		UserEncrKey: userEncrKey,
		UserSignKey: userSignKey,
	}
	if err = persist.InsertKeySet(userId, k); err != nil {
		return nil, err
	}

	return &keySet{
		privSign:   privSign,
		pubSignKey: pubSign,
		pubEncrKey: pubEncrStr,
	}, nil
}

func VerifyPassphrase(userPass string) (int64, bool, error) {
	return persist.GetUserPassphrase(userPass)
}

func VerifyRefreshToken(userId int64, encr string) error {
	k, err := persist.GetKeys(userId)
	if err != nil {
		return err
	}
	t, err := decrypt(k.PrivEncrKey, []byte(encr))
	if err != nil {
		return err
	}
	spl := strings.Split(t.String(), ".")
	token, sig := spl[0], spl[1]
	if ok := ed25519.Verify(k.PubSignKey, sig, ); !ok {
		return fmt.Errorf("token does not match")
	}

	dbToken, err := persist.GetRefreshToken(userId)
	if err != nil {
		return err
	}

	// match
	if dbToken == string(strings.Split(token, ".")[0]) {
		return nil
	}

	return fmt.Errorf("token does not match")
}

func MakeRefreshToken(userId int64, encryptKey string, signKey []byte) (string, error) {
	token := GenerateToken(50)
	err := persist.InsertRefreshToken(userId, token)
	if err != nil {
		return "", err
	}

	buf, err := encrypt([]byte(encryptKey), ed25519.Sign(signKey, []byte(token)))
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-!@#$%^&*()-=+[]")

func GenerateToken(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func SignId(userId int64, signKey []byte) []byte {
	return ed25519.Sign(signKey, []byte(string(userId)))
}
