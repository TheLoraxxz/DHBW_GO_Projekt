package authentifizierung

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// tests pon create user
func TestCreateUser(t *testing.T) {
	user := "admin"
	password := "admin"
	assert.Equal(t, 0, len(users.users))
	_ = CreateUser(&user, &password)
	assert.Equal(t, 1, len(users.users))
}

// tests that on authentication it returns a true user
func TestAuthenticateUser_True(t *testing.T) {
	users.users = make(map[string]string)
	user := "admin"
	password := "admin"
	_ = CreateUser(&user, &password)
	wahr, _ := AuthenticateUser(&user, &password)
	assert.Equal(t, true, wahr)
}

// checks that if it fails that it returns a wrong user
func TestAuthenticateUser_False(t *testing.T) {
	users.users = make(map[string]string)
	user := "admin"
	password := "admin"
	_ = CreateUser(&user, &password)
	passwordWrong := "user"
	wahr, cookie := AuthenticateUser(&user, &passwordWrong)
	assert.Equal(t, false, wahr)
	assert.NotEqual(t, nil, cookie)
}

// test that the cookie which is given back is the right one
func TestAuthenticateUser_CookieMngmt(t *testing.T) {
	users.users = make(map[string]string)
	user := "admin"
	password := "admin"
	// nutzer erstellen
	_ = CreateUser(&user, &password)
	// den coookie zurückholen aufgebaut wie:
	_, cookie := AuthenticateUser(&user, &password)
	// auslesen des cookies und des uernames
	username := cookie[:strings.Index(cookie, "|")]
	cookie = cookie[strings.Index(cookie, "|")+1:]
	// schauen, dass der neue hash richtig generiert ist
	isSame := bcrypt.CompareHashAndPassword([]byte(cookie), []byte("admin"+users.users["admin"]))
	assert.Equal(t, nil, isSame)
	assert.Equal(t, user, username)
}

// tests that cookie check is true on right input
func TestCheckCookie_True(t *testing.T) {
	users.users = make(map[string]string)
	user := "admin"
	password := "admin"
	_ = CreateUser(&user, &password)
	//cookie is manually created
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin"+users.users["admin"]), 2)
	cookie := "admin|" + string(hashedPassword)
	//checked whether cookie is equal
	isAllowed, username := CheckCookie(&cookie)
	assert.Equal(t, true, isAllowed)
	assert.Equal(t, "admin", username)
}

// checks that authenticate user and check cookie work with each other
func TestCheckCookieAndAuthenticateUser(t *testing.T) {
	users.users = make(map[string]string)
	user := "admin"
	password := "admin"
	_ = CreateUser(&user, &password)
	_, cookie := AuthenticateUser(&user, &password)
	isAllowed, _ := CheckCookie(&cookie)
	assert.Equal(t, true, isAllowed)
}

// tests if user already exists whether it returns an error
func TestCreateUser_AlreadyExists(t *testing.T) {
	users.users = make(map[string]string)
	user := "admin"
	password := "admin"
	_ = CreateUser(&user, &password)
	err := CreateUser(&user, &password)
	assert.NotEqual(t, nil, err)
}

// tests if multiple users are created at the same time whether it saves them all or whether there are overlapping
func TestCreateUser_MultipleUserAtOnce(t *testing.T) {
	users.users = make(map[string]string)
	//create sync group for making done
	var wg sync.WaitGroup
	// function that checks whether the function has run through successfully
	addUser := func(name *string, passw *string) {
		err := CreateUser(name, passw)
		assert.Equal(t, nil, err)
		wg.Done()
	}
	wg.Add(20)
	for i := 0; i < 20; i++ {
		test := "user" + strconv.Itoa(i)
		go addUser(&test, &test)
	}
	wg.Wait()
	//check at the end that the length of user is 20 and the information is secured
	assert.Equal(t, 20, len(users.users))
}

func TestCreateUser_NoSpecialChractersInUsername(t *testing.T) {
	users.users = make(map[string]string)
	username := "test|"
	password := "test"
	err := CreateUser(&username, &password)
	assert.NotEqual(t, nil, err, "Create User works but shoudnt")
	assert.Equal(t, 0, len(users.users))
}

func TestCreateUser_SpecialCharacktersInPasswordAllowed(t *testing.T) {
	users.users = make(map[string]string)
	username := "test"
	password := "test|$"
	err := CreateUser(&username, &password)
	assert.Equal(t, nil, err, "should work")
	assert.Equal(t, 1, len(users.users))
}

// should work correctly
func TestChangeUser(t *testing.T) {
	users.users = make(map[string]string)
	username := "test"
	password := "test"
	_ = CreateUser(&username, &password)
	newPassword := "bla"
	//change password
	newCookie, err := ChangeUser(&username, &password, &newPassword)
	assert.Equal(t, err, nil)
	//check  that cookie is allowed --> the right cookie should be given back
	isCorrect, _ := AuthenticateUser(&username, &newPassword)
	assert.Equal(t, true, isCorrect)

	isCorrect, _ = CheckCookie(&newCookie)
	assert.Equal(t, true, isCorrect)
}

// should return an error if the wrong user is given
func TestChangeUser_WrongUser(t *testing.T) {
	users.users = make(map[string]string)
	username := "test"
	password := "test"
	_ = CreateUser(&username, &password)
	newPassword := "bla"
	_, err := ChangeUser(&newPassword, &password, &newPassword)
	assert.NotEqual(t, nil, err)
}

// should return an error if it is the wrong old password
func TestChangeUser_WrongOldPassword(t *testing.T) {
	users.users = make(map[string]string)
	username := "test"
	password := "test"
	_ = CreateUser(&username, &password)
	newPassword := "bla"
	_, err := ChangeUser(&username, &newPassword, &password)
	assert.NotEqual(t, nil, err)
}
