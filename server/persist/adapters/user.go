package adapters

import (
	"fmt"
	"time"

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
	INSERT INTO key_set (
  	public_key,
  	private_key,
  	public_sign_key,
  	private_sign_key,
  	user_public_key,
  	user_sign_key,
  	user_id
	)
	VALUES (
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
	err = rows.Scan(
		&id,
		&k.Recipient,
		&k.Identity,
		&k.PubSignKey,
		&k.PrivSignKey,
		&k.UserEncrKey,
		&k.UserSignKey,
		uId,
	)
	if err != nil {
		return nil, err
	}

	return &k, nil
}

func (p UserPersist) InsertRefreshToken(userId int64, token string) (*time.Time, error) {
	sql := `
	INSERT INTO refresh_token (user_id, token)
		VALUES ($1, $2)
	RETURNING expiry;
	`
	rows, err := p.c.Query(sql, userId, token)
	if err != nil {
		return nil, err
	}
	var expiry time.Time
	err = rows.Scan(&expiry)
	return &expiry, err
}

func (p UserPersist) GetUserPassphrase(passphrase string) (int64, bool, error) {
	if passphrase == "" {
		return -1, false, fmt.Errorf("passphrase must not be empty")
	}
	sql := `
	SELECT (
		passphrase, user_id
	)
	FROM passphrase
	WHERE passphrase = $1
		AND expiry > CURRENT_TIMESTAMP
		AND is_unused = FALSE;
	`
	rows, err := p.c.Query(sql, passphrase)
	if err != nil {
		return -1, false, err
	}

	var userId int64
	var dbPassphrase string
	err = rows.Scan(&dbPassphrase, &userId)
	if err != nil {
		return -1, false, err
	}
	if passphrase == dbPassphrase {
		return userId, true, nil
	}

	return -1, false, fmt.Errorf("passphrase does not match")
}

func (p UserPersist) InsertUserPassphrase(userId int64, passphrase string) error {
	if passphrase == "" {
		return fmt.Errorf("passphrase must not be empty")
	}

	sql := `
	INSERT INTO passphrase (passphrase, user_id)
	  VALUES ($1, $2);
	`
	_, err := p.c.Execute(sql, passphrase, userId)
	return err
}
