package mail

import (
	"encoding/json"
	"github.com/rdavidnota/tiered-synchronous/source/commands/mail"
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
	message := domain.RequestSendMail{}
	err := json.Unmarshal(request, &message)
	utils.Check(err)

	if message.Base.Action == "send" {
		return sendMail(request)
	} else {
		return domain.Result{
			Code:    1,
			Message: "Action not found",
			Json:    "{}",
		}
	}
}

func sendMail(request []byte) domain.Result {

	requestSend := domain.RequestSendMail{}
	err := json.Unmarshal(request, &requestSend)
	utils.Check(err)

	mail.SendMail(requestSend.From, requestSend.To, requestSend.Message)

	return domain.Result{
		Code:    0,
		Message: "OK",
		Json:    "{}",
	}
}
