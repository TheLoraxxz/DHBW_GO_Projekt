package authentifizierung

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserData struct {
	user     string `json:"user"`
	password string `json:"password"`
}

// TODO: muss darauf warten ein threadsafe weg zu implementieren --> ich fixe das nach der Vorlesng: Channels und Go Routinen ^^
var users = []UserData{}

func AuthenticateUser(user *string, pasw *string) (correctPassw bool, newCookie string) {
	for _, oneOfUsers := range users {
		err := bcrypt.CompareHashAndPassword([]byte(oneOfUsers.password), []byte(*pasw))
		if err == nil && strings.Compare(*user, oneOfUsers.user) == 0 {
			bytes, _ := bcrypt.GenerateFromPassword([]byte(oneOfUsers.password), 2)
			return true, string(bytes)
		}
	}
	return false, ""
}

func CheckCookie() {

}

func CreateUser(user *string, pasw *string) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*pasw), 14)
	if err != nil {
		fmt.Println(err)
	}
	newUser := UserData{
		user:     *user,
		password: string(bytes),
	}
	users = append(users, newUser)
}
