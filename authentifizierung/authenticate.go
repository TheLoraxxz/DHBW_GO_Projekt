package authentifizierung

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"sync"
)

// UserData used for dateisystem
type UserData struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

// struct for user synchronisation
type usersync struct {
	lock  sync.Mutex
	users map[string]string
}

var users = usersync{
	lock:  sync.Mutex{},
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

	}
	return "", errors.New("Wrong User")
}
