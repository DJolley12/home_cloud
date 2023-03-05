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

func NewKeyService(persist ports.UserPersist) (*KeyService, error) {
	if persist == nil {
		return nil, fmt.Errorf("persist cannot be nil")
	}
	return &KeyService{
		userPersist: persist,
	}, nil
}

func decrypt(privK, data []byte) ([]byte, error) {
	identity, err := age.ParseX25519Identity(string(privK))
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(data)

	r, err := age.Decrypt(buf, identity)
	if err != nil {
		return nil, err
	}

	out := &bytes.Buffer{}
	if _, err := io.Copy(out, r); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func (s *KeyService) Encrypt(userId int64, data []byte) ([]byte, error) {
	keys, err := s.userPersist.GetKeys(userId)
	if err != nil {
		return nil, err
	}
	return encrypt(keys.UserEncrKey, data)
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
	Recipient  string
}

// returns sign keys, encryption key
func (s *KeyService) GenKeyPairForUser(userId int64, userEncrKey, userSignKey []byte) (*KeySet, error) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		return nil, err
	}
	recipient := identity.Recipient()

	pubSign, privSign, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	k := ports.Keys{
		Recipient:   []byte(recipient.String()),
		Identity:    []byte(identity.String()),
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
		Recipient:  recipient.String(),
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
	token, err := DecryptAndVerify(tokenSig.Token, k.Identity, tokenSig.Signature, sigKey)

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

func DecryptAndVerify(data, cryptoKey, sig, sigKey []byte) ([]byte, error) {
	t, err := decrypt(cryptoKey, data)
	if err != nil {
		return nil, err
	}
	token := t
	if ok := ed25519.Verify(sigKey, token, sig); !ok {
		return nil, fmt.Errorf("token does not match")
	}
	return token, nil
}

func encryptAndSign(data []byte, encryptionKey, priv []byte) (*TokenSig, error) {
	out := &bytes.Buffer{}
	sig := ed25519.Sign(priv, data)
	signer := bytes.NewBuffer(data)

	rec, err := age.ParseX25519Recipient(string(encryptionKey))
	w, err := age.Encrypt(out, rec)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(w, signer); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	return &TokenSig{
		Token:     out.Bytes(),
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
