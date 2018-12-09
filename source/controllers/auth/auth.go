package auth

import "net/http"

func Authorizer(w http.ResponseWriter, r *http.Request) bool {
	username, password, _ := r.BasicAuth()

	if username != "rnota" || password != "mercado.nota" {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	return true
}
