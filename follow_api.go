package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

)


const (

	REMOVE_FOLLOW	= "DELETE from follows " +
		"WHERE user_id=? and follow_id=?"
		
)


func removeFollow(uid string, fid string) (error) {

	_, err := data.Exec(
		REMOVE_FOLLOW, uid, fid,
	)

	if err != nil {
		log.Printf("%s removeFollow(): %s", APP_NAME, err.Error())
		return err
	} else {
		return nil
	}

} // removeFollow


func followHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id := vars["id"]

  switch r.Method {
	case http.MethodPost:

		u := checkToken(r)

		if u == nil {
			w.WriteHeader(http.StatusUnauthorized)
		} else {

			if id == u.ID {

				fid := r.FormValue("fid")

				if fid != "" {

					err := removeFollow(u.ID, fid)

					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
					}

				} else {
					w.WriteHeader(http.StatusBadRequest)
				}

			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}

		}
		
	case http.MethodGet:
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // followHandler
