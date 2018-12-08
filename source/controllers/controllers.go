package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rdavidnota/tiered-synchronous/source/commands"
	"github.com/rdavidnota/tiered-synchronous/source/domain"
	"io"
	"log"
	"net/http"
	"strconv"
)

func Authorizer(w http.ResponseWriter, r *http.Request) bool {
	username, password, _ := r.BasicAuth()

	if username != "rnota" || password != "mercado.nota" {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	return true
}

func GetDocuments(w http.ResponseWriter, r *http.Request) {

	if Authorizer(w, r) {
		var docs []domain.Document
		docs = commands.ListFiles()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(docs)
	}
}

func GetDocument(w http.ResponseWriter, r *http.Request) {
	if Authorizer(w, r) {
		var doc domain.Document
		params := mux.Vars(r)
		doc = commands.GetFile(params["id"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(doc)
	}
}

func DelDocument(w http.ResponseWriter, r *http.Request) {
	if Authorizer(w, r) {
		var doc domain.Document
		params := mux.Vars(r)
		doc = commands.GetFile(params["id"])

		commands.DeleteFile(doc.Name)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(doc)
	}
}

func CreatedDocument(w http.ResponseWriter, r *http.Request) {

	if Authorizer(w, r) {
		var filename = r.URL.Query().Get("name")
		fmt.Println(filename)
		file, _, err := r.FormFile("document")

		commands.Check(err)

		defer file.Close()

		var doc = commands.CreatedFile(filename, file)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(doc)
	}
}

func DownloadDocument(w http.ResponseWriter, r *http.Request) {
	if Authorizer(w, r) {
		params := mux.Vars(r)

		file, filename := commands.GetFileById(params["id"])

		log.Println(filename)

		fileHeader := make([]byte, 512)
		file.Read(fileHeader)
		fileContentType := http.DetectContentType(fileHeader)

		fileStat,err := file.Stat()

		commands.Check(err)

		fileSize := strconv.FormatInt(fileStat.Size(), 10)

		w.Header().Set("Content-Disposition", "Attachment; filename="+filename)
		w.Header().Set("Content-Type", fileContentType)
		w.Header().Set("Content-Length", fileSize)

		defer file.Close()

		io.Copy(w, file)
	}
}
