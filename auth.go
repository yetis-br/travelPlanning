package main

import "github.com/ant0ine/go-json-rest/rest"

//Authenticator manage user authentication
func Authenticator(login string, password string) bool {
	return login == "admin" && password == "admin"
}

//CheckCondition manage routes authentication to define permission levels
func CheckCondition(request *rest.Request) bool {
	return request.URL.Path != "/login"
}
