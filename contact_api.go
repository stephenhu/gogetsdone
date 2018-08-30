package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"

)


const (

	CREATE_CONTACT	= "INSERT into contacts(" +
		"user_id, contact_id, contact_state_id) " +
		"VALUES(?, ?, (" +
		"SELECT id from contact_states where name=?))"

	REMOVE_CONTACT	= "DELETE from contacts " +
		"WHERE user_id=? and contact_id=?"

	ACCEPT_CONTACT  = "UPDATE contacts " +
		"SET contact_state_id=(" +
		"SELECT id from contact_states where name='accepted') " +
		"WHERE user_id=? and contact_id=?"

	DECLINE_CONTACT  = "UPDATE contacts " +
		"SET contact_state_id=(" +
		"SELECT id from contact_states where name='declined') " +
		"WHERE user_id=? and contact_id=?"

	GET_CONTACTS = "SELECT contacts.id, contacts.contact_id, " +
	  "contact_states.name, users.name, users.icon " +
	  "from contacts, contact_states, users " +
		"WHERE contacts.user_id=? and contacts.contact_id=users.id and " +
		"contact_states.id=contacts.contact_state_id"

	GET_CONTACT_REQUESTS = "SELECT contacts.id, contacts.contact_id, " +
		"contacts.user_id, contact_states.name, users.name, users.icon " +
		"from contacts, contact_states, users " +
		"WHERE contacts.contact_id=? and contacts.user_id=users.id and " +
		"contact_states.id=contacts.contact_state_id and " +
		"contacts.contact_state_id=(SELECT id from contact_states where name='requested')"

)


func getContactRequests(uid string) []Contact {

	rows, err := data.Query(
		GET_CONTACT_REQUESTS, uid,
	)

	defer rows.Close()

	contacts := []Contact{}

	if err != nil || err == sql.ErrNoRows {
		log.Printf("%s getContactRequests(): %s", APP_NAME, err.Error())
	} else {

		for rows.Next() {

			c := Contact{}

			err := rows.Scan(&c.ID, &c.ContactID, &c.UserID, &c.State, &c.ContactName, 
				&c.ContactIcon)
			
			if err != nil {

				log.Printf("%s getContactRequests(): %s", APP_NAME, err.Error())
				return contacts

			}
			
			contacts = append(contacts, c)

		}

	}

	return contacts

} // getContactRequests


func getContacts(uid string) []Contact {

	rows, err := data.Query(
		GET_CONTACTS, uid,
	)

	defer rows.Close()

	contacts := []Contact{}

	if err != nil || err == sql.ErrNoRows {
		log.Printf("%s getContacts(): %s", APP_NAME, err.Error())
	} else {

		for rows.Next() {

			c := Contact{}

			err := rows.Scan(&c.ID, &c.ContactID, &c.State, &c.ContactName, 
				&c.ContactIcon)
			
			if err != nil {

				log.Printf("%s getContacts(): %s", APP_NAME, err.Error())
				return contacts

			}
			
			contacts = append(contacts, c)

		}

	}

	return contacts

} // getContacts


func addContact(uid string, cid string, state string) (error) {

	_, err := data.Exec(
		CREATE_CONTACT, uid, cid, state,
	)

	if err != nil {

		log.Printf("%s addContact(): %s", APP_NAME, err.Error())
		return err

	} else {
		return nil
	}

} // addContact


func removeContact(uid string, cid string) (error) {

	_, err := data.Exec(
		REMOVE_CONTACT, uid, cid,
	)

	if err != nil {
		log.Printf("%s removeContact(): %s", APP_NAME, err.Error())
		return err
	} else {
		return nil
	}

} // removeContact


func updateContact(uid string, cid string, action string) (error) {

	var statement string

	if action == CONTACT_ACCEPTED {
		statement = ACCEPT_CONTACT
	} else if action == CONTACT_DECLINED {
		statement = DECLINE_CONTACT
	} else {
		return errors.New("Non-registered contact event") 
	}

	_, err := data.Exec(
		statement, uid, cid,
	)

	if err != nil {

		log.Printf("%s updateContact(): %s", APP_NAME, err.Error())
		return err

	} else {
		return nil
	}

} // updateContact


func contactHandler(w http.ResponseWriter, r *http.Request) {

	u := checkToken(r)

	if u == nil {
		w.WriteHeader(http.StatusUnauthorized)
	} else {

		vars := mux.Vars(r)

		id 	:= vars["id"]
		cid := vars["cid"]

		if id == u.ID {

			switch r.Method {
			case http.MethodPost:
	
				email := r.FormValue("email")

				user := getUserByEmail(email)

				if user == nil {
					w.WriteHeader(http.StatusNotFound)
				} else {

					if user.ID == u.ID {
						w.WriteHeader(http.StatusBadRequest)
					} else {

						tx, err := data.Begin()

						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
						} else {
							
							err := addContact(u.ID, user.ID, CONTACT_REQUESTED)

							if err != nil {
								
								tx.Rollback()
								w.WriteHeader(http.StatusInternalServerError)
		
							} else {
		
								err := addContact(user.ID, u.ID, CONTACT_PENDING)
		
								if err != nil {
		
									tx.Rollback()
									w.WriteHeader(http.StatusInternalServerError)
		
								} else {
									tx.Commit()
								}
		
							}

						}
		
					}
	
				}
	
			case http.MethodGet:

				contacts := getContacts(id)

				j, err := json.Marshal(contacts)

				if err != nil {

					log.Printf("%s contactHandler(): %s", APP_NAME, err.Error())
					w.WriteHeader(http.StatusInternalServerError)

				} else {
					w.Write(j)
				}
			
			case http.MethodPut:

				if cid == "" {
					w.WriteHeader(http.StatusBadRequest)
				} else {


					action := r.FormValue("action")

					tx, err := data.Begin()

					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
					} else {

						err := updateContact(u.ID, cid, action)

						if err != nil {
							
							tx.Rollback()
							log.Printf("%s contactHandler(): %s", APP_NAME, err.Error())
							w.WriteHeader(http.StatusBadRequest)

						} else {

							err := updateContact(cid, u.ID, action)

							if err != nil {
							
								tx.Rollback()
								log.Printf("%s contactHandler(): %s", APP_NAME, err.Error())
								w.WriteHeader(http.StatusBadRequest)
								
							} else {

								tx.Commit()

							} 

						}

					}				

				}

			case http.MethodDelete:
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
	
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}

	}

} // contactHandler
