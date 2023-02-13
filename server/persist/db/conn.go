package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Conn struct {
	host     string
	port     int
	user     string
	dbname   string
	password string
}

func NewConn(port int, host, user, password, dbname string) Conn {
	return Conn{
		host:     host,
		port:     port,
		user:     user,
		dbname:   dbname,
		password: password,
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
		c.host, c.port, c.user, c.password, c.dbname)
}
