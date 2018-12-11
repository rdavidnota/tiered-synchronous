package documents

import (
	"encoding/json"
	"fmt"
	"github.com/rdavidnota/tiered-synchronous/source/commands/documents"
	"github.com/rdavidnota/tiered-synchronous/source/commands/utils"
	"github.com/rdavidnota/tiered-synchronous/source/domain"
	"log"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}


func Analyze(request []byte) domain.Result {
	fmt.Println(string(request))
	message := domain.RequestListDocument{}
	err := json.Unmarshal(request, &message)
	utils.Check(err)
	fmt.Println(message.Base.Action)
	if message.Base.Action == "create" {
		return create(request)
	} else if message.Base.Action == "remove" {
		return remove(request)
	} else if message.Base.Action == "get" {
		return get(request)
	} else if message.Base.Action == "list" {
		return list()
	} else {
		return domain.Result{
			Code:    1,
			Message: "Action not found",
			Json:    "{}",
		}
	}
}

func get(request []byte) domain.Result {
	requestGet := domain.RequestGetDocument{}
	err := json.Unmarshal(request, &requestGet)
	utils.Check(err)

	document := documents.GetFile(requestGet.ID)
	jsonResult, _ := json.Marshal(document)

	return domain.Result{
		Code:    0,
		Message: "OK",
		Json:    string(jsonResult),
	}
}

func create(request []byte) domain.Result {
	requestCreate := domain.RequestCreateDocument{}
	err := json.Unmarshal(request, &requestCreate)
	utils.Check(err)

	documents.CreatedFileFromBytes(requestCreate.Name, requestCreate.Content)

	return list()
}

func remove(request []byte) domain.Result {
	requestDelete := domain.RequestDeleteDocument{}
	err := json.Unmarshal(request, &requestDelete)
	utils.Check(err)
	documents.DeleteFile(requestDelete.ID)

	return list()
}

func list() domain.Result {
	listDocuments := documents.ListFiles()
	jsonResult, err := json.Marshal(listDocuments)

	utils.Check(err)

	return domain.Result{
		Code:    0,
		Message: "OK",
		Json:    string(jsonResult),
	}
}
