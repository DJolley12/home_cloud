package adapters

import (
	"github.com/DJolley12/home_cloud/server/persist/db"
	"github.com/DJolley12/home_cloud/server/persist/ports"
)

type UserPersist struct {
	c db.Conn
}

func NewUserPersist(c ports.DBConfig) UserPersist {
	return UserPersist{
		c: db.NewConn(c),
	}
}

func (p UserPersist) InsertKeySet(userId int64, keys ports.Keys) error {
	sql := `
	INSERT INTO key_set(
  	public_key,
  	private_key,
  	public_sign_key,
  	private_sign_key,
  	user_public_key,
  	user_sign_key,
  	user_id
	)
	VALUES(
  	$1,
  	$2,
  	$3,
  	$4
  	$5,
  	$6,
  	$7
	);
	`
	_, err := p.c.Execute(sql,
		keys.Recipient,
		keys.Identity,
		keys.PubSignKey,
		keys.PrivSignKey,
		keys.UserEncrKey,
		keys.UserSignKey,
		userId,
	)

	return err
}

func (p UserPersist) GetRefreshToken(userId int64) (string, error) {
	panic("unimplemented")
}

func (p UserPersist) GetKeys(userId int64) (*ports.Keys, error) {
	sql := `
	SELECT (
		id,
		public_key,
		private_key,
		public_sign_key,
		private_sign_key,
		user_public_key,
		user_sign_key,
		user_id
	)
	FROM key_set
	WHERE user_id = $1;
	`
	rows, err := p.c.Query(sql, userId)
	if err != nil {
		return nil, err
	}

	var id int
	var uId int
	k := ports.Keys{}
	rows.Scan(
		&id,
		&k.Recipient,
		&k.Identity,
		&k.PubSignKey,
		&k.PrivSignKey,
		&k.UserEncrKey,
		&k.UserSignKey,
		uId,
	)

	return &k, nil
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
