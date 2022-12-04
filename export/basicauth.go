package export

/*
BasicAuth für Download der Ical bereitzustellen --> weites gehend entnommen aus der Vorlesung
*/

//Mat-Nr. 8689159
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
		isUser, cookieText := authenticator.Authenticate(user, pswd)
		if isUser == true {
			// wenn user authentifiziert ist, dann wird ein cookie erstellt um später auf den Username zugreifen zu können
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

// CheckUserValid prüft, ob zugriff Valide
func CheckUserValid(user, pswd string) (bool, string) {
	check, cookie := authentifizierung.AuthenticateUser(&user, &pswd)
	return check, cookie
}
