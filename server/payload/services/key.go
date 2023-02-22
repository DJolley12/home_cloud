package services

import (
	"bytes"
	"crypto/ed25519"
	"fmt"
	"io"
	"math/rand"
	"os"

	"filippo.io/age"
	"github.com/DJolley12/home_cloud/server/persist/ports"
)

type KeyService struct {
	userPersist ports.UserPersist
}

func decrypt(privK, data []byte) ([]byte, error) {
	fmt.Fprintf(os.Stderr, "in decrypt data: %v\n", string(data))
	rec, err := age.ParseX25519Identity(string(privK))
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(data)
	out := &bytes.Buffer{}

	r, err := age.Decrypt(buf, rec)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(out, r); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (s *KeyService) Encrypt(userId int64, data []byte) ([]byte, error) {
	key, err := s.userPersist.GetPublicKey(userId)
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
	dBuff := bytes.NewBuffer(data)

	fmt.Fprintf(os.Stderr, "in enc %v\n", string(data))
	w, err := age.Encrypt(out, rec)
	if err != nil {
		return nil, err
	}
	defer w.Close()
	if _, err := io.Copy(w, dBuff); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

type KeySet struct {
	PrivSign   []byte
	PubSignKey []byte
	PubEncrKey string
}

// returns server's generated public encryption key, public sign key, error
func (s *KeyService) GenKeyPairForUser(userId int64, userEncrKey, userSignKey []byte) (*KeySet, error) {
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
	k := ports.Keys{
		PubEncrKey:  []byte(pubEncrStr),
		PrivEncrKey: []byte(privEncrId.String()),
		PubSignKey:  pubSign,
		PrivSignKey: privSign,
		UserEncrKey: userEncrKey,
		UserSignKey: userSignKey,
	}
	if err = s.userPersist.InsertKeySet(userId, k); err != nil {
		return nil, err
	}

	return &KeySet{
		PrivSign:   privSign,
		PubSignKey: pubSign,
		PubEncrKey: pubEncrStr,
	}, nil
}

func (s *KeyService) VerifyPassphrase(userPass string) (int64, bool, error) {
	return s.userPersist.GetUserPassphrase(userPass)
}

func (s *KeyService) VerifyRefreshToken(userId int64, tokenSig TokenSig, sigKey []byte) error {
	k, err := s.userPersist.GetKeys(userId)
	if err != nil {
		return err
	}
	token, err := DecryptAndVerify(k.PrivEncrKey, tokenSig.Token, tokenSig.Signature, sigKey)

	dbToken, err := s.userPersist.GetRefreshToken(userId)
	if err != nil {
		return err
	}

	// match
	if dbToken == string(token) {
		return nil
	}

	return fmt.Errorf("token does not match")
}

type TokenSig struct {
	Token     []byte
	Signature []byte
}

func (s *KeyService) MakeRefreshToken(userId int64, encryptKey string, signKey []byte) (*TokenSig, error) {
	t := GenerateToken(50)
	err := s.userPersist.InsertRefreshToken(userId, t)
	if err != nil {
		return nil, err
	}

	b := []byte(t)
	return encryptAndSign(b, []byte(encryptKey), signKey)
}

func (s *KeyService) MakeAccessToken(userId int64) (*TokenSig, string, error) {
	k, err := s.userPersist.GetKeys(userId)
	if err != nil {
		return nil, "", err
	}
	t := GenerateToken(50)

	b := []byte(t)
	ts, err := encryptAndSign(b, k.UserEncrKey, k.PrivSignKey)
	return ts, t, err
}

func DecryptAndVerify(cryptoKey, b, sig, sigKey []byte) ([]byte, error) {
	t, err := decrypt(cryptoKey, b)
	if err != nil {
		return nil, err
	}
	token := t
	if ok := ed25519.Verify(sigKey, token, sig); !ok {
		return nil, fmt.Errorf("token does not match")
	}
	return token, nil
}

func encryptAndSign(b []byte, cryptoKey, signKey []byte) (*TokenSig, error) {
	sig := ed25519.Sign(signKey, b)
	fmt.Fprintf(os.Stderr, "sig before enc %v\n", string(b))
	token, err := encrypt([]byte(cryptoKey), b)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(os.Stderr, "encrypted %v\n", string(token))

	return &TokenSig{
		Token:     token,
		Signature: sig,
	}, nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-!@#$%^&*()-=+[]")

func GenerateToken(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
