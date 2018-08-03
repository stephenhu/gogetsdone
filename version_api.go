package main

import (
	"encoding/json"
	"log"
	"net/http"

)


func versionHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
	case http.MethodGet:

		vi := VersionInfo{
			Version:	VERSION,
		}

		j, err := json.Marshal(vi)

		if err != nil {
			log.Println(err)
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

} // userHandler