package models

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/mahfuzan/image-downloader/config"
)

var schema = `
CREATE TABLE IF NOT EXISTS images (
	id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	file_name text,
	file_path text,
	source_url text
);
`

type Image struct {
	Id        int64  `db:"id" json:"id"`
	FileName  string `db:"file_name" json:"file_name"`
	FilePath  string `db:"file_path" json:"file_path"`
	SourceUrl string `db:"source_url" json:"source_url"`
}

var db *sqlx.DB

func init() {
	config.Connect()
	db = config.GetDb()
	_, err := db.Query("SELECT * FROM images")
	if err != nil {
		_, err = db.Exec(schema)
		if err != nil {
			log.Println(err)
		}
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
	}
}

func SaveToDatabase(url string, filename string, filePath string) (Image, error) {
	image := Image{
		FileName:  filename,
		FilePath:  filePath,
		SourceUrl: url,
	}

	result, err := db.NamedExec("INSERT INTO images (file_name, file_path, source_url) VALUES (:file_name, :file_path, :source_url)", &image)
	if err != nil {
		return Image{}, err
	}

	image.Id, err = result.LastInsertId()
	if err != nil {
		return Image{}, err
	}

	return image, nil
}

func GetList() ([]Image, error) {
	imageList := []Image{}
	err := db.Select(&imageList, "SELECT * FROM images")
	if err != nil {
		return []Image{}, err
	}

	return imageList, nil
}

func GetById(id int64) (Image, error) {
	image := Image{}
	err := db.Get(&image, "SELECT * FROM images WHERE id = ?", id)
	if err != nil {
		return image, err
	}

	return image, nil
}
