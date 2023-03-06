package payload

import (
	// "bytes"
	"crypto/ed25519"
	"encoding/json"
	"os"
	"path"
	"time"

	// "fmt"
	// "io"
	// "math/rand"
	// "os"

	"filippo.io/age"
)

type KeyService struct {
	baseDir string
}

type UserKeySet struct {
	PrivSignKey []byte `json:"privSignKey"`
	PubSignKey  []byte `json:"pubSignKey"`
	Identity    string `json:"identity"`
	Recipient   string `json:"recipient"`
}

type ServerKeySet struct {
	UserId     int64  `json:"userId"`
	PubSignKey []byte `json:"pubSignKey"`
	Recipient  string `json:"recipient"`
}

type Token struct {
	RefreshToken []byte `json:"refresh-token"`
	Expiry       time.Time `json:"expiry"`
}

// returns server's generated public encryption key, public sign key, error
func (s *KeyService) GenKeyPair() (*UserKeySet, error) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		return nil, err
	}
	recipient := identity.Recipient()

	pubSign, privSign, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	k := UserKeySet{
		Recipient:   recipient.String(),
		Identity:    identity.String(),
		PubSignKey:  pubSign,
		PrivSignKey: privSign,
	}
	if err = s.saveClientKeySet(k); err != nil {
		return nil, err
	}

	return &UserKeySet{
		PrivSignKey: privSign,
		PubSignKey:  pubSign,
		Recipient:   recipient.String(),
	}, nil
}

func (s *KeyService) saveClientKeySet(k UserKeySet) error {
	return save(k, path.Join(s.baseDir, "client-keys.json"))
}

func (s *KeyService) SaveServerKeys(k ServerKeySet) error {
	return save(k, path.Join(s.baseDir, "server-keys.json"))
}

func (s *KeyService) SaveToken(token Token) error {
	return save(token, path.Join(s.baseDir, "refresh-token.json"))
}

func save(data any, path string) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = f.Write(json)
	if err != nil {
		return err
	}
	return nil
}
