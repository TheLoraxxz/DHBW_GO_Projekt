package export

import (
	"net/http"
)

type Authenticator interface {
	Authenticate(user, password string) bool
}
type AuthenticatorFunc func(user, password string) bool

func (af AuthenticatorFunc) Authenticate(user, password string) bool {
	return af(user, password)
}

func WrapperAuth(authenticator Authenticator, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pswd, ok := r.BasicAuth()
		if ok && authenticator.Authenticate(user, pswd) {
			handler(w, r)
		} else {
			w.Header().Set("WWW-Authenticate",
				"Basic realm=\"My Simple Server\"")
			http.Error(w,
				http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized)
		}
	}
}

// CheckUserValid ToDo User Prüfung umsetzen
func CheckUserValid(user, pswd string) bool { //prüft, ob zugriff Valide
	return true
}
