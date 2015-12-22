package main

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
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
	log.Infof("Starting in %s mode", GetKeyValue("server", "mode"))
	log.Infof("Listening on port: %s", GetKeyValue("server", "port"))

	jwt := &jwt.JWTMiddleware{
		Key:        SecretKey,
		Realm:      Realm,
		Timeout:    time.Hour,
		MaxRefresh: time.Hour * 24,
		Authenticator: func(userId string, password string) bool {
			return Authenticator(userId, password)
		}}

	api := rest.NewApi()
	api.Use(&AccessLogTPApiMiddleware{
		Format: DefaultLogFormat,
	})
	api.Use(rest.DefaultCommonStack...)
	api.Use(&rest.ContentTypeCheckerMiddleware{})
	api.Use(&rest.JsonIndentMiddleware{})
	api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			return CheckCondition(request)
		},
		IfTrue: jwt,
	})

	router := NewRouter(jwt, NewDBSession("travelPlanning"))
	api.SetApp(router)

	elapsed := time.Since(start)
	log.Infof("Started in %fs", elapsed.Seconds())

	log.Fatalln(http.ListenAndServe(":"+GetKeyValue("server", "port"), api.MakeHandler()))
}
