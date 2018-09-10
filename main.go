package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/golang-migrate/migrate/database/sqlite3"
	"github.com/gorilla/mux"
)

var database		= flag.String("database", "./data/getsdone.db", "database address")
var port				= flag.String("port", "8888", "service port")
var domain      = flag.String("domain", LOCAL_HOST, "domain address")

var data *sql.DB = nil


func logf(msg string, fname string) {
	log.Printf("%s %s(): %s", version(), fname, msg)
} // logf


func version() string {
  return fmt.Sprintf("%s v%s", APP_NAME, VERSION)
} // version


func connectDatabase() {
	
	_, err := os.Stat(*database)

	if err != nil || os.IsNotExist(err) {
		log.Println(err)
		log.Fatal("Database not found, please initialize database")
	} else {

		db, err := sql.Open("sqlite3", *database)

		if err != nil {
			log.Fatal("Database connection error: ", err)
		}

		data = db

	}

} // connectDatabase


func initRoutes() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/auth", authHandler)

	router.HandleFunc("/api/version", versionHandler)
	router.HandleFunc("/api/users", userHandler)
	
	router.HandleFunc("/api/users/{id:[0-9]+}/tasks", taskHandler)
	router.HandleFunc("/api/users/{id:[0-9]+}/tasks/{tid:[0-9]+}", taskHandler)
	router.HandleFunc("/api/users/{id:[0-9]+}/tasks/{tid:[0-9]+}/comments",
	  commentHandler)
	
	router.HandleFunc("/api/users/{id:[0-9]+}/contacts", contactHandler)
	router.HandleFunc("/api/users/{id:[0-9]+}/contacts/{cid:[0-9]+}",
	  contactHandler)

	router.HandleFunc("/api/users/{id:[0-9]+}/notifications", notificationHandler)

	return router

} // initRoutes

func main() {

	flag.Parse()

	connectDatabase()

	router := initRoutes()

	addr := fmt.Sprintf(":%s", *port)

	log.Printf("%s listening on port %s", APP_NAME, *port)

	log.Fatal(http.ListenAndServe(addr, router))


} // main
