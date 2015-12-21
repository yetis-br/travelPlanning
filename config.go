package main

import (
	"log"

	"github.com/go-ini/ini"
)

var cfg *ini.File

func init() {
	var err error
	cfg, err = ini.Load("config.ini")
	if err != nil {
		log.Println(err)
	}
}

//GetKeyValue return the value of a key from the configuration file
func GetKeyValue(section string, key string) string {
	return cfg.Section(section).Key(key).String()
}
