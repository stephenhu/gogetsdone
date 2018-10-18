package main

import (
	//"encoding/json"
	"log"
	"net/http"

	"github.com/madsportslab/glbs"

)

const (

	UPDATE_USER_ICON = "UPDATE users " +
		"SET icon=? " +
		"WHERE id=?"

)


func updateUserIcon(id string, key string) bool {

	_, err := data.Exec(
		UPDATE_USER_ICON, key, id,
	)

	if err != nil {
		log.Printf("%s updateUserIcon(): %s", APP_NAME, err.Error())
		return false
	} else {
		return true
	}

} // updateUserIcon


func iconHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
	case http.MethodPost:
		
		u := checkToken(r)

		if u == nil {
			w.WriteHeader(http.StatusUnauthorized)
		} else {

			file, header, err := r.FormFile("icon")

			if err != nil {
				log.Println(err)
			} else {

				defer file.Close()

				if header.Size > ICON_MAX_BYTES {
				
					log.Println("Icon exceeds maximum supported size.")
					w.WriteHeader(http.StatusRequestEntityTooLarge)
				
				} else {

					glbs.SetNamespace("icons")

					k := glbs.Put(file)
	
					if k == nil {
						w.WriteHeader(http.StatusInternalServerError)
					} else {

						if !updateUserIcon(u.ID, *glbs.GetPath(*k)) {
							w.WriteHeader(http.StatusInternalServerError)
						}


					}

				}
				
			}


		}


	case http.MethodGet:
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // iconHandler
