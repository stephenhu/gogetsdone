package main

import (
	"net/http"
	"testing"

)

func TestIconUpload(t *testing.T) {

	req, err := http.NewRequest("GET", versionApi, nil)
	
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, received: %d", res.StatusCode)
	}

} // TestIconUpload
