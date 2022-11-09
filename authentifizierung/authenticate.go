package authentifizierung

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserData struct {
	user     string `json:"user"`
	password string `json:"password"`
}

var users = make([]UserData, 5)

func AuthenticateUser(user *string, pasw *string) {
	fmt.Println(*user)
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
