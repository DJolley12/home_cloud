package services

import (
	"crypto/ed25519"
	"fmt"
	"math/rand"
	"time"

	"filippo.io/age"
	"github.com/DJolley12/home_cloud/server/persist/ports"
	encryption "github.com/DJolley12/home_cloud/shared/encryption"
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

func (s *KeyService) Encrypt(userId int64, data []byte) ([]byte, error) {
	keys, err := s.userPersist.GetKeys(userId)
	if err != nil {
		return nil, err
	}
	return encryption.Encrypt(keys.UserEncrKey, data)
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

func (s *KeyService) VerifyRefreshToken(userId int64, tokenSig encryption.TokenSig, sigKey []byte) error {
	k, err := s.userPersist.GetKeys(userId)
	if err != nil {
		return err
	}
	token, err := encryption.DecryptAndVerify(tokenSig.Token, k.Identity, tokenSig.Signature, sigKey)

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

func (s *KeyService) MakeRefreshToken(userId int64, encryptKey string, signKey []byte) (*time.Time, *encryption.TokenSig, error) {
	t := GenerateToken(50)
	expiry, err := s.userPersist.InsertRefreshToken(userId, t)
	if err != nil {
		return nil, nil, err
	}

	b := []byte(t)
	ts, err := encryption.EncryptAndSign(b, []byte(encryptKey), signKey)
	if err != nil {
		return nil, nil, err
	}
	return expiry, ts, nil
}

func (s *KeyService) MakeAccessToken(userId int64) (*encryption.TokenSig, string, error) {
	k, err := s.userPersist.GetKeys(userId)
	if err != nil {
		return nil, "", err
	}
	t := GenerateToken(50)

	b := []byte(t)
	ts, err := encryption.EncryptAndSign(b, k.UserEncrKey, k.PrivSignKey)
	return ts, t, err
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-!@#$%^&*()-=+[]")

func GenerateToken(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
