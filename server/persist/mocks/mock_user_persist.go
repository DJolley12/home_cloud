package mocks

import (
	"github.com/DJolley12/home_cloud/server/persist/ports"
)

type UserPersist struct {
}

func NewUserPersist() UserPersist {
	return UserPersist{}
}

func (p UserPersist) InsertKeySet(userId int64, keys ports.Keys) error {
	panic("unimplemented")
}

func (p UserPersist) GetRefreshToken(userId int64) (string, error) {
	panic("unimplemented")
}

func (p UserPersist) GetKeys(userId int64) (ports.Keys, error) {
	panic("unimplemented")
}

func (p UserPersist) GetPrivateKey(userId int64) ([]byte, error) {
	panic("unimplemented")
}

func (p UserPersist) GetPublicKey(userId int64) ([]byte, error) {
	panic("unimplemented")
}

func (p UserPersist) InsertRefreshToken(userId int64, token string) error {
	panic("unimplemented")
}

func (p UserPersist) GetUserPassphrase(passphrase string) (int64, bool, error) {
	panic("unimplemented")
}
