package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/stephenhu/gowdl"

)

const (

	CREATE_TASK = "INSERT into tasks(" +
		"owner_id, delegate_id, task) " +
		"VALUES(?, ?, ?)"
	
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

	GET_ASSIGNED_TASKS_BY_USER = "SELECT tasks.id, tasks.owner_id, tasks.delegate_id, " +
		"tasks.origin_id, tasks.task, tasks.actual, tasks.created, users.name " +
		"FROM tasks, users " +
		"WHERE tasks.owner_id=? and tasks.delegate_id=users.id " +
		"and tasks.actual IS NULL and tasks.deferred=0"

	GET_COMPLETED_TASKS_BY_USER = "SELECT tasks.id, tasks.owner_id, " +
		"tasks.delegate_id, tasks.origin_id, tasks.task, tasks.actual, " +
		"tasks.created, users.name " +
		"FROM tasks, users " +
		"WHERE (tasks.owner_id=? or tasks.delegate_id=?) and (tasks.owner_id=users.id) " +
		"and tasks.actual IS NOT NULL " +
		"ORDER BY tasks.created DESC"

	GET_DEFERRED_TASKS_BY_USER = "SELECT tasks.id, tasks.owner_id, " +
		"tasks.delegate_id, tasks.origin_id, tasks.task, tasks.actual, " +
		"tasks.created, users.name " +
		"FROM tasks, users " +
		"WHERE (tasks.owner_id=? or tasks.delegate_id=?) and tasks.deferred=1 and (tasks.owner_id=users.id) " +
		"and tasks.actual IS NULL"

	GET_TASK = "SELECT tasks.id, tasks.owner_id, tasks.delegate_id, " +
		"tasks.origin_id, tasks.task, tasks.actual, tasks.created, users.name " +
		"FROM tasks, users " +
		"WHERE tasks.id=? and tasks.owner_id=users.id"

	COMPLETE_TASK = "UPDATE tasks " +
		"SET actual=CURRENT_TIMESTAMP " +
		"WHERE id=?"

	UPDATE_TASK = "UPDATE tasks " +
		"SET deferred=? " +
		"WHERE id=?"

)


func shortenURLs(s string) (string, error) {
	
	if len(s) < 1 {
		return "", errors.New("Invalid string")
	} else {
		return strings.Replace(s, "\n", " ", 1), nil
	}

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


func createTasksForDelegates(id string, mentions []string,
	task string) (error) {

	for _, m := range mentions {

		log.Println(m)
		u := getUserByName(strings.ToLower(m))
		log.Println(u)
		
		if u == nil {
			return errors.New("Delegate user not found.")
		} else {
			
			err := createTask(id, &u.ID, task)

			if err != nil {
				return err
			}

		}

	}

	return nil

} // createTasksForDelegates


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


func createTask(uid string, did *string, task string) (error) {

	shortenedTask, err := shortenURLs(task)

	if err != nil {
		log.Println(err)
		return err
	} else {
			
		_, err := data.Exec(
			CREATE_TASK, uid, did, shortenedTask,
		)

		if err != nil {

			log.Println(err)
			return err

		} else {

			addHashtags(uid, task)

			return nil

		}

	}

} // createTask


func getTasksByUser(id string, view string) []Task {

	var rows *sql.Rows
	var err error

	tasks := []Task{}

	if view == TASK_COMPLETED {
	
		rows, err = data.Query(
			GET_COMPLETED_TASKS_BY_USER, id, id,
		)

	} else if view == TASK_DEFERRED {
		
		rows, err = data.Query(
			GET_DEFERRED_TASKS_BY_USER, id, id,
		)

	} else if view == TASK_ASSIGNED {

		rows, err = data.Query(
			GET_ASSIGNED_TASKS_BY_USER, id,
		)

	} else {
		
		rows, err = data.Query(
			GET_OPEN_TASKS_BY_USER, id, id,
		)

	}

	defer rows.Close()

	if err != nil || err == sql.ErrNoRows {
		
		log.Printf("%s getTasksByUser(): %s", APP_NAME, err.Error())
		return tasks

	} else {

		for rows.Next() {

			t := Task{}

			err := rows.Scan(&t.ID, &t.OwnerID, &t.DelegateID, &t.OriginID,
				&t.Task, &t.Actual, &t.Created, &t.OwnerName)
			
			if err != nil || err == sql.ErrNoRows {

				log.Printf("%s getTasksByUser(): %s", APP_NAME, err.Error())
				return tasks

			} else {

				comments := getCommentsByTask(t.ID)

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


func updateTask(tid string, state int) bool {

	if state != UPDATE_TASK_DEFERRED && state != UPDATE_TASK_UNDEFERRED {
		return false
	}

	_, err := data.Exec(
		UPDATE_TASK, state, tid,
	)

	if err != nil {
		log.Printf("%s deferTask(): %s", APP_NAME, err.Error())
		return false
	} else {
		return true
	}

} // updateTask


func taskHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
	case http.MethodGet:
		
		u := checkToken(r)

		if u == nil {
			w.WriteHeader(http.StatusUnauthorized)
		} else {

			vars := mux.Vars(r)

			id	:= vars["id"]
			tid	:= vars["tid"]
			
			if id == u.ID {

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

			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}


		}
		
	case http.MethodPost:
		
		u := checkToken(r)

		if u == nil {
			w.WriteHeader(http.StatusUnauthorized)
		} else {

			task := r.FormValue("task")

			if task == "" {
				w.WriteHeader(http.StatusBadRequest)
			} else {

				mentions := gowdl.ExtractMentions(task)

				if len(mentions) > 0 {

					if checkContacts(u.ID, mentions) {

						tx, err := data.Begin()

						if err != nil {

							log.Printf("%s taskHandler(): %s", APP_NAME, err.Error())
							w.WriteHeader(http.StatusInternalServerError)

						} else {

							err := createTasksForDelegates(u.ID, mentions, task)

							if err != nil {

								logf(err.Error(), "taskHandler")
								tx.Rollback()
								w.WriteHeader(http.StatusInternalServerError)

							} else {
								tx.Commit()
							}

						}

					} else {
						
						log.Printf(
							"%s taskHandler(): One or more delegates are not " +
							"listed in your contacts or you have delegated to yourself which should not be allowed.", APP_NAME)
						w.WriteHeader(http.StatusConflict)
						
					}

				} else {

					tx, err := data.Begin()

					if err != nil {

						log.Printf("%s taskHandler(): %s", APP_NAME, err.Error())
						w.WriteHeader(http.StatusInternalServerError)

					} else {
						
						err := createTask(u.ID, nil, task)

						if err != nil {
							
							tx.Rollback()
							log.Printf("%s taskHandler(): %s", APP_NAME, err.Error())
							w.WriteHeader(http.StatusInternalServerError)

						} else {
							tx.Commit()
						}

					}
	
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

				action := r.FormValue("action")

				if action == ACTION_COMPLETED {

					if !completeTask(tid) {
						w.WriteHeader(http.StatusInternalServerError)
					}
	
				} else if action == ACTION_DEFERRED {

					if !updateTask(tid, UPDATE_TASK_DEFERRED) {
						w.WriteHeader(http.StatusInternalServerError)
					}

				} else if action == ACTION_UNDEFERRED {

					if !updateTask(tid, UPDATE_TASK_UNDEFERRED) {
						w.WriteHeader(http.StatusInternalServerError)
					}

				} else {
					w.WriteHeader(http.StatusBadRequest)
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
