package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/StephanDollberg/go-json-rest-middleware-jwt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var Response *httptest.ResponseRecorder
var Request *http.Request
var loginCredentials = map[string]string{"username": "admin", "password": "admin"}
var tst *testing.T

func TestTravelPlanning(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TravelPlanning Suite")
}

type resultToken struct {
	TokenString string `json:"token"`
}

func Login() string {
	APIRequest("/login", "POST", loginCredentials, "")

	token := resultToken{}
	json.Unmarshal(Response.Body.Bytes(), &token)

	return token.TokenString
}

func APIRequest(url string, method string, model interface{}, token string) {
	jwtMiddleware := &jwt.JWTMiddleware{
		Key:        SecretKey,
		Realm:      Realm,
		Timeout:    time.Hour,
		MaxRefresh: time.Hour * 24,
		Authenticator: func(userId string, password string) bool {
			return Authenticator(userId, password)
		},
	}

	api := rest.NewApi()
	api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			return CheckCondition(request)
		},
		IfTrue: jwtMiddleware,
	})

	api.SetApp(NewRouter(jwtMiddleware))

	Request = test.MakeSimpleRequest(method, url, model)
	if token != "" {
		Request.Header.Set("Authorization", "Bearer "+token)
	}
	recorded := test.RunRequest(tst, api.MakeHandler(), Request)
	Response = recorded.Recorder
}
