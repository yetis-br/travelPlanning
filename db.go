package main

import (
	"log"

	r "github.com/dancannon/gorethink"
)

//NewDBSession creates a new session to manage database data
func NewDBSession(Database string) *r.Session {
	conn, err := r.Connect(r.ConnectOpts{
		Address:  "192.168.116.128:28015",
		Database: Database,
		MaxOpen:  40,
	})
	if err != nil {
		log.Println(err)
		return nil
	}

	err = r.DBCreate(Database).Exec(conn)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Database created")
	}

	_, err = r.DB(Database).TableCreate("Trip").RunWrite(conn)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Table created")
	}

	return conn
}
