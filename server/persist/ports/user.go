package ports

import ()

type Keys struct {
	PubEncrKey  []byte
	PrivEncrKey []byte
	PubSignKey  []byte
	PrivSignKey []byte
	UserEncrKey []byte
	UserSignKey []byte
}

type UserPersist interface {
	InsertKeySet(userId int64, keys Keys) error
	GetRefreshToken(userId int64) (string, error)
	GetKeys(userId int64) (Keys, error)
	GetPrivateKey(userId int64) ([]byte, error)
	GetPublicKey(userId int64) ([]byte, error)
	InsertRefreshToken(userId int64, token string) error
	GetUserPassphrase(passphrase string) (int64, bool, error)
}
