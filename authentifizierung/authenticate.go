package authentifizierung

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type UserJSON struct {
	User  string `json:"user"`
	Passw string `json:"password"`
}

// struct for authentication and start and endtime
type authentication struct {
	endTime time.Time
	user    string
}

type cookieStruct struct {
	cookies map[string]authentication
	lock    sync.RWMutex
}

var cookies = cookieStruct{
	cookies: map[string]authentication{},
	lock:    sync.RWMutex{},
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

// struct for user synchronisation
type usersync struct {
	lock  sync.RWMutex
	users map[string]string
	keys  map[string]authentication
}

var users = usersync{
	lock:  sync.RWMutex{},
	users: make(map[string]string, 5),
}

func DeleteOldCookies() {
	cookies.lock.Lock()
	defer cookies.lock.Unlock()
	newCookies := make(map[string]authentication)
	for key, val := range cookies.cookies {
		if val.endTime.After(time.Now()) {
			newCookies[key] = val
		}
	}
	cookies.cookies = newCookies
}

// AuthenticateUser --> called on login --> creates cookie
func AuthenticateUser(user *string, pasw *string) (correctPassw bool, newCookie string) {
	val, found := users.users[*user]
	if found {
		for _, oneOfUsers := range users.users {
			err := bcrypt.CompareHashAndPassword([]byte(oneOfUsers), []byte(*pasw))
			if err == nil && strings.Compare(val, oneOfUsers) == 0 {
				varbytes := make([]byte, 100)
				for i := range varbytes {
					varbytes[i] = letters[rand.Int63()%int64(len(letters))]
				}
				if err != nil {
					return false, ""
				}
				randomString := string(varbytes)
				cookies.lock.Lock()
				cookies.cookies[randomString] = authentication{
					user:    *user,
					endTime: time.Now().Add(15 * time.Minute),
				}
				correctPassw = true
				newCookie = *user + "|" + randomString
				cookies.lock.Unlock()
				return
			}
		}
	}
	return false, *pasw
}

// CheckCookie checks whether cookie is the right one
// returns false if it isnt the cookie to the user
func CheckCookie(cookie *string) (isAllowed bool, userName string) {
	//check the cookie to prevent any panics
	if len(*cookie) == 0 || strings.Index(*cookie, "|") == -1 || strings.Index(*cookie, "|") == len(*cookie)-1 {
		return false, ""
	}
	//get the username and cookie string from the cookie given
	username := (*cookie)[:strings.Index(*cookie, "|")]
	cookieString := (*cookie)[strings.Index(*cookie, "|")+1:]
	// if the username is as key in the map then it checks whether key is the same as the cookie
	if val, found := cookies.cookies[cookieString]; found == true && username == val.user && val.endTime.After(time.Now()) {
		cookies.lock.Lock()
		defer cookies.lock.Unlock()
		cookies.cookies[cookieString] = authentication{
			user:    userName,
			endTime: time.Now().Add(15 * time.Minute),
		}
		return true, username

	}
	return false, ""
}

// CreateUser
// checks if the user already exists and if not it creates a new user and hashes the password
func CreateUser(user *string, pasw *string) error {
	//check whether it contains $ or | --> | is not allowed because it is used in the cookie
	notAllowed := strings.ContainsAny(*user, "|$")
	if notAllowed {
		return errors.New("Username can't contain | or $")
	}
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

// ChangeUser
// Function checks whether the user exists, the old password given is the same
// if so it changes it to the new password
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

// LoadUserData
// function beeing called only once on startup
// loads up already existing users or if it is the first time
// it adds admin admin to the user
func LoadUserData(firstuser *string, firstPassword *string, path *string) (err error) {
	//loads working direcotry
	userLoaded := []UserJSON{}
	pathAbs := filepath.Join(*path, "data", "user", "user-data.json")
	file, err := os.Open(pathAbs)
	//if it cant open it because there are any problems it just creates a new user
	if err != nil {
		err := CreateUser(firstuser, firstPassword)
		if err != nil {
			return fmt.Errorf("Error on creating user %w", err)
		}
		return nil
	}
	// if it can open it it reads all bytes and tries to convert it to json
	// to then read it in and load it into the cache
	bytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Problem with Reading %w", err)
	}
	err = json.Unmarshal(bytes, &userLoaded)
	if err != nil {
		// if this is not the case and if there are any errors
		// it just creates per default what is given by the flags and then
		// continues on
		err := CreateUser(firstuser, firstPassword)
		if err != nil {
			return fmt.Errorf("Error on creating user %w", err)
		}
		// it closes the file to prevent any wrong read/right permissions
		err = file.Close()
		if err != nil {
			log.Fatal("coudnt close file")
		}
		//it removes the current file which it cant read
		err = os.Remove(pathAbs)
		if err != nil {
			return fmt.Errorf("Error on deleting file %w", err)
		}
		return nil
	}
	// if it works according to plan it loads everything in which is in the file
	for _, element := range userLoaded {
		users.users[element.User] = element.Passw
	}
	// should close it  last
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("coudnt close file")
		}
	}(file)
	return nil

}

// SaveUserData
// saves the data under ../data/users/user-data.json
func SaveUserData(path *string) error {
	//locks all data to prevent any problems
	user := []UserJSON{}
	users.lock.Lock()
	defer users.lock.Unlock()
	// iterate over and add to user json
	for key, elem := range users.users {
		user = append(user, UserJSON{User: key, Passw: elem})
	}
	// open path from basepath and save it as json
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
