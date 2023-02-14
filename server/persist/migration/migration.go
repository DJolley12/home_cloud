package migrations

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/DJolley12/home_cloud/server/persist/db"
	// randgen "github.com/DJolley12/home_cloud/server/rand_gen"
)

type Migration struct {
	Id             int64
	Name           string
	Description    string
	MigrationState MigrationType
}

type MigrationType int

const (
	Down MigrationType = iota
	Up
)

const (
	base   = "server/persist/migration/sql"
	check  = "core/check_schema.sql"
	create = "core/create_migration.sql"
	get    = "core/get_migration.sql"
	down   = "core/down_migration.sql"
	up     = "core/up_migration.sql"
	start  = "core/start_migration.sql"
)

func FromString(mType string) MigrationType {
	switch mType {
	case "DOWN":
		return Down
	case "UP":
		return Up
	default:
		panic(fmt.Sprintf("unrecognized migration type %#v", mType))
	}
}
func (t *MigrationType) String() string {
	switch *t {
	case Down:
		return "DOWN"
	case Up:
		return "UP"
	default:
		panic(fmt.Sprintf("unrecognized migrationType: %#v", t))
	}
}

func Init(c db.Conn) (string, error) {
	sql, err := getSql(start)
	log.Println(sql)
	if err != nil {
		return "", err
	}
	_, err = c.Execute(sql)
	if err != nil {
		return "", err
	}
	return "sucessfully initialized migrations", nil
}

func CreateMigration(c db.Conn) (string, error) {
	sql, err := getSql(create)
	if err != nil {
		return "", err
	}
	res, err := c.Query(sql)
	if err != nil {
		return "", err
	}

	var id int64
	_ = res.Next()
	if err := res.Scan(&id); err != nil {
		return "", err
	}
	u, d := Up, Down
	uName, dName := u.fmt(id), d.fmt(id)

	dir := fmt.Sprintf("%v/user_migrations/%v", base, id)
	err = os.Mkdir(dir, 0755)
	if err != nil {
		return "", err
	}

	f, err := os.Create(filepath.Join(dir, uName))
	if err != nil {
		return "", err
	}
	f.Close()

	f, err = os.Create(filepath.Join(dir, dName))
	if err != nil {
		return "", err
	}
	f.Close()

	return fmt.Sprintf("created migration id: %v", id), nil
}

func (t *MigrationType) RunMigration(c db.Conn, id int64, name, desc, dbName string) (string, error) {
	if id < 0 {
		return "", fmt.Errorf("id cannot be empty")
	}
	fName := t.fmt(id)
	sql, err := getSql(fmt.Sprintf("/user_migrations/%v/%v", id, fName))
	if err != nil {
		return "", err
	}
	if strings.Trim(sql, " \r\n") == "" {
		return "", fmt.Errorf("migration file %v is empty", fName)
	}

	uSql, err := getSql(up)
	if err != nil {
		return "", err
	}
	dSql, err := getSql(down)

	validateId(c, id)

	spl := strings.Split(sql, ";")
	for _, sq := range spl {
		_, err = c.Execute(fmt.Sprintf("%v;", sq))
		if err != nil {
			return "", err
		}
	}

	if *t == Up {
		_, err = c.Execute(uSql, id, name, desc, t.String())
		if err != nil {
			return "", err
		}
	} else if *t == Down {
		_, err = c.Execute(dSql)
	}

	return "successfully ran migration %v", nil
}

func validateId(c db.Conn, id int64) error {
	g, err := getSql(get)
	if err != nil {
		return err
	}
	res, err := c.Query(g, id)
	m := Migration{}
	for {
		var ms string
		err = res.Scan(&m.Id, &m.Name, &m.Description, &ms)
		if err != nil {
			return err
		}
		m.MigrationState = FromString(ms)
		if !res.Next() {
			break
		} else {
			return fmt.Errorf("too many rows returned for id %v", id)
		}
	}

	if m.Id != id {
		panic("wtf")
	}
	return nil
}

func (t *MigrationType) fmt(id int64) string {
	switch *t {
	case Down:
		return fmt.Sprintf("DOWN_%v.sql", id)
	case Up:
		return fmt.Sprintf("UP_%v.sql", id)
	default:
		panic(fmt.Sprintf("unrecognized migration type: %#v", *t))
	}
}

func getSql(name string) (string, error) {
	// curr, err := os.Getwd()
	// if err != nil {
	// 	return "", err
	// }
	fmt.Fprintf(os.Stderr, "name: %v\n", name)
	b, err := os.ReadFile(filepath.Join(base, name))
	if err != nil {
		return "", err
	}
	return string(b), nil
}
