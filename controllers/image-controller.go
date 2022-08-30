package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

type Response struct {
	Success bool           `json:"success"`
	Data    []models.Image `json:"data"`
	Error   ErrorResponse  `json:"error"`
}

type ErrorResponse struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}

func DownloadImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Success: false,
	}

	// read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Error.Code = "FailedReadData"
		response.Error.Desc = "Failed to read response body data"

		responseData, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unknown error"))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseData)
		return
	}

	// decode json
	var url Url
	err = json.Unmarshal(body, &url)
	if err != nil {
		response.Error.Code = "FailedUnmarshal"
		response.Error.Desc = "Failed to unmarshal data"

		responseData, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unknown error"))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseData)
		return
	}

	// save file to storage
	filename := filepath.Base(url.Url)
	filePath := "./images/" + filename
	err = saveFile(url.Url, filename, filePath)
	if err != nil {
		response.Error.Code = "FailedSaveFile"
		response.Error.Desc = "Failed to save file to storage"

		responseData, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unknown error"))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseData)
		return
	}
	fmt.Printf("File %s successfully downloaded\n", filename)

	// save to database
	imageData, err := models.SaveToDatabase(url.Url, filename, filePath)
	if err != nil {
		response.Error.Code = "FailedInsertDb"
		response.Error.Desc = "Failed to insert data to database"

		responseData, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unknown error"))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseData)
		return
	}

	image := []models.Image{imageData}
	response.Success = true
	response.Data = image
	responseData, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unknown error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
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
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Success: false,
	}

	imageList, err := models.GetList()
	if err != nil {
		response.Error.Code = "RecordNotFound"
		response.Error.Desc = "Record not found in database"
		res, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unknown Error"))
			return
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write(res)
			return
		}
	}

	response.Success = true
	response.Data = imageList
	responseData, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unknown Error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func GetImageById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Success: false,
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 0, 0)
	if err != nil {
		response.Error.Code = "FailedParsing"
		response.Error.Desc = "Failed to parse parameter"

		res, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unknown error"))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	imageData, err := models.GetById(id)
	if err != nil {
		response.Error.Code = "RecordNotFound"
		response.Error.Desc = "Record not found in database"
		res, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unknown error"))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(res)
		return
	}

	image := []models.Image{imageData}
	response.Success = true
	response.Data = image
	responseData, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unknown Error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}
