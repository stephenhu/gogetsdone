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

const (
	APP_NAME				= "getsdone"
	GETSDONE        = APP_NAME
	HASH_LENGTH     = 64
	HMAC_KEY        = "spain this summer"
	LOCALHOST       = "127.0.0.1"
	PEPPER          = "getsdone is the bomb"
	SALT_LENGTH     = 32
	TOKEN_LENGTH    = 32
	VERSION					= "0.1"
)

var database		= flag.String("database", "./data/getsdone.db", "database address")
var port				= flag.String("port", "8888", "service port")
var domain      = flag.String("domain", LOCALHOST, "domain address")

var data *sql.DB = nil

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

	router.HandleFunc("/api/users", userHandler)
	router.HandleFunc("/api/tasks", taskHandler)
	router.HandleFunc("/api/users/{id:[0-9]+}/follows", followHandler)

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
