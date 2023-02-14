package persist

import (

)

func InsertKeySet(userId int64, pubK, privK, userKey []byte) error {
  panic("unimplemented")
}

func GetRefreshToken(userId int64) (string, error) {
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
