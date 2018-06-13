package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/stephenhu/gowdl"
)

const (

	CREATE_USER	= "INSERT into users(" +
		"email, salt, password) " +
		"VALUES(?, ?, ?)"

	GET_USER_BY_EMAIL = "SELECT " +
		"id, email, mobile, icon, password, salt, registered " +
		"from users " +
	  "WHERE email=?"

	)

func createUser(email string, password string) (error) {

	hash, salt, err := gowdl.GenerateHashAndSalt(password, HMAC_KEY, PEPPER,   
		HASH_LENGTH)

	if err != nil {
		
		log.Println(err)
		return err

	} else {

		_, err := data.Exec(
			CREATE_USER, email, salt, hash,
		)

		if err != nil {

			log.Println(err)
			return err

		} else {
			return nil
		}

	}

} // createUser


func getUserByEmail(email string) *User {

	row := data.QueryRow(
		GET_USER_BY_EMAIL, email,
	)

	u := User{}

	err := row.Scan(&u.ID, &u.Email, &u.Mobile, &u.Icon, &u.Password,
		&u.Salt, &u.Registered)

	if err == sql.ErrNoRows {
		log.Println("gogetsdone getUserByEmail(): ", err)
		return nil
	}

	return &u

} // getUserByEmail


func userHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
	case http.MethodPost:

		email 		:= r.FormValue("email")
		password  := r.FormValue("password")

		if email == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
		} else {

			err := createUser(email, password)

			if err != nil {
				
				log.Printf("%s userHandler(): %s", APP_NAME, err.String)
				w.WriteHeader(http.StatusForbidden)

			} else {

				
			}

		}

	case http.MethodGet:
  
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // userHandler