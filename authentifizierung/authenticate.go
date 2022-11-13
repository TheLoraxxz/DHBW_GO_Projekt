package authentifizierung

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserData struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

// TODO: muss darauf warten ein threadsafe weg zu implementieren --> ich fixe das nach der Vorlesng: Channels und Go Routinen ^^
var users = make(map[string]string)

func AuthenticateUser(user *string, pasw *string) (correctPassw bool, newCookie string) {
	val, found := users[*user]
	if found {
		for _, oneOfUsers := range users {
			err := bcrypt.CompareHashAndPassword([]byte(oneOfUsers), []byte(*pasw))
			if err == nil && strings.Compare(val, oneOfUsers) == 0 {
				bytes, hashError := bcrypt.GenerateFromPassword([]byte(*user+oneOfUsers), 2)
				if hashError != nil {
					return false, ""
				}
				return true, *user + "|" + string(bytes)
			}
		}
	}
	return false, *pasw
}

func CheckCookie(cookie *string) bool {
	cookieDeRef := *cookie
	username := cookieDeRef[:strings.Index(cookieDeRef, "|")]
	cookieString := cookieDeRef[strings.Index(cookieDeRef, "|")+1:]
	if _, found := users[username]; found == true {
		err := bcrypt.CompareHashAndPassword([]byte(cookieString), []byte(username+users[username]))
		if err == nil {
			return true
		}

	}
	return false
}

func CreateUser(user *string, pasw *string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*pasw), 14)
	if err != nil {
		return errors.New("Fehlschlag des Hashings")
	}
	notAllowed := strings.ContainsAny(*user, "|$")
	if notAllowed {
		return errors.New("Username darf keine Sonderzeichen enthalten")
	}
	users[*user] = string(bytes)
	return nil
}
