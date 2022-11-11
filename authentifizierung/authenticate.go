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

func AuthenticateUser(user *string, pasw *string) (correctPassw bool, authUser string) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*pasw), 14)
	if err == nil {
		passwordHashed := string(bytes)
		for _, oneOfUsers := range users {
			if strings.Compare(oneOfUsers.password, passwordHashed) == 0 && strings.Compare(*user, oneOfUsers.user) == 0 {
				return true, oneOfUsers.user
			}
		}
	}
	return false, ""
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
