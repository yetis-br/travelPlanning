package main

import (
	"log"
	"net/http"
	"time"

	"github.com/StephanDollberg/go-json-rest-middleware-jwt"
	"github.com/ant0ine/go-json-rest/rest"
)

var (
	//SecretKey default secret key to create tokens
	SecretKey = []byte("secret key")
	//Realm default to use in request and response header
	Realm = "jwt auth"
)

func main() {
	start := time.Now()
	log.SetPrefix("[Travel Planning API] ")
	log.Printf("Starting in %s mode", GetKeyValue("server", "mode"))
	log.Printf("Listening on port: %s", GetKeyValue("server", "port"))

	jwt := &jwt.JWTMiddleware{
		Key:        SecretKey,
		Realm:      Realm,
		Timeout:    time.Hour,
		MaxRefresh: time.Hour * 24,
		Authenticator: func(userId string, password string) bool {
			return Authenticator(userId, password)
		}}

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			return CheckCondition(request)
		},
		IfTrue: jwt,
	})

	router := NewRouter(jwt, NewDBSession("travelPlanning"))
	api.SetApp(router)

	elapsed := time.Since(start)
	log.Printf("Started in %fs", elapsed.Seconds())

	log.Fatal(http.ListenAndServe(":"+GetKeyValue("server", "port"), api.MakeHandler()))
}
