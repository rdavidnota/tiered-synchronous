package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type Document struct {
	ID   string
	Name string
	Size int64
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/documents", getDocuments).Methods(http.MethodGet)
	router.HandleFunc("/documents/{id}", getDocument).Methods(http.MethodGet)
	router.HandleFunc("/documents", createdDocument).Queries("name", "{name}").Methods(http.MethodPost)
	router.HandleFunc("/documents/{id}", delDocument).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":9000", router))
}

func getDocuments(w http.ResponseWriter, r *http.Request) {
	var docs []Document
	docs = listFiles()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(docs)
}

func getDocument(w http.ResponseWriter, r *http.Request) {
	var doc Document
	params := mux.Vars(r)
	doc = getFile(params["id"])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(doc)
}

func delDocument(w http.ResponseWriter, r *http.Request) {
	var doc Document
	params := mux.Vars(r)
	doc = getFile(params["id"])

	deleteFile(doc.Name)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(doc)
}

func createdDocument(w http.ResponseWriter, r *http.Request) {

	var filename = r.URL.Query().Get("name")
	fmt.Println(filename)
	file, _, err := r.FormFile("document")

	check(err)

	defer file.Close()

	var doc = createdFile(filename, file)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(doc)
}

func createdFile(filename string, file multipart.File) Document{
	f, err := os.OpenFile(pathFolder+"\\"+filename, os.O_WRONLY|os.O_CREATE, 0666)

	check(err)
	defer f.Close()

	io.Copy(f, file)

	return getFile(calculateChecksum(filename))
}

func deleteFile(filename string) {
	if len(filename) > 0 {
		error := os.Remove(pathFolder + "\\" + filename)
		check(error)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const pathFolder = "C:\\Users\\UTI01\\Downloads\\Programs"

func getFile(id string) Document {
	var result Document
	var docs []Document
	docs = listFiles()

	for _, doc := range docs {
		if doc.ID == id {
			result = doc
		}
	}

	return result
}

func listFiles() []Document {
	files, err := ioutil.ReadDir(pathFolder)
	if err != nil {
		log.Fatal(err)
	}

	var docs []Document

	for _, f := range files {
		if !f.IsDir() {
			docs = append(docs, Document{ID: calculateChecksum(f.Name()), Name: f.Name(), Size: f.Size()})
		}
	}

	return docs
}

func calculateChecksum(filename string) string {
	var checksum = ""
	file, err := os.Open(pathFolder + "\\" + filename)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		panic(err)
	}

	checksum = hex.EncodeToString(hash.Sum(nil))

	return checksum
}
