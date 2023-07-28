package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func TestMain(t *testing.M) {
	router = gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	fmt.Println("test init")

	t.Run()
}

func TestGetAlbums(t *testing.T) {
	req, err := http.NewRequest("GET", "/albums", nil)
	if err != nil {
		t.Fatal(err)
	}

	reqRecorder := httptest.NewRecorder()

	router.ServeHTTP(reqRecorder, req)

	if status := reqRecorder.Code; status != http.StatusOK {
		t.Errorf("incorrect status code: got %v want %v", status, http.StatusOK)
	}

	str, err := json.MarshalIndent(albums, "", "    ")
	expected := string(str)
	if reqRecorder.Body.String() != expected {
		t.Errorf("incorrect response body: got %v want %v", reqRecorder.Body.String(), expected)
	}
}

func TestGetAlbumByID(t *testing.T) {
	req, err := http.NewRequest("GET", "/albums/2", nil)
	if err != nil {
		t.Fatal(err)
	}

	reqRecorder := httptest.NewRecorder()

	router.ServeHTTP(reqRecorder, req)

	if status := reqRecorder.Code; status != http.StatusOK {
		t.Errorf("incorrect status code: got %v want %v", status, http.StatusOK)
	}

	str, err := json.MarshalIndent(albums[1], "", "    ")
	expected := string(str)
	if reqRecorder.Body.String() != expected {
		t.Errorf("incorrect response body: got %v want %v", reqRecorder.Body.String(), expected)
	}
}

func TestPostAlbums(t *testing.T) {
	jsonBody := []byte(`{"id":"4","title":"test title", "artist":"testman", "price":99.99}`)
	req, err := http.NewRequest("POST", "/albums", bytes.NewReader(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	reqRecorder := httptest.NewRecorder()

	router.ServeHTTP(reqRecorder, req)

	if status := reqRecorder.Code; status != http.StatusCreated {
		t.Errorf("incorrect status code: got %v want %v", status, http.StatusOK)
	}

	var expected album
	json.Unmarshal(jsonBody, &expected)

	var res album
	json.Unmarshal(reqRecorder.Body.Bytes(), &res)
	if res.ID != expected.ID {
		t.Errorf("incorrect response body: got %v want %v", res.ID, expected.ID)
	}

	if res.Title != expected.Title {
		t.Errorf("incorrect response body: got %v want %v", res.ID, expected.Title)
	}

	if res.Artist != expected.Artist {
		t.Errorf("incorrect response body: got %v want %v", res.ID, expected.Artist)
	}

	if res.Price != expected.Price {
		t.Errorf("incorrect response body: got %v want %v", res.ID, expected.Price)
	}
}
