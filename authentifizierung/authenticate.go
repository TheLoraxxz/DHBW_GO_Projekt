package authentifizierung

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"sync"
)

type UserData struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type usersync struct {
	lock  sync.Mutex
	users map[string]string
}

var users = usersync{
	lock:  sync.Mutex{},
	users: make(map[string]string, 5),
}

func AuthenticateUser(user *string, pasw *string) (correctPassw bool, newCookie string) {
	val, found := users.users[*user]
	if found {
		for _, oneOfUsers := range users.users {
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

func CheckCookie(cookie *string) (isAllowed bool, userName string) {
	cookieDeRef := *cookie
	username := cookieDeRef[:strings.Index(cookieDeRef, "|")]
	cookieString := cookieDeRef[strings.Index(cookieDeRef, "|")+1:]
	if _, found := users.users[username]; found == true {
		err := bcrypt.CompareHashAndPassword([]byte(cookieString), []byte(username+users.users[username]))
		if err == nil {
			return true, username
		}

	}
	return false, ""
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
	users.lock.Lock()
	defer users.lock.Unlock()
	_, found := users.users[*user]
	if found {
		return errors.New("User already created")
	}
	users.users[*user] = string(bytes)
	return nil
}
func ChangeUser(olduser *string, newuser *string, oldPassw *string, newPassw *string) (newCookie string, err error) {

}
