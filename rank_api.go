package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

)


const (

	GET_RANKS = "SELECT id, rank, count " +
	"FROM ranks "

)


func getRanks() []Rank {

	ranks := []Rank{}

	rows, err := data.Query(
		GET_RANKS,
	)

	if err != nil || err == sql.ErrNoRows {
		log.Printf("%s getRanks(): %s", APP_NAME, err.Error())
	} else {

		for rows.Next() {

			r := Rank{}

			err := rows.Scan(&r.ID, &r.Name, &r.Count)
			
			if err != nil {
				log.Printf("%s getRanks(): %s", APP_NAME, err.Error())
			} else {
				ranks = append(ranks, r)
			}

		}

	}

	return ranks

} // getRanks


func rankHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
	case http.MethodGet:

		ranks := getRanks()

		j, err := json.Marshal(ranks)

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

} // rankHandler
