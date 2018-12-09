package mail

import (
	"github.com/rdavidnota/tiered-synchronous/source/commands/mail"
	"github.com/rdavidnota/tiered-synchronous/source/controllers/auth"
	"net/http"
	"strings"
)

func SendMail(w http.ResponseWriter, r *http.Request) {
	if auth.Authorizer(w, r) {
		from := r.FormValue("from")
		to := strings.Split(r.FormValue("to"), ",")
		msg := r.FormValue("message")

		mail.SendMail(from, to, msg)

		w.WriteHeader(http.StatusOK)
	}
}
