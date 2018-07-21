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

	GET_OPEN_TASKS_BY_USER = "SELECT tasks.id, tasks.owner_id, tasks.delegate_id, " +
		"tasks.origin_id, tasks.task, tasks.actual, tasks.created, users.name " +
		"FROM tasks, users " +
		"WHERE (tasks.owner_id=? or tasks.delegate_id=?) and (tasks.owner_id=users.id) " +
		"and tasks.actual IS NULL and tasks.deferred=0"

	GET_COMPLETED_TASKS_BY_USER = "SELECT tasks.id, tasks.owner_id, " +
		"tasks.delegate_id, tasks.origin_id, tasks.task, tasks.actual, " +
		"tasks.created, users.name " +
		"FROM tasks, users " +
		"WHERE (tasks.owner_id=? or tasks.delegate_id=?) and (tasks.owner_id=users.id) " +
		"and tasks.actual IS NOT NULL"

	GET_DEFERRED_TASKS_BY_USER = "SELECT tasks.id, tasks.owner_id, " +
		"tasks.delegate_id, tasks.origin_id, tasks.task, tasks.actual, " +
		"tasks.created, users.name " +
		"FROM tasks, users " +
		"WHERE (tasks.owner_id=? or tasks.delegate_id=?) and tasks.deferred=1 " +
		"and tasks.actual IS NULL"

	GET_TASK = "SELECT tasks.id, tasks.owner_id, tasks.delegate_id, " +
		"tasks.origin_id, tasks.task, tasks.actual, tasks.created, users.name " +
		"FROM tasks, users " +
		"WHERE tasks.id=? and tasks.owner_id=users.id"

	COMPLETE_TASK = "UPDATE tasks " +
		"SET actual=CURRENT_TIMESTAMP " +
		"WHERE id=?"

	DEFER_TASK = "UPDATE tasks " +
		"SET deferred=1 " +
		"WHERE id=?"

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

			mentions := gowdl.ExtractMentions(task)

			if len(mentions) > 0 {
				return nil
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
		
					log.Println(oid)

					if err != nil {
						log.Printf("%s createTask(): %s", APP_NAME, err.Error())
						return err
					} else {
		
						addHashtags(uid, task)
		
						//addDelegates(uid, fmt.Sprintf("%d", oid), shortenedTask)
		
						tx.Commit()
		
						return nil
		
					}
				
				}

			}
			
	
		}

	}


} // createTask


func getTasksByUser(id string, view string) []Task {

	var rows *sql.Rows
	var err error

	tasks := []Task{}
	
	log.Println(view)

	if view == TASK_COMPLETED {
	
		rows, err = data.Query(
			GET_COMPLETED_TASKS_BY_USER, id, id,
		)

	} else if view == TASK_DEFERRED {
		
		rows, err = data.Query(
			GET_DEFERRED_TASKS_BY_USER, id, id,
		)

	} else {
		
		rows, err = data.Query(
			GET_OPEN_TASKS_BY_USER, id, id,
		)

	}

	defer rows.Close()

	if err != nil || err == sql.ErrNoRows {
		
		log.Printf("%s getOpenTasksByUser(): %s", APP_NAME, err.Error())
		return tasks

	} else {

		for rows.Next() {

			t := Task{}

			err := rows.Scan(&t.ID, &t.OwnerID, &t.DelegateID, &t.OriginID,
				&t.Task, &t.Actual, &t.Created, &t.OwnerName)
				
			if err != nil || err == sql.ErrNoRows {

				log.Printf("%s getOpenTasksByUser(): %s", APP_NAME, err.Error())
				return tasks

			} else {

				comments := getCommentsByTask(t.ID)

				log.Println(comments)

				t.Comments = comments

				tasks = append(tasks, t)
			}

		}

		return tasks

	}

} // getTasksByUser


func getTask(id string) *Task {

	row := data.QueryRow(
		GET_TASK, id,
	)

	t := Task{}

	err := row.Scan(&t.ID, &t.OwnerID, &t.DelegateID, &t.OriginID,
		&t.Task, &t.Actual, &t.Created, &t.OwnerName)

	if err != nil || err == sql.ErrNoRows {
		log.Println("gogetsdone getTask(): ", err)
		return nil
	}

	comments := getCommentsByTask(t.ID)

	t.Comments = comments
	
	return &t

} // getTask


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

			// TODO: authorization?

			vars := mux.Vars(r)

			id	:= vars["id"]
			tid	:= vars["tid"]

			log.Println(id)
			
			if tid == "" {

				view := r.FormValue("view")

				var tasks []Task
	
				tasks = getTasksByUser(u.ID, view)	
	
				j, err := json.Marshal(tasks)
	
				if err != nil {
					
					log.Printf("%s taskHandler(): %s", APP_NAME, err.Error())
					w.WriteHeader(http.StatusInternalServerError)
	
				} else {
					w.Write(j)
				}

			} else {

				task := getTask(tid)

				j, err := json.Marshal(task)

				if err != nil {

					log.Printf("%s taskHandler(): %s", APP_NAME, err.Error())
					w.WriteHeader(http.StatusInternalServerError)

				} else {
					w.Write(j)
				}

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
