package authentifizierung

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type UserJSON struct {
	User  string `json:"user"`
	Passw string `json:"password"`
}

// struct for user synchronisation
type usersync struct {
	lock  sync.RWMutex
	users map[string]string
}

var users = usersync{
	lock:  sync.RWMutex{},
	users: make(map[string]string, 5),
}

// AuthenticateUser --> called on login --> creates cookie
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

// TODO: if it doesnt find the slice it should return instand zero or error + test cases

// CheckCookie checks whether cookie is the right one/**
func CheckCookie(cookie *string) (isAllowed bool, userName string) {
	//get the username and cookie string from the cookie given
	username := (*cookie)[:strings.Index(*cookie, "|")]
	cookieString := (*cookie)[strings.Index(*cookie, "|")+1:]
	// if the username is as key in the map then it checks whether key is the same as the cookie
	if _, found := users.users[username]; found == true {
		//always checks whether the username given in the cookie is the same as the hashed value --> so one cant change the
		// username and get more access rights or different access rights
		err := bcrypt.CompareHashAndPassword([]byte(cookieString), []byte(username+users.users[username]))
		// if it is the same then it returns true and the username
		if err == nil {
			return true, username
		}

	}
	return false, ""
}

func CreateUser(user *string, pasw *string) error {
	//check whether it contains $ or | --> | is not allowed because it is used in the cookie
	notAllowed := strings.ContainsAny(*user, "|$")
	if notAllowed {
		return errors.New("Username can't contain | or $")
	}
	//TODO: setup new test case for this
	if len(*user) == 0 || len(*pasw) == 0 {
		return errors.New("Password and User can not be empty")
	}
	//lock user because now we are looking into the user and check whether the username is the same
	users.lock.Lock()
	// on end unlock the user
	defer users.lock.Unlock()
	//if user is found then it returns
	_, found := users.users[*user]
	if found {
		return errors.New("User already created")
	}
	// now do the performance costing hashing algorithms and check whether error is nil
	bytes, err := bcrypt.GenerateFromPassword([]byte(*pasw), 14)
	if err != nil {
		return err
	}
	// actually create the user and then return nil to show that everything worked
	users.users[*user] = string(bytes)
	return nil
}
func ChangeUser(user *string, oldPassw *string, newPassw *string) (newCookie string, err error) {
	val, found := users.users[*user]
	if found {
		err := bcrypt.CompareHashAndPassword([]byte(val), []byte(*oldPassw))
		if err != nil {
			return "", err
		}
		// generate new password
		newHash, errorHash := bcrypt.GenerateFromPassword([]byte(*newPassw), 14)
		if errorHash != nil {
			return "", errorHash
		}
		users.lock.Lock()
		defer users.lock.Unlock()
		users.users[*user] = string(newHash)
		_, cookie := AuthenticateUser(user, newPassw)
		return cookie, nil
	}
	return "", errors.New("Wrong User")
}

func LoadUserData(firstuser *string, firstPassword *string, path *string) (err error) {
	ex, err := os.Getwd()
	fmt.Println(ex)
	userLoaded := []UserJSON{}
	pathAbs := filepath.Join(*path, "data", "user", "user-data.json")
	file, err := os.Open(pathAbs)
	if err != nil {
		err := CreateUser(firstuser, firstPassword)
		if err != nil {
			return fmt.Errorf("Error on creating user %w", err)
		}
		return nil
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Problem with Reading %w", err)
	}
	err = json.Unmarshal(bytes, &userLoaded)
	if err != nil {
		err := CreateUser(firstuser, firstPassword)
		if err != nil {
			return fmt.Errorf("Error on creating user %w", err)
		}
		err = file.Close()
		if err != nil {
			log.Fatal("coudnt close file")
		}
		err = os.Remove(pathAbs)
		if err != nil {
			return fmt.Errorf("Error on deleting file %w", err)
		}
		return nil
	}
	for _, element := range userLoaded {
		users.users[element.User] = element.Passw
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("coudnt close file")
		}
	}(file)
	return nil

}

func SaveUserData(path *string) error {
	users.lock.Lock()
	user := []UserJSON{}
	defer users.lock.Unlock()
	for key, elem := range users.users {
		user = append(user, UserJSON{User: key, Passw: elem})
	}
	pathAbs := filepath.Join(*path, "data", "user", "user-data.json")
	file, err := json.MarshalIndent(user, "", "")
	if err != nil {
		return fmt.Errorf("Error on creating json file %w", err)
	}
	err = os.WriteFile(pathAbs, file, 0644)
	if err != nil {
		return fmt.Errorf("Error on writing to Json file %w", err)
	}

	return err
}
