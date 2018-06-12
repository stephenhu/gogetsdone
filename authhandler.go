package main

import (
	"net/http"

)

const (

	SIGNUP_USER		= "INSERT into users(" +
		"email, password" +
		") VALUES(" + 
		"?, ?);"

)



func authHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
	case http.MethodPost:

  case http.MethodGet:
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // authHandler