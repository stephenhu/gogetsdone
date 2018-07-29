package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/stephenhu/gowdl"
)

const (

	CREATE_USER	= "INSERT into users(" +
		"email, salt, password, name, token) " +
		"VALUES(?, ?, ?, ?, ?)"

	GET_USER_BY_EMAIL = "SELECT " +
		"id, email, name, mobile, icon, password, salt, registered, token " +
		"from users " +
	  "WHERE email=?"

	GET_USER_BY_NAME = "SELECT " +
		"id, email, name, icon " +
		"from users " +
		"WHERE name=?"
		
	GET_USER_BY_TOKEN = "SELECT " +
		"users.id, users.email, users.name, users.mobile, " +
		"users.icon, users.password, users.salt, users.registered, " + 
		"users.token, ranks.rank " +
		"from users, ranks " +
		"WHERE users.token=? and users.rank_id=ranks.id"
		
)


func createUser(email string, password string) (error) {

	hash, salt, err := gowdl.GenerateHashAndSalt(password, HMAC_KEY, PEPPER,   
		HASH_LENGTH)

	if err != nil {
		
		log.Println(err)
		return err

	} else {

		name, err := gowdl.ParseNameFromEmail(email)

		if err != nil {
			log.Printf("%s createUser(): %s", APP_NAME, err.Error())
			return err
		} else {

			token, err := gowdl.GenerateToken(HMAC_KEY, TOKEN_LENGTH)

			if err != nil {
				return err
			} else {

				_, err := data.Exec(
					CREATE_USER, email, salt, hash, name, token,
				)
		
				if err != nil {
		
					log.Println(err)
					return err
		
				} else {
					return nil
				}

			}
	
		}

	}

} // createUser


func getUserByName(name string) *UserInfo {

	row := data.QueryRow(
		GET_USER_BY_NAME, name,
	)

	u := UserInfo{}

	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Icon)

	if err != nil || err == sql.ErrNoRows {
		log.Println("gogetsdone getUserByName(): ", err)
		return nil
	}

	return &u

} // getUserByName


func getUserByEmail(email string) *User {

	row := data.QueryRow(
		GET_USER_BY_EMAIL, email,
	)

	u := User{}

	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Mobile, &u.Icon, &u.Password,
		&u.Salt, &u.Registered, &u.Token)

	if err != nil || err == sql.ErrNoRows {
		log.Println("gogetsdone getUserByEmail(): ", err)
		return nil
	}

	return &u

} // getUserByEmail


func getUserByToken(token string) *User {

	row := data.QueryRow(
		GET_USER_BY_TOKEN, token,
	)

	u := User{}

	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.Mobile, &u.Icon, &u.Password,
		&u.Salt, &u.Registered, &u.Token, &u.RankName)

	if err != nil || err == sql.ErrNoRows {
		log.Println("gogetsdone getUserByToken(): ", err)
		return nil
	}

	return &u

} // getUserByToken


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
				
				log.Printf("%s userHandler(): %s", APP_NAME, err.Error())
				w.WriteHeader(http.StatusForbidden)

			} else {

				u := getUserByEmail(email)

				if u == nil {
					w.WriteHeader(http.StatusNotFound)
				} else {

					if u.Token.Valid {

						encryptedData, err := encryptCookieData(u.ID, u.Token.String,
							u.Icon.String)

						if err != nil {
							log.Printf("%s userHandler(): %s", APP_NAME, err.Error())
							w.WriteHeader(http.StatusInternalServerError)
						} else {

							cookie := &http.Cookie{
								Name: GETSDONE,
								Value: encryptedData,
								Domain: *domain,
								Path: "/",
							}

							http.SetCookie(w, cookie)

						}						

					} else {
						log.Printf("%s userHandler(): %s", APP_NAME, "no token")
					}

								
				}

			}

		}

	case http.MethodGet:

		u := checkToken(r)

		if u != nil {

			info := UserInfo{
				ID:				u.ID,
				Name:  		u.Name,
				RankName: u.RankName,
			}

			j, err := json.Marshal(info)

			if err != nil {
				
				log.Printf("%s userHandler(): %s", APP_NAME, err.Error())
				w.WriteHeader(http.StatusInternalServerError)

			} else {
				w.Write(j)
			}

		} else {
			
			log.Printf("%s userHandler(): %s", APP_NAME, "user not authenticated")
			w.WriteHeader(http.StatusUnauthorized)

		}
  
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // userHandler