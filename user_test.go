package main

import (
	"log"
	"net/http"
	"strings"
	"testing"

)

func TestCreateUser(t *testing.T) {

	req, err := http.NewRequest("POST", userApi,
		strings.NewReader("email=sinbad@aol.com&password=abcxyz"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, received: %d", res.StatusCode)
	}

	cookies := res.Cookies()

	log.Println(cookies)

} // TestCreateUser


func TestCreateUserEmptyEmail(t *testing.T) {

	req, err := http.NewRequest("POST", userApi,
		strings.NewReader("email=&password=abcd"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected %d, received: %d", http.StatusBadRequest, res.StatusCode)
	}

} // TestCreateUserEmptyEmail


func TestCreateUserEmptyPassword(t *testing.T) {

	req, err := http.NewRequest("POST", userApi,
		strings.NewReader("email=sinbad@aol.com&password="))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected %d, received: %d", http.StatusBadRequest, res.StatusCode)
	}

} // TestCreateUserEmptyPassword


func TestCreateUserExisting(t *testing.T) {

	req, err := http.NewRequest("POST", userApi,
		strings.NewReader("email=sinbad@aol.com&password=abcd"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusForbidden {
		t.Errorf("Expected %d, received: %d", http.StatusForbidden, res.StatusCode)
	}

} // TestCreateUserExisting
