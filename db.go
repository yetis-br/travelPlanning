package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	r "github.com/dancannon/gorethink"
)

//NewDBSession creates a new session to manage database data
func NewDBSession(database string) *r.Session {
	conn, err := r.Connect(r.ConnectOpts{
		Address:  GetKeyValue("database", "address"),
		Database: database,
		MaxOpen:  40,
	})
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	log.Infof("Connected to RethinkDB %s", GetKeyValue("database", "address"))

	status, err := createDB(database, conn)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Info(status)
	}

	tables := []string{"Trip"}
	total, err := createTables(database, tables, conn)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Infof("%d new table(s) created", total)
	}

	return conn
}

func createDB(database string, conn *r.Session) (string, error) {
	var response interface{}
	res, _ := r.DBList().Contains(database).Run(conn)
	res.One(&response)
	if response == false {
		return fmt.Sprintf("New database %s created", database), r.DBCreate(database).Exec(conn)
	}
	return fmt.Sprintf("Database %s already exists", database), nil
}

func createTables(database string, tables []string, conn *r.Session) (int, error) {
	var (
		response interface{}
		total    int
	)
	for _, table := range tables {
		res, _ := r.DB(database).TableList().Contains(table).Run(conn)
		res.One(&response)
		if response == false {
			_, err := r.DB(database).TableCreate(table).RunWrite(conn)
			if err != nil {
				return 0, err
			}
			total++
		}
	}
	return total, nil
}
