package main

import (
	"github.com/gorilla/mux"
	"github.com/rdavidnota/tiered-synchronous/source/controllers"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/documents", controllers.GetDocuments).Methods(http.MethodGet)
	router.HandleFunc("/documents/download/{id}", controllers.DownloadDocument).Methods(http.MethodGet)
	router.HandleFunc("/documents/{id}", controllers.GetDocument).Methods(http.MethodGet)
	router.HandleFunc("/documents", controllers.CreatedDocument).Queries("name", "{name}").Methods(http.MethodPost)
	router.HandleFunc("/documents/{id}", controllers.DelDocument).Methods(http.MethodDelete)



	log.Fatal(http.ListenAndServe(":9000", router))
}
