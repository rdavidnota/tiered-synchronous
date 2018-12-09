package documents

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rdavidnota/tiered-synchronous/source/commands/utils"
	"github.com/rdavidnota/tiered-synchronous/source/commands/documents"
	"github.com/rdavidnota/tiered-synchronous/source/controllers/auth"
	"github.com/rdavidnota/tiered-synchronous/source/domain"
	"io"
	"log"
	"net/http"
	"strconv"
)

func GetDocuments(w http.ResponseWriter, r *http.Request) {

	if auth.Authorizer(w, r) {
		var docs []domain.Document
		docs = documents.ListFiles()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(docs)
	}
}

func GetDocument(w http.ResponseWriter, r *http.Request) {
	if auth.Authorizer(w, r) {
		var doc domain.Document
		params := mux.Vars(r)
		doc = documents.GetFile(params["id"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(doc)
	}
}

func DelDocument(w http.ResponseWriter, r *http.Request) {
	if auth.Authorizer(w, r) {
		var doc domain.Document
		params := mux.Vars(r)
		doc = documents.GetFile(params["id"])

		documents.DeleteFile(doc.Name)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(doc)
	}
}

func CreatedDocument(w http.ResponseWriter, r *http.Request) {

	if auth.Authorizer(w, r) {
		var filename = r.URL.Query().Get("name")
		fmt.Println(filename)
		file, _, err := r.FormFile("document")

		utils.Check(err)

		defer file.Close()

		var doc = documents.CreatedFileFromFile(filename, file)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(doc)
	}
}

func DownloadDocument(w http.ResponseWriter, r *http.Request) {
	if auth.Authorizer(w, r) {
		params := mux.Vars(r)

		file, filename := documents.GetFileById(params["id"])

		log.Println(filename)

		fileHeader := make([]byte, 512)
		file.Read(fileHeader)
		fileContentType := http.DetectContentType(fileHeader)

		fileStat, err := file.Stat()

		utils.Check(err)

		fileSize := strconv.FormatInt(fileStat.Size(), 10)

		w.Header().Set("Content-Disposition", "Attachment; filename="+filename)
		w.Header().Set("Content-Type", fileContentType)
		w.Header().Set("Content-Length", fileSize)

		defer file.Close()

		io.Copy(w, file)
	}
}
