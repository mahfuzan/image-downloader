package controllers_test

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/mahfuzan/image-downloader/controllers"
	"github.com/mahfuzan/image-downloader/models"
)

func TestGetImageList(t *testing.T) {
	imageData := []models.Image{
		{
			Id:        1,
			FileName:  "testing.png",
			FilePath:  "\\images\\testing.png",
			SourceUrl: "images/testing.png",
		},
		{
			Id:        2,
			FileName:  "testing.png",
			FilePath:  "\\images\\testing.png",
			SourceUrl: "images/testing.png",
		},
	}

	response := controllers.Response{
		Success: true,
		Data:    imageData,
	}

	responseData, _ := json.Marshal(response)
	expected := `{"success":true,"data":[{"id":1,"file_name":"testing.png","file_path":"\\images\\testing.png","source_url":"images/testing.png"},{"id":2,"file_name":"testing.png","file_path":"\\images\\testing.png","source_url":"images/testing.png"}],"error":null}`

	if string(responseData) != expected {
		t.Errorf("Error expectation not met, want %v, get %v", expected, responseData)
	}
}

func TestEmptyTable(t *testing.T) {
	response := controllers.Response{
		Success: true,
		Data:    []models.Image{},
	}

	responseData, _ := json.Marshal(response)
	expected := `{"success":true,"data":[],"error":null}`
	if string(responseData) != expected {
		t.Errorf("Error expectation not met, want %v, get %v", expected, responseData)
	}
}

func TestGetImageById(t *testing.T) {
	id := "1"
	imageId, _ := strconv.ParseInt(id, 0, 0)
	imageData := models.Image{
		Id:        imageId,
		FileName:  "testing.png",
		FilePath:  "\\images\\testing.png",
		SourceUrl: "images/testing.png",
	}

	response := controllers.Response{
		Success: true,
		Data:    imageData,
	}

	responseData, _ := json.Marshal(response)
	expected := `{"success":true,"data":{"id":1,"file_name":"testing.png","file_path":"\\images\\testing.png","source_url":"images/testing.png"},"error":null}`
	if string(responseData) != expected {
		t.Errorf("Error expectation not met, want %v, get %v", expected, responseData)
	}
}

func TestGetImageNotExistent(t *testing.T) {
	response := controllers.Response{
		Success: false,
		Error: &controllers.ErrorResponse{
			Code: controllers.ERR_NOT_FOUND_CODE,
			Desc: controllers.ERR_NOT_FOUND_DESC,
		},
	}

	responseData, _ := json.Marshal(response)
	expected := `{"success":false,"data":null,"error":{"code":"RecordNotFound","desc":"Record not found in database"}}`
	if string(responseData) != expected {
		t.Errorf("Error expectation not met, want %v, get %v", expected, responseData)
	}
}

func TestFailedParsing(t *testing.T) {
	id := "'1'"
	_, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		response := controllers.Response{
			Success: false,
			Error: &controllers.ErrorResponse{
				Code: controllers.ERR_FAILED_PARSING_CODE,
				Desc: controllers.ERR_FAILED_PARSING_DESC,
			},
		}

		responseData, _ := json.Marshal(response)
		expected := `{"success":false,"data":null,"error":{"code":"FailedParsing","desc":"Failed to parse parameter"}}`
		if string(responseData) != expected {
			t.Errorf("Error expectation not met, want %v, get %v", expected, responseData)
		}
	} else {
		t.Error("Error expectation not met, if ParseInt does not have correct value, it should've got error")
	}
}

func TestDownloadImage(t *testing.T) {
	url := "https://i.imgur.com/ONsnEhy.jpeg"
	filename := path.Base(url)
	homeDir, _ := os.UserHomeDir()
	filePath := filepath.Join(homeDir, "Downloads", filename)
	err := controllers.SaveFile(url, filename, filePath)
	if err != nil {
		response := controllers.Response{
			Success: false,
			Error: &controllers.ErrorResponse{
				Code: controllers.ERR_FAILED_SAVE_FILE_CODE,
				Desc: controllers.ERR_FAILED_SAVE_FILE_DESC,
			},
		}

		responseData, _ := json.Marshal(response)
		expected := `{"success":false,"data":null,"error":{"code":"FailedSaveFile","desc":"Failed to save file to storage"}}`
		if string(responseData) != expected {
			t.Errorf("Error expectation not met, want %v, get %v", expected, responseData)
		}
	}

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		t.Errorf("File is not found in given file path")
	}

	imageData := models.Image{
		Id:        1,
		FileName:  "ONsnEhy.jpeg",
		FilePath:  "C:\\Users\\SIRCLO\\Downloads\\ONsnEhy.jpeg",
		SourceUrl: "https://i.imgur.com/ONsnEhy.jpeg",
	}

	response := controllers.Response{
		Success: true,
		Data:    imageData,
	}

	responseData, _ := json.Marshal(response)
	expected := `{"success":true,"data":{"id":1,"file_name":"ONsnEhy.jpeg","file_path":"C:\\Users\\SIRCLO\\Downloads\\ONsnEhy.jpeg","source_url":"https://i.imgur.com/ONsnEhy.jpeg"},"error":null}`
	if string(responseData) != expected {
		t.Errorf("Error expectation not met, want %v, get %v", expected, responseData)
	}
}

func TestFailedUnmarshal(t *testing.T) {
	var url controllers.Url
	jsonStr := []byte(`{\}`)
	err := json.Unmarshal(jsonStr, &url)

	if err != nil {
		response := controllers.Response{
			Success: false,
			Error: &controllers.ErrorResponse{
				Code: controllers.ERR_FAILED_UNMARSHAL_CODE,
				Desc: controllers.ERR_FAILED_UNMARSHAL_DESC,
			},
		}

		responseData, _ := json.Marshal(response)
		expected := `{"success":false,"data":null,"error":{"code":"FailedUnmarshal","desc":"Failed to unmarshal data"}}`
		if string(responseData) != expected {
			t.Errorf("Error expectation not met, want %v, get %v", expected, responseData)
		}
	} else {
		t.Error("Error expectation not met, if json string is not correct, it would fail unmarshal")
	}
}

func TestFailedSaveFile(t *testing.T) {
	url := ""
	filename := path.Base(url)
	homeDir, _ := os.UserHomeDir()
	filePath := filepath.Join(homeDir, "Downloads", filename)
	err := controllers.SaveFile(url, filename, filePath)
	if err != nil {
		response := controllers.Response{
			Success: false,
			Error: &controllers.ErrorResponse{
				Code: controllers.ERR_FAILED_SAVE_FILE_CODE,
				Desc: controllers.ERR_FAILED_SAVE_FILE_DESC,
			},
		}

		responseData, _ := json.Marshal(response)
		expected := `{"success":false,"data":null,"error":{"code":"FailedSaveFile","desc":"Failed to save file to storage"}}`
		if string(responseData) != expected {
			t.Errorf("Error expectation not met, want %v, get %v", expected, responseData)
		}
	} else {
		t.Error("Error expectation not met, if url is empty, it should've error on saving file to storage")
	}
}
