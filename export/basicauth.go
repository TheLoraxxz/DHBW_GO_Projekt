package export

import (
	"DHBW_GO_Projekt/authentifizierung"
	"net/http"
)

type Authenticator interface {
	Authenticate(user, password string) (bool, string)
}
type AuthenticatorFunc func(user, password string) (bool, string)

func (af AuthenticatorFunc) Authenticate(user, password string) (bool, string) {
	return af(user, password)
}

func WrapperAuth(authenticator Authenticator, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pswd, ok := r.BasicAuth()
		isUser, cookieText := authentifizierung.AuthenticateUser(&user, &pswd)
		if isUser == true {
			// wenn user authentifiziert ist dann wird cookie erstellt und
			cookie := &http.Cookie{
				Name:     "Download-Kalender",
				Value:    cookieText,
				Path:     "/",
				MaxAge:   3600,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			}
			r.AddCookie(cookie)
		}
		if ok && isUser {
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

// CheckUserValid User Prüfung
func CheckUserValid(user, pswd string) (bool, string) { //prüft, ob zugriff Valide
	check, cookie := authentifizierung.AuthenticateUser(&user, &pswd)
	return check, cookie
}
