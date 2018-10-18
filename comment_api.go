package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"

)


const (

	GET_COMMENTS_BY_TASK	= "SELECT comments.id, comments.comment, " +
		"comments.created, comments.task_id, users.name, users.icon " +
		"FROM comments, users " +
		"WHERE comments.task_id=? and comments.user_id=users.id " +
		"ORDER BY comments.created DESC"

	CREATE_COMMENT = "INSERT into comments(user_id, task_id, comment) " +
	  "VALUES(?, ?, ?)"
		
)


func getCommentsByTask(id string) []Comment {

	comments := []Comment{}
	
	rows, err := data.Query(
		GET_COMMENTS_BY_TASK, id,
	)

	defer rows.Close()

	if err != nil || err == sql.ErrNoRows {
		
		log.Printf("%s getCommentsByTask(): %s", APP_NAME, err.Error())
		return comments

	} else {

		for rows.Next() {

			c := Comment{}

			err := rows.Scan(&c.ID, &c.Comment, &c.Created, &c.TaskID, &c.UserName,
			  &c.UserIcon)
				
			if err != nil || err == sql.ErrNoRows {

				log.Printf("%s getCommentsByTask(): %s", APP_NAME, err.Error())
				return comments

			} else {
				comments = append(comments, c)
			}

		}

		return comments

	}

} // getCommentsByTask


func addComment(id string, tid string, comment string) error {

	_, err := data.Exec(
		CREATE_COMMENT, id, tid, comment,
	)

	if err != nil {

		log.Printf("%s addComment(): %s", APP_NAME, err.Error())
		return err

	} else {
		return nil
	}

} // addComment


func commentHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id 		:= vars["id"]
	tid 	:= vars["tid"]

  switch r.Method {
	case http.MethodPost:

		u := checkToken(r)

		if u == nil {
			w.WriteHeader(http.StatusUnauthorized)
		} else {

			if u.ID == id {

				comment := r.FormValue("comment")

				if comment == "" {
	
					log.Printf("%s Post Comment commentHandler(): %s", APP_NAME, "Comment cannot be empty string")
					w.WriteHeader(http.StatusBadRequest)
	
				} else {
		
					err := addComment(u.ID, tid, comment)
	
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
					}
	
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

} // commentHandler
