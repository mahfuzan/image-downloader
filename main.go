package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mahfuzan/image-downloader/controllers"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/download-image/", controllers.DownloadImage).Methods("POST")
	r.HandleFunc("/download-image/", controllers.GetImageList).Methods("GET")
	r.HandleFunc("/download-image/{id}", controllers.GetImageById).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("127.0.0.1:3306", r))
}
