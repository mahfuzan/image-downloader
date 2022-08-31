package controllers_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/mahfuzan/image-downloader/config"
	"github.com/mahfuzan/image-downloader/controllers"
)

var db *sqlx.DB

func TestGetImageList(t *testing.T) {
	truncateTable()
	addImage()
	// create a new HTTP request
	request, err := http.NewRequest("GET", "/download-image/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create new recorder to record response received by endpoint
	response := httptest.NewRecorder()

	// assign HTTP handler function
	handler := http.HandlerFunc(controllers.GetImageList)

	// hits endpoint with response recorder and request
	handler.ServeHTTP(response, request)

	// check if response is ok
	checkResponseCode(t, http.StatusOK, response.Code)

	// expected output from the endpoint
	expected := `{"success":true,"data":[{"id":1,"file_name":"testing.png","file_path":"./images/testing.png","source_url":"testing.png"}],"error":null}`

	// check response body if it is what we expect
	if response.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", response.Body.String(), expected)
	}
}

func TestEmptyTable(t *testing.T) {
	truncateTable()
	request, err := http.NewRequest("GET", "/download-image/", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.GetImageList)
	handler.ServeHTTP(response, request)
	checkResponseCode(t, http.StatusOK, response.Code)

	expected := `{"success":true,"data":[],"error":null}`
	if response.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", response.Body.String(), expected)
	}
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("handler returned wrong status code: got %v want %v", actual, expected)
	}
}

func TestGetImageById(t *testing.T) {
	truncateTable()
	addImage()

	request, err := http.NewRequest("GET", "/download-image/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.GetImageById)
	handler.ServeHTTP(response, request)
	checkResponseCode(t, http.StatusOK, response.Code)

	expected := `{"success":true,"data":{"id":1,"file_name":"testing.png","file_path":"./images/testing.png","source_url":"testing.png"},"error":null}`

	if response.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", response.Body.String(), expected)
	}
}

func TestGetImageNotExistent(t *testing.T) {
	truncateTable()
	request, err := http.NewRequest("GET", "/download-image/15", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.GetImageById)
	handler.ServeHTTP(response, request)
	checkResponseCode(t, http.StatusNotFound, response.Code)

	expected := `{"success":false,"data":null,"error":{"code":"RecordNotFound","desc":"Record not found in database"}}`
	if response.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", response.Body.String(), expected)
	}
}

func addImage() {
	db = config.GetDb()
	_, err := db.Exec("INSERT INTO images (file_name, file_path, source_url) VALUES ('testing.png', './images/testing.png', 'testing.png')")
	if err != nil {
		log.Fatal(err)
	}
}

func truncateTable() {
	db = config.GetDb()
	if _, err := db.Exec("TRUNCATE TABLE images"); err != nil {
		log.Fatal(err)
	}
}

func TestFailedParsing(t *testing.T) {
	request, err := http.NewRequest("GET", "/download-image/'1'", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.GetImageById)
	handler.ServeHTTP(response, request)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	expected := `{"success":false,"data":null,"error":{"code":"FailedParsing","desc":"Failed to parse parameter"}}`
	if response.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", response.Body.String(), expected)
	}
}

func TestDownloadImage(t *testing.T) {
	truncateTable()
	jsonStr := []byte(`{"url":"https://i.imgur.com/ONsnEhy.jpeg"}`)
	request, err := http.NewRequest("POST", "/download-image/", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.DownloadImage)
	handler.ServeHTTP(response, request)
	checkResponseCode(t, http.StatusOK, response.Code)

	expected := `{"success":true,"data":{"id":1,"file_name":"ONsnEhy.jpeg","file_path":"C:\\Users\\SIRCLO\\Downloads\\ONsnEhy.jpeg","source_url":"https://i.imgur.com/ONsnEhy.jpeg"},"error":null}`

	if response.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", response.Body.String(), expected)
	}
}

func TestFailedUnmarshal(t *testing.T) {
	jsonStr := []byte(`{\}`)
	request, err := http.NewRequest("POST", "/download-image/", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.DownloadImage)
	handler.ServeHTTP(response, request)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	expected := `{"success":false,"data":null,"error":{"code":"FailedUnmarshal","desc":"Failed to unmarshal data"}}`
	if response.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", response.Body.String(), expected)
	}
}

func TestFailedSaveFile(t *testing.T) {
	jsonStr := []byte(`{}`)
	request, err := http.NewRequest("POST", "/download-image/", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.DownloadImage)
	handler.ServeHTTP(response, request)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	expected := `{"success":false,"data":null,"error":{"code":"FailedSaveFile","desc":"Failed to save file to storage"}}`
	if response.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", response.Body.String(), expected)
	}
}

func TestFailedInsertToDb(t *testing.T) {
	db := config.GetDb()
	db.Close()
	jsonStr := []byte(`{"url":"https://i.imgur.com/ONsnEhy.jpeg"}`)
	request, err := http.NewRequest("POST", "/download-image/", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.DownloadImage)
	handler.ServeHTTP(response, request)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	expected := `{"success":false,"data":null,"error":{"code":"FailedInsertDb","desc":"Failed to insert data to database"}}`
	if response.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", response.Body.String(), expected)
	}
}
