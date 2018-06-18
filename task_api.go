package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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

	GET_HASHTAG = "SELECT id, hashtag, description " +
		"FROM hashtags " +
		"WHERE hashtag=?"

	GET_TASKS_BY_USER = "SELECT id, owner_id, delegate_id, origin_id, task, " +
		"estimate, actual, created " +
		"FROM tasks " +
		"WHERE owner_id=? or delegate_id=?"
		
)


const (
	MENTION_SYMBOL				= "@"
	HASHTAG_SYMBOL				= "#"
	SPACE_SYMBOL          = " "
	COMMA_SYMBOL          = ","
)


func shortenURLs(s string) (string, error) {
  return s, nil
} // shortenURLs


func addHashtags(id string, task string, hashtags string) {

	if hashtags != "" {

		h := strings.Split(hashtags, COMMA_SYMBOL)

		for _, hashtag := range h {

			err := createHashtag(strings.TrimSpace(hashtag))

			if err != nil {
				log.Printf("%s addHashtags(): %s", APP_NAME, err.Error())
			}
	
		}

	}

} // addHashtags


func addDelegates(id string, oid string, task string,
	mentions string, hashtags string) {

	if mentions != "" {
		
		m := strings.Split(mentions, COMMA_SYMBOL)

		for _, mention := range m {

			log.Println(mention)
			u := getUserByName(strings.TrimSpace(mention))

			if u == nil {
				log.Printf("%s addDelegates(): delegate user not found", APP_NAME)
			} else {
				cloneTask(id, u.ID, oid, task, hashtags)
			}

		}

	}

} // addDelegates


func cloneTask(id string, delegate_id string, origin_id string, task string,
	hashtags string) (error) {

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

			tags := strings.Split(hashtags, COMMA_SYMBOL)

			for _, hashtag := range tags {

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

	err := row.Scan(&h.ID, &h.Tag, &h.Description)

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


func createTask(id string, task string, mentions string, hashtags string, 
	urls string) (error) {

	tx, err := data.Begin()

	if err != nil {
		log.Printf("%s createTask(): %s", APP_NAME, err.Error())
		return err
	} else {

		// TODO: shorten urls, use a db table or shortener service
		shortenedTask, err := shortenURLs(task)

		if err != nil {
			log.Println(err)
			return err
		} else {

			res, err := data.Exec(
				CREATE_TASK, id, shortenedTask,
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
	
					addHashtags(id, task, hashtags)
	
					addDelegates(id, fmt.Sprintf("%d", oid), shortenedTask, mentions,
					  hashtags)
	
					tx.Commit()
	
					return nil
	
				}
			
			}
	
		}

	}


} // createTask


func getTasksByUser(id string) []Task {

	tasks := []Task{}

	rows, err := data.Query(
		GET_TASKS_BY_USER, id, id,
	)

	defer rows.Close()

	if err != nil || err == sql.ErrNoRows {
		
		log.Printf("%s getTasksByUser(): %s", err.Error())
		return tasks

	} else {

		for rows.Next() {

			t := Task{}

			err := rows.Scan(&t.ID, &t.OwnerID, &t.DelegateID, &t.OriginID,
				&t.Task, &t.Estimate, &t.Actual, &t.Created)
				
			if err != nil || err == sql.ErrNoRows {

				log.Printf("%s getTasksByUser(): %s", err.Error())
				return tasks

			} else {
				tasks = append(tasks, t)
			}

		}

		return tasks

	}

} // getTasksByUser


func taskHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
	case http.MethodGet:
		
		u := checkToken(r)

		if u == nil {
			w.WriteHeader(http.StatusUnauthorized)
		} else {

			tasks := getTasksByUser(u.ID)

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
			mentions 	:= r.FormValue("mentions")
			hashtags  := r.FormValue("hashtags")
			urls      := r.FormValue("urls")

			if task == "" {
				w.WriteHeader(http.StatusBadRequest)
			} else {

				err := createTask(u.ID, task, mentions, hashtags, urls)

				if err != nil {
					log.Println("%s taskHandler(): %s", APP_NAME, err.Error())
					w.WriteHeader(http.StatusConflict)
				}

			}
			
		}

  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // taskHandler
