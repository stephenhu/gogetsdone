package main

import (
	"log"
	"net/http"

	"github.com/stephenhu/gowdl"
)

const (

	CREATE_USER	= "INSERT into users(" +
		"email, salt, password) " +
		"VALUES(?, ?, ?);"

	)

func createUser(email string, hash string, salt string) (error) {


	_, err := data.Exec(
		CREATE_USER, email, salt, hash,
	)

	if err != nil {
		log.Println(err)
		return err
	} else {
		return nil
	}

} // createUser

func userHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
	case http.MethodPost:

		email 		:= r.FormValue("email")
		password  := r.FormValue("password")

		if email == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
		} else {

			hash, salt, err := gowdl.GenerateHashAndSalt(password, HMAC_KEY, PEPPER,   
				HASH_LENGTH)

			if err != nil {
				
				log.Println(err)
				w.WriteHeader(http.StatusConflict)

			} else {

				err := createUser(email, hash, salt)

				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusForbidden)
				} else {
	
				  
				}
	
			}

		}

	case http.MethodGet:
  
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // userHandler