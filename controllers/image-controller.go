package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/mahfuzan/image-downloader/models"
)

type Url struct {
	Url string `json:"url"`
}

type Response struct {
	Success bool           `json:"success"`
	Data    interface{}    `json:"data"`
	Error   *ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}

const ERR_FAILED_READ_DATA_CODE string = "FailedReadData"
const ERR_FAILED_READ_DATA_DESC string = "Failed to read response body data"
const ERR_FAILED_UNMARSHAL_CODE string = "FailedUnmarshal"
const ERR_FAILED_UNMARSHAL_DESC string = "Failed to unmarshal data"
const ERR_FAILED_SAVE_FILE_CODE string = "FailedSaveFile"
const ERR_FAILED_SAVE_FILE_DESC string = "Failed to save file to storage"
const ERR_FAILED_INSERT_DB_CODE string = "FailedInsertDb"
const ERR_FAILED_INSERT_DB_DESC string = "Failed to insert data to database"
const ERR_FAILED_GET_DATA_DB_CODE string = "FailedGetDb"
const ERR_FAILED_GET_DATA_DB_DESC string = "Failed to get data from database"
const ERR_NOT_FOUND_CODE string = "RecordNotFound"
const ERR_NOT_FOUND_DESC string = "Record not found in database"
const ERR_FAILED_PARSING_CODE string = "FailedParsing"
const ERR_FAILED_PARSING_DESC string = "Failed to parse parameter"
const ERR_UNKOWN string = "Unknown Error"

func DownloadImage(w http.ResponseWriter, r *http.Request) {
	// read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		returnErrorResponse(w, http.StatusBadRequest, ERR_FAILED_READ_DATA_CODE, ERR_FAILED_READ_DATA_DESC)
		return
	}

	// decode json
	var url Url
	err = json.Unmarshal(body, &url)
	if err != nil {
		returnErrorResponse(w, http.StatusBadRequest, ERR_FAILED_UNMARSHAL_CODE, ERR_FAILED_UNMARSHAL_DESC)
		return
	}

	// save file to storage
	filename := path.Base(url.Url)
	homeDir, _ := os.UserHomeDir()
	filePath := filepath.Join(homeDir, "Downloads", filename)
	err = SaveFile(url.Url, filename, filePath)
	if err != nil {
		returnErrorResponse(w, http.StatusBadRequest, ERR_FAILED_SAVE_FILE_CODE, ERR_FAILED_SAVE_FILE_DESC)
		return
	}
	fmt.Printf("File %s successfully downloaded\n", filename)

	// save to database
	imageData, err := models.SaveToDatabase(url.Url, filename, filePath)
	if err != nil {
		returnErrorResponse(w, http.StatusBadRequest, ERR_FAILED_INSERT_DB_CODE, ERR_FAILED_INSERT_DB_DESC)
		return
	}

	response := Response{
		Success: true,
		Data:    imageData,
	}
	returnResponse(w, http.StatusOK, response)
}

func SaveFile(url string, filename string, filePath string) error {
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
	imageList, err := models.GetList()
	if err != nil {
		returnErrorResponse(w, http.StatusBadRequest, ERR_FAILED_GET_DATA_DB_CODE, ERR_FAILED_GET_DATA_DB_DESC)
		return
	}

	response := Response{
		Success: true,
		Data:    imageList,
	}
	returnResponse(w, http.StatusOK, response)
}

func GetImageById(w http.ResponseWriter, r *http.Request) {
	varId := path.Base(r.URL.String())
	id, err := strconv.ParseInt(varId, 0, 0)
	if err != nil {
		returnErrorResponse(w, http.StatusBadRequest, ERR_FAILED_PARSING_CODE, ERR_FAILED_PARSING_DESC)
		return
	}

	imageData, err := models.GetById(id)
	if err != nil {
		returnErrorResponse(w, http.StatusNotFound, ERR_NOT_FOUND_CODE, ERR_NOT_FOUND_DESC)
		return
	}

	response := Response{
		Success: true,
		Data:    imageData,
	}
	returnResponse(w, http.StatusOK, response)
}

func returnResponse(w http.ResponseWriter, httpStatus int, response any) {
	responseData, err := json.Marshal(response)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(ERR_UNKOWN))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)
		w.Write(responseData)
	}
}

func returnErrorResponse(w http.ResponseWriter, httpStatus int, code string, desc string) {
	errorData := ErrorResponse{
		Code: code,
		Desc: desc,
	}

	response := Response{
		Success: false,
		Error:   &errorData,
	}

	returnResponse(w, httpStatus, response)
}
