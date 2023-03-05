package db

import (
	"database/sql"
	"fmt"

	"github.com/DJolley12/home_cloud/server/persist/ports"
	_ "github.com/lib/pq"
)

type Conn struct {
	config ports.DBConfig
}

func NewConn(c ports.DBConfig) Conn {
	return Conn{
		c,
	}
}

func (c *Conn) Execute(q string, params ...any) (sql.Result, error) {
	db, err := sql.Open("postgres", c.getInfo())
	if err != nil {
		return nil, err
	}
	defer db.Close()
	if len(params) > 0 {
		return db.Exec(q, params...)
	}
	return db.Exec(q)
}

func (c *Conn) Query(q string, params ...any) (*sql.Rows, error) {
	db, err := sql.Open("postgres", c.getInfo())
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return db.Query(q, params...)
}

func (c *Conn) getInfo() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.config.Host, c.config.Port, c.config.User, c.config.Password, c.config.DBName)
}
