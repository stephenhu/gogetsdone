package main

import (
	"fmt"
	"log"
	"net/http/httptest"
	"os"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/sqlite3"
	_ "github.com/golang-migrate/migrate/source/file"

)

const (
	TEST_DRIVER 		= "sqlite3://data/getsdone.test.db"
	TEST_DATABASE		= "./data/getsdone.test.db"
	TEST_SOURCE     = "file://db/migrations"
)

var server *httptest.Server
var userApi string
var db string

func newDatabase(source string, driver string) {

  m, err := migrate.New(
		source,
		driver,
	)

	if err != nil {
		log.Println(err)
	} else {
		m.Steps(2)
	}

} // newDatabase


func setDatabase(name string) {

  *database = name

} // setDatabase


func init() {

	cleanup()

  server = httptest.NewServer(initRoutes())
	
	newDatabase(TEST_SOURCE, TEST_DRIVER)

  setDatabase(TEST_DATABASE)

	connectDatabase()

  userApi = fmt.Sprintf("%s/api/users", server.URL)

} // init

func cleanup() {

	_, err := os.Stat(TEST_DATABASE)

	if err == nil {
		os.Remove(TEST_DATABASE)
	}

} // cleanup
