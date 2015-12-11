package main_test

import (
	"encoding/json"
	"log"
	"net/http/httptest"
	"time"

	"github.com/StephanDollberg/go-json-rest-middleware-jwt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/yetis-br/travelPlanning"

	"testing"
)

var Response *httptest.ResponseRecorder
var Request *rest.Request
var tst *testing.T

func TestTravelPlanning(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TravelPlanning Suite")
}

type resultToken struct {
	TokenString string `json:"token"`
}

func Login(loginCreds map[string]string) string {
	// the middleware to test
	jwtMiddleware := &jwt.JWTMiddleware{
		Key:        main.SecretKey,
		Realm:      main.Realm,
		Timeout:    time.Hour,
		MaxRefresh: time.Hour * 24,
		Authenticator: func(userId string, password string) bool {
			return main.Authenticator(userId, password)
		},
	}
	// api for login purpose
	api := rest.NewApi()
	api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			return request.URL.Path != "/login"
		},
		IfTrue: jwtMiddleware,
	})
	api.SetApp(rest.AppSimple(jwtMiddleware.LoginHandler))
	Request := test.MakeSimpleRequest("POST", "/login", loginCreds)
	recorded := test.RunRequest(tst, api.MakeHandler(), Request)
	Response = recorded.Recorder

	token := resultToken{}
	json.Unmarshal(recorded.Recorder.Body.Bytes(), &token)

	return token.TokenString
}

func APIRequest(url string, handlerFunc rest.HandlerFunc, method string, json interface{}, token string) {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(main.JWT)

	var restRoute *rest.Route
	switch method {
	case "GET":
		restRoute = rest.Get(url, handlerFunc)
	case "POST":
		restRoute = rest.Post(url, handlerFunc)
	}

	router, err := rest.MakeRouter(
		restRoute,
	)
	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)

	r := test.MakeSimpleRequest(method, url, json)
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}

	recorded := test.RunRequest(tst, api.MakeHandler(), r)
	Response = recorded.Recorder
}
