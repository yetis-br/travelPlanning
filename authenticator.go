package main

//Authenticator delas with user authentication
func Authenticator(login string, password string) bool {
	return login == "admin" && password == "admin"
}
