package payload

import (
	// "bytes"
	"crypto/ed25519"
	"encoding/json"
	"os"
	"path"

	// "fmt"
	// "io"
	// "math/rand"
	// "os"

	"filippo.io/age"
)

type KeyService struct {
	baseDir string
}

type KeySet struct {
	PrivSignKey []byte `json:privSignKey`
	PubSignKey  []byte `json:pubSignKey`
	Identity    string `json:identity`
	Recipient   string `json:reciptient`
}

// returns server's generated public encryption key, public sign key, error
func (s *KeyService) GenKeyPair() (*KeySet, error) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		return nil, err
	}
	recipient := identity.Recipient()

	pubSign, privSign, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	k := KeySet{
		Recipient:   recipient.String(),
		Identity:    identity.String(),
		PubSignKey:  pubSign,
		PrivSignKey: privSign,
	}
	if err = s.saveClientKeySet(k); err != nil {
		return nil, err
	}

	return &KeySet{
		PrivSignKey: privSign,
		PubSignKey:  pubSign,
		Recipient:   recipient.String(),
	}, nil
}

func (s *KeyService) saveClientKeySet(k KeySet) error {
	json, err := json.Marshal(k)
	if err != nil {
		return err
	}

	f, err := os.Create(path.Join(s.baseDir, "client-keys.json"))
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
