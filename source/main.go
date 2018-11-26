package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
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
	router.HandleFunc("/documents", getDocuments).Methods("GET")
	router.HandleFunc("/documents/{id}", getDocument).Methods("GET")
	log.Fatal(http.ListenAndServe(":9000", router))
}

func getDocuments(w http.ResponseWriter, r *http.Request) {
	var docs []Document
	docs = listFiles()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}

func getDocument(w http.ResponseWriter, r *http.Request) {
	var doc Document
	params := mux.Vars(r)
	doc = getFile(params["id"])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

const pathFolder = "C:\\Users\\rdavi\\Downloads\\Programs"

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
