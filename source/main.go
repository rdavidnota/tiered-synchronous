package main

import (
	"github.com/gorilla/mux"
	"github.com/rdavidnota/tiered-synchronous/source/controllers/documents"
	"github.com/rdavidnota/tiered-synchronous/source/controllers/mail"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/documents", documents.GetDocuments).Methods(http.MethodGet)
	router.HandleFunc("/documents/download/{id}", documents.DownloadDocument).Methods(http.MethodGet)
	router.HandleFunc("/documents/{id}", documents.GetDocument).Methods(http.MethodGet)
	router.HandleFunc("/documents", documents.CreatedDocument).Queries("name", "{name}").Methods(http.MethodPost)
	router.HandleFunc("/documents/{id}", documents.DelDocument).Methods(http.MethodDelete)

	router.HandleFunc("/mail/send", mail.SendMail).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":9000", router))
}
