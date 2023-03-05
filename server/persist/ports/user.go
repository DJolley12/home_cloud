package ports

import ()

type Keys struct {
	Recipient   []byte
	Identity    []byte
	PubSignKey  []byte
	PrivSignKey []byte
	UserEncrKey []byte
	UserSignKey []byte
}

type UserPersist interface {
	GetRefreshToken(userId int64) (string, error)
	InsertRefreshToken(userId int64, token string) error
	GetKeys(userId int64) (*Keys, error)
	InsertKeySet(userId int64, keys Keys) error
	GetUserPassphrase(passphrase string) (int64, bool, error)
	InsertUserPassphrase(userId int64, passphrase string) error
}
