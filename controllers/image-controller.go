package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mahfuzan/image-downloader/models"
)

type Url struct {
	Url string `json:"url"`
}

func DownloadImage(w http.ResponseWriter, r *http.Request) {
	// read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	// decode json
	var url Url
	err = json.Unmarshal(body, &url)
	if err != nil {
		panic(err)
	}

	// save file to storage
	filename := filepath.Base(url.Url)
	filePath := "./images/" + filename
	err = saveFile(url.Url, filename, filePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File %s successfully downloaded\n", filename)

	// save to database
	res, err := models.SaveToDatabase(url.Url, filename, filePath)
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	result, err := json.Marshal(res)
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func saveFile(url string, filename string, filePath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New("received non 200 response code")
	}

	// create empty file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func GetImageList(w http.ResponseWriter, r *http.Request) {
	result, err := models.GetList()
	if err != nil {
		var failedResponse = map[string]string{
			"success": "false",
			"message": "Data not found",
		}
		res, err := json.Marshal(failedResponse)
		if err != nil {
			w.Header().Set("Content-Type", "pkglication/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte{})
			return
		} else {
			w.Header().Set("Content-Type", "pkglication/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write(res)
			return
		}
	}

	res, err := json.Marshal(result)
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetImageById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 0, 0)
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	result, err := models.GetById(id)
	if err != nil {
		var failedResponse = map[string]string{
			"success": "false",
			"message": "Data not found",
		}
		res, err := json.Marshal(failedResponse)
		if err != nil {
			w.Header().Set("Content-Type", "pkglication/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte{})
			return
		}

		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(res)
		return
	}

	res, err := json.Marshal(result)
	if err != nil {
		w.Header().Set("Content-Type", "pkglication/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
