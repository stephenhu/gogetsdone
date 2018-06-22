package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/stephenhu/gowdl"

)

const (

	CREATE_TASK = "INSERT into tasks(" +
		"owner_id, task) " +
		"VALUES(?, ?)"
	
	CLONE_TASK = "INSERT into tasks(" +
		"owner_id, delegate_id, origin_id, task) " +
		"VALUES(?, ?, ?, ?)"

	CREATE_HASHTAG = "INSERT into hashtags(" +
		"hashtag) " +
		"VALUES(?)"

	ADD_TASK_HASHTAG = "INSERT into task_hashtags(" +
		"task_id, hashtag_id) " +
		"VALUES(?, ?)"

	GET_HASHTAG = "SELECT id, hashtag " +
		"FROM hashtags " +
		"WHERE hashtag=?"

	GET_OPEN_TASKS_BY_USER = "SELECT id, owner_id, delegate_id, origin_id, " +
	  "task, actual, created " +
		"FROM tasks " +
		"WHERE (owner_id=? or delegate_id=?) and actual IS NULL"
	
	COMPLETE_TASK = "UPDATE tasks " +
		"SET actual=CURRENT_TIMESTAMP " +
		"where id=?"

)


func shortenURLs(s string) (string, error) {
  return s, nil
} // shortenURLs


func addHashtags(id string, task string) {

	hashtags := gowdl.ExtractHashtags(task)

	for _, hashtag := range hashtags {

		err := createHashtag(strings.TrimSpace(hashtag))

		if err != nil {
			log.Printf("%s addHashtags(): %s", APP_NAME, err.Error())
		}

	}

} // addHashtags


func addDelegates(id string, oid string, task string) {

	mentions := gowdl.ExtractMentions(task)
	
	for _, mention := range mentions {

		u := getUserByName(strings.TrimSpace(mention))

		if u == nil {
			log.Printf("%s addDelegates(): delegate user not found", APP_NAME)
		} else {
			cloneTask(id, u.ID, oid, task)
		}

	}

} // addDelegates


func cloneTask(id string, delegate_id string, origin_id string,
	task string) (error) {

	res, err := data.Exec(
		CLONE_TASK, id, delegate_id, origin_id, task,
	)

	if err != nil {
		return err
	} else {

		tid, err := res.LastInsertId()

		if err != nil {
			log.Printf("%s cloneTask(): %s", APP_NAME, err.Error())
			return err
		} else {

			hashtags := gowdl.ExtractHashtags(task)

			for _, hashtag := range hashtags {

				h := getHashtag(strings.TrimSpace(hashtag))

				if h != nil {
					addHashtagToTask(fmt.Sprintf("%d", tid), h.ID)
				}

			}

			return nil

		}

	}

} // cloneTask


func createHashtag(hashtag string) (error) {

	_, err := data.Exec(
		CREATE_HASHTAG, hashtag,
	)

	if err != nil {
		return err
	} else {
		return nil
	}

} // createHashtag


func getHashtag(hashtag string) *Hashtag {

	row := data.QueryRow(
		GET_HASHTAG, hashtag,
	)

	h := Hashtag{}

	err := row.Scan(&h.ID, &h.Tag)

	if err != nil {
		log.Printf("%s getHashtag(): %s", APP_NAME, err.Error())
		return nil
	} else {
		return &h
	}

} // getHashtag


func addHashtagToTask(tid string, hid string) (error) {

	_, err := data.Exec(
		ADD_TASK_HASHTAG, tid, hid,
	)

	if err != nil {
		return err
	} else {
		return nil
	}

} // addHashtagToTask


func createTask(uid string, task string) (error) {

	tx, err := data.Begin()

	if err != nil {
		log.Printf("%s createTask(): %s", APP_NAME, err.Error())
		return err
	} else {

		shortenedTask, err := shortenURLs(task)

		if err != nil {
			log.Println(err)
			return err
		} else {

			res, err := data.Exec(
				CREATE_TASK, uid, shortenedTask,
			)
		
			if err != nil {
		
				tx.Rollback()
				log.Println(err)
				return err
		
			} else {
	
				oid, err := res.LastInsertId()
	
				if err != nil {
					log.Printf("%s createTask(): %s", APP_NAME, err.Error())
					return err
				} else {
	
					addHashtags(uid, task)
	
					addDelegates(uid, fmt.Sprintf("%d", oid), shortenedTask)
	
					tx.Commit()
	
					return nil
	
				}
			
			}
	
		}

	}


} // createTask


func getOpenTasksByUser(id string) []Task {

	tasks := []Task{}

	rows, err := data.Query(
		GET_OPEN_TASKS_BY_USER, id, id,
	)

	defer rows.Close()

	if err != nil || err == sql.ErrNoRows {
		
		log.Printf("%s getOpenTasksByUser(): %s", APP_NAME, err.Error())
		return tasks

	} else {

		for rows.Next() {

			t := Task{}

			err := rows.Scan(&t.ID, &t.OwnerID, &t.DelegateID, &t.OriginID,
				&t.Task, &t.Actual, &t.Created)
				
			if err != nil || err == sql.ErrNoRows {

				log.Printf("%s getOpenTasksByUser(): %s", APP_NAME, err.Error())
				return tasks

			} else {
				tasks = append(tasks, t)
			}

		}

		return tasks

	}

} // getTasksByUser


func completeTask(tid string) bool {

	_, err := data.Exec(
		COMPLETE_TASK, tid,
	)

	if err != nil {
		log.Printf("%s completeTask(): %s", APP_NAME, err.Error())
		return false
	} else {
		return true
	}

} // completeTask


func taskHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
	case http.MethodGet:
		
		u := checkToken(r)

		if u == nil {
			w.WriteHeader(http.StatusUnauthorized)
		} else {

			tasks := getOpenTasksByUser(u.ID)

			j, err := json.Marshal(tasks)

			if err != nil {
				
				log.Printf("%s taskHandler(): %s", APP_NAME, err.Error())
				w.WriteHeader(http.StatusInternalServerError)

			} else {
				w.Write(j)
			}

		}
		
	case http.MethodPost:
		
		u := checkToken(r)

		if u == nil {
			w.WriteHeader(http.StatusUnauthorized)
		} else {

			task 			:= r.FormValue("task")

			if task == "" {
				w.WriteHeader(http.StatusBadRequest)
			} else {

				err := createTask(u.ID, task)

				if err != nil {
					log.Printf("%s taskHandler(): %s", APP_NAME, err.Error())
					w.WriteHeader(http.StatusConflict)
				}

			}
			
		}

	case http.MethodPut:

		vars := mux.Vars(r)

		id 		:= vars["id"]
		tid 	:= vars["tid"]

		u := checkToken(r)

		if u == nil {
			w.WriteHeader(http.StatusUnauthorized)
		} else {

			if u.ID == id {
				
				if !completeTask(tid) {
					w.WriteHeader(http.StatusInternalServerError)
				}

			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}


		}



  case http.MethodDelete:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // taskHandler
