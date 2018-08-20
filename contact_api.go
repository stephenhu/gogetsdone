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
		"user_id, contact_id) " +
		"VALUES(?, ?)"

	REMOVE_CONTACT	= "DELETE from contacts " +
		"WHERE user_id=? and contact_id=?"

	ACCEPT_CONTACT  = "UPDATE contacts " +
		"SET state=1 " +
		"WHERE user_id=? and contact_id=?"

	DECLINE_CONTACT  = "UPDATE contacts " +
		"SET state=2 " +
		"WHERE user_id=? and contact_id=?"

	GET_CONTACTS = "SELECT contacts.id, contacts.contact_id, " +
	  "contacts.accepted, users.name, users.icon " +
	  "from contacts, users " +
	  "WHERE contacts.user_id=? and contacts.contact_id=users.id"

	GET_CONTACT_REQUESTS = "SELECT id, contact_id from contacts " +
	  "WHERE user_id=? and contact_id=? and accepted=0"

)


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


func addContact(uid string, cid string) (error) {

	_, err := data.Exec(
		CREATE_CONTACT, uid, cid,
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
	
				if cid != "" {
	
					tx, err := data.Begin()

					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
					} else {

						err := addContact(u.ID, cid)

						if err != nil {
							
							tx.Rollback()
							w.WriteHeader(http.StatusInternalServerError)

						} else {

							err := addContact(cid, u.ID)

							if err != nil {

								tx.Rollback()
								w.WriteHeader(http.StatusInternalServerError)

							} else {
								tx.Commit()
							}

						}
						
					}

				} else {
					w.WriteHeader(http.StatusBadRequest)
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
