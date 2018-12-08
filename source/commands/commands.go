package commands

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/rdavidnota/tiered-synchronous/source/domain"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
)

func CreatedFile(filename string, file multipart.File) domain.Document {
	f, err := os.OpenFile(domain.PathFolder+"\\"+filename, os.O_WRONLY|os.O_CREATE, 0666)

	Check(err)
	defer f.Close()

	io.Copy(f, file)

	return GetFile(CalculateChecksum(filename))
}

func DeleteFile(filename string) {
	if len(filename) > 0 {
		err := os.Remove(domain.PathFolder + "\\" + filename)
		Check(err)
	}
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func GetFileById(id string) (*os.File, string) {
	files, err := ioutil.ReadDir(domain.PathFolder)
	Check(err)

	var doc domain.Document

	for _, f := range files {
		if !f.IsDir() {
			checksum := CalculateChecksum(f.Name())

			if checksum == id {
				doc = domain.Document{Name: f.Name(), ID: checksum, Size: f.Size()}
				break
			}
		}
	}

	if doc.ID == "" {
		return nil, ""
	}

	file, err := os.Open(domain.PathFolder + "\\" + doc.Name)

	Check(err)

	return file, file.Name()
}

func GetFile(id string) domain.Document {
	var result domain.Document
	var docs []domain.Document
	docs = ListFiles()

	for _, doc := range docs {
		if doc.ID == id {
			result = doc
		}
	}

	return result
}

func ListFiles() []domain.Document {
	files, err := ioutil.ReadDir(domain.PathFolder)

	Check(err)

	var docs []domain.Document

	for _, f := range files {
		if !f.IsDir() {
			docs = append(docs, domain.Document{ID: CalculateChecksum(f.Name()), Name: f.Name(), Size: f.Size()})
		}
	}

	return docs
}

func CalculateChecksum(filename string) string {
	var checksum = ""
	file, err := os.Open(domain.PathFolder + "\\" + filename)

	Check(err)

	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	Check(err)

	checksum = hex.EncodeToString(hash.Sum(nil))

	return checksum
}
