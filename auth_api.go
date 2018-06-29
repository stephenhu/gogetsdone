package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/stephenhu/gowdl"
)

const (

	UPDATE_TOKEN = "UPDATE users " +
		"SET token=? " +
		"WHERE id=?"

	DELETE_TOKEN = "UPDATE users " +
		"SET token='' " +
		"WHERE id=?"

)


func encryptCookieData(id string, token string, icon string) (string, error) {

	block, err := aes.NewCipher([]byte(BLOCK_KEY))

	if err != nil {
		return "", err
	} else {

		cfb := cipher.NewCFBEncrypter(block, []byte(IV))

		c := CookieData{
			ID: id,
			Token: token,
		}

		j, err := json.Marshal(c)

		if err != nil {
			log.Printf("%s encryptCookieData(): %s", APP_NAME, err.Error())
			return "", err
		} else {

			ciphertext := make([]byte, len(j))

			cfb.XORKeyStream(ciphertext, j)

			return hex.EncodeToString(ciphertext), nil

		}
	  
	}

} // encryptCookieData


func decryptCookieData(ciphertext string) *CookieData {

	block, err := aes.NewCipher([]byte(BLOCK_KEY))

	if err != nil {
		
		log.Printf("%s decryptCookieData(): %s", APP_NAME, err.Error())
		return nil

	} else {

		cfb := cipher.NewCFBDecrypter(block, []byte(IV))

		decoded, err := hex.DecodeString(ciphertext)

		if err != nil {
			
			log.Printf("%s decryptCookieData(): %s", APP_NAME, err.Error())
			return nil

		} else {

			j := make([]byte, len(decoded))

			cfb.XORKeyStream(j, decoded)

			c := CookieData{}
	
			err := json.Unmarshal(j, &c)
	
			if err != nil {
				
				log.Printf("%s decryptCookieData(): %s", APP_NAME, err.Error())
				return nil

			} else {
				return &c
			}

		}


	}
} // decryptCookieData


func checkToken(r *http.Request) *User {

	cookie, err := r.Cookie(GETSDONE)

	if err != nil {
		return nil
	} else {

		cookieData := decryptCookieData(cookie.Value)

		log.Println(cookieData)
		if cookieData != nil {

			u := getUserByToken(cookieData.Token)

			return u

		} else {
			return nil
		}
		
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


func deleteToken(u *User) (error) {

	if u == nil {
		return errors.New(fmt.Sprintf("%s deleteToken(): %s", APP_NAME,
		  "cannot delete token for nil user"))
	} else {

		_, err := data.Exec(
			DELETE_TOKEN, u.ID,
		)

		if err != nil {
			return err
		} else {
			return nil
		}

	}

} // deleteToken


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
				log.Printf("%s authHandler(): %s", APP_NAME, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
			} else {

				encryptedData, err := encryptCookieData(u.ID, token, u.Icon.String)

				if err != nil {

					log.Println(err)
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
				
			}

		}

	case http.MethodPost:
  case http.MethodGet:
	case http.MethodDelete:
		
		cookie := &http.Cookie{
			Name:   GETSDONE,
			Value:  "",
			Domain: *domain,
			Path:   "/",
			MaxAge: -1,
		}

		http.SetCookie(w, cookie)

		u := checkToken(r)

		if u != nil {

			err := deleteToken(u)

			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}

		}

	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // authHandler
