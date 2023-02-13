package main

import (
	"flag"
	"log"

	"github.com/DJolley12/home_cloud/server/persist/db"
	mg "github.com/DJolley12/home_cloud/server/persist/migrations"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "userispassword"
	password = "passwordispassword"
	dbname   = "home_cloud"
)

func main() {
	ma := flag.String("action", "", "init run or create")
	mt := flag.String("type", "", "type of migration, UP or DOWN")
	name := flag.String("name", "", "name for migration, should be short and desciptive")
	desc := flag.String("description", "", "description for the migration")
	id := flag.Int64("id", -1, "id for the migration")

	flag.Parse()

	if *ma == "run" {
		if *name == "" {
			log.Fatal("name cannot be empty")
		} else if *desc == "" {
			log.Fatal("description cannot be empty")
		}
	}

	conn := db.NewConn(port, host, user, password, dbname)
	if *ma == "init" {
		log.Println("init")
		if res, err := mg.Init(conn); err != nil {
			log.Fatal(err)
		} else {
			log.Println(res)
		}
	} else if *ma == "create" {
		res, err := mg.CreateMigration(conn)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(res)
	} else if *ma == "run" {
		mType := mg.FromString(*mt)
		res, err := mType.RunMigration(conn, *id, *desc, *name, dbname)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(res)
	}
}
