package authentifizierung

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserData struct {
	user     string `json:"user"`
	password string `json:"password"`
}

// TODO: muss darauf warten ein threadsafe weg zu implementieren --> ich fixe das nach der Vorlesng: Channels und Go Routinen ^^
var users = make(map[string]string)

func AuthenticateUser(user *string, pasw *string) (correctPassw bool, newCookie string) {
	val, found := users[*user]
	for _, oneOfUsers := range users {
		err := bcrypt.CompareHashAndPassword([]byte(oneOfUsers), []byte(*pasw))
		if err == nil && found && strings.Compare(val, oneOfUsers) == 0 {
			bytes, hashError := bcrypt.GenerateFromPassword([]byte(oneOfUsers+val), 2)
			if hashError != nil {
				return false, ""
			}
			return true, *user + "|" + string(bytes)
		}
	}
	return false, ""
}

func CheckCookie(cookie *string) bool {
	cookieDeRef := *cookie
	username := cookieDeRef[:strings.Index(cookieDeRef, "|")]
	cookieString := cookieDeRef[strings.Index(cookieDeRef, "|")+1:]
	fmt.Println(username)
	fmt.Println(cookieString)
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
