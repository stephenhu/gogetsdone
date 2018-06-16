package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/stephenhu/gowdl"
)

const (

	UPDATE_TOKEN = "UPDATE users" +
		"SET token='?' " +
		"WHERE id=?"

)


func checkToken(r *http.Request) *User {

	cookie, err := r.Cookie("madsportslab")

	if err != nil {
		return nil
	} else {

		u := getUserByToken(cookie.Value)

		return u

	}

} // checkToken


func authenticate(email string, password string) *User {

	if email == "" || password == "" {
		return nil
	} else {

		u := getUserByEmail(email)

		if u == nil {
			return nil
		} else {

			hash, err := gowdl.GenerateHash(password, u.Salt, HMAC_KEY,
				PEPPER, HASH_LENGTH)

			if err != nil {
				log.Printf("%s authenticate(): %s", APP_NAME, err)
				return nil
			} else {

				if hash == u.Password {
					return u
				} else {
					return nil
				}
	
			}
				

		}

	}

} // authenticate


func updateToken(u *User) (string, error) {

	if u == nil {
		return "", errors.New(fmt.Sprintf("%s updateToken(): %s", APP_NAME,
		  "cannot update token for nil user"))
	} else {

		token, err := gowdl.GenerateToken(HMAC_KEY, TOKEN_LENGTH)

		if err != nil {
			return "", err
		} else {

			_, err := data.Exec(
				UPDATE_TOKEN, token, u.ID,
			)

			if err != nil {
				return "", err
			} else {
				return token, nil
			}

		}

	}

} // updateToken


func authHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
	case http.MethodPut:

		email 		:= r.FormValue("email")
		password	:= r.FormValue("password")

		u := authenticate(email, password)

		if u == nil {
			w.WriteHeader(http.StatusUnauthorized)
		} else {

			token, err := updateToken(u)

			if err != nil {
				log.Printf("%s authHandler(): %s", APP_NAME, "unable to update token")
				w.WriteHeader(http.StatusInternalServerError)
			} else {

				cookie := http.Cookie{
          Name: GETSDONE,
          Value: token,
          Domain: "127.0.0.1",
          Path: "/",
        }
        
				http.SetCookie(w, &cookie)
				
			}

		}

	case http.MethodPost:
  case http.MethodGet:
	case http.MethodDelete:
		
		cookie := http.Cookie{
			Name:   GETSDONE,
			Value:  "",
			Domain: *domain,
			Path:   "/",
			MaxAge: -1,
		}

		http.SetCookie(w, &cookie)

	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // authHandler