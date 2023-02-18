package persist

import ()

type Keys struct {
	PubEncrKey    []byte
	PrivEncrKey   []byte
	PubSignKey  []byte
	PrivSignKey []byte
	UserEncrKey []byte
	UserSignKey []byte
}

func InsertKeySet(userId int64, keys Keys) error {
	panic("unimplemented")
}

func GetRefreshToken(userId int64) (string, error) {
	panic("unimplemented")
}

func GetKeys(userId int64) (Keys, error) {
	panic("unimplemented")
}

func GetPrivateKey(userId int64) ([]byte, error) {
	panic("unimplemented")
}

func GetPublicKey(userId int64) ([]byte, error) {
	panic("unimplemented")
}

func InsertRefreshToken(userId int64, token string) error {
	panic("unimplemented")
}

func GetUserPassphrase(passphrase string) (int64, bool, error) {
	panic("unimplemented")
}
