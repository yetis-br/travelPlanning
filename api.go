package main

import (
	"log"
	"net/http"
	"time"

	"github.com/StephanDollberg/go-json-rest-middleware-jwt"
	"github.com/ant0ine/go-json-rest/rest"
)

var (
	//JWT middleware to manage token authentication
	JWT *jwt.JWTMiddleware
	//SecretKey default secret key to create tokens
	SecretKey = []byte("secret key")
	//Realm default to use in request and response header
	Realm = "jwt auth"
)

func main() {
	JWT = &jwt.JWTMiddleware{
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
		IfTrue: JWT,
	})

	router := NewRouter()
	api.SetApp(router)

	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
