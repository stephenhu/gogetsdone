package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

)


func notificationHandler(w http.ResponseWriter, r *http.Request) {

	u := checkToken(r)

	if u == nil {
		w.WriteHeader(http.StatusUnauthorized)
	} else {

		vars := mux.Vars(r)

		id := vars["id"]

		if id != u.ID {
			w.WriteHeader(http.StatusUnauthorized)
		} else {

			switch r.Method {
			case http.MethodGet:
		
				contacts := getContactRequests(u.ID)
		
				j, err := json.Marshal(contacts)
	
				if err != nil {
					log.Printf("%s notificationHandler(): %s", APP_NAME, err.Error())
					w.WriteHeader(http.StatusInternalServerError)
				} else {
					w.Write(j)
				}
		
			case http.MethodPost:
			case http.MethodDelete:
			case http.MethodPut:
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}

		}

	}

} // notificationHandler
