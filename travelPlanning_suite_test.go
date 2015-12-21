package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"testing"

	"github.com/StephanDollberg/go-json-rest-middleware-jwt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
	r "github.com/dancannon/gorethink"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const database = "travelPlanningTest"

var (
	Response         *httptest.ResponseRecorder
	Request          *http.Request
	testConn         *r.Session
	tst              *testing.T
	loginCredentials = map[string]string{"username": "admin", "password": "admin"}
)

func TestTravelPlanning(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TravelPlanning Suite")
}

var _ = BeforeSuite(func() {
	testConn = NewDBSession(database)
})

var _ = AfterSuite(func() {
	resp, err := r.DBDrop(database).RunWrite(testConn)
	if err != nil {
		log.Print(err)
	}

	log.Printf("%d DB dropped, %d tables dropped", resp.DBsDropped, resp.TablesDropped)
	testConn.Close()
})

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

	api.SetApp(NewRouter(jwtMiddleware, testConn))

	Request = test.MakeSimpleRequest(method, url, model)
	if token != "" {
		Request.Header.Set("Authorization", "Bearer "+token)
	}
	recorded := test.RunRequest(tst, api.MakeHandler(), Request)
	Response = recorded.Recorder
}
