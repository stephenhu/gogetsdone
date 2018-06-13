package main

import (
	"net/http"

)
func taskHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
  case http.MethodGet:
  case http.MethodPost:
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // taskHandler