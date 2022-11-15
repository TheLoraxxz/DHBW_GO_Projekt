package authentifizierung

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"testing"
)

// tests pon create user
func TestCreateUser(t *testing.T) {
	user := "admin"
	password := "admin"
	assert.Equal(t, 0, len(users.users))
	CreateUser(&user, &password)
	assert.Equal(t, 1, len(users.users))
}

// tests that on authentication it returns a true user
func TestAuthenticateUserTrue(t *testing.T) {
	users.users = make(map[string]string)
	user := "admin"
	password := "admin"
	CreateUser(&user, &password)
	wahr, _ := AuthenticateUser(&user, &password)
	assert.Equal(t, true, wahr)
}

// checks that if it fails that it returns a wrong user
func TestAuthenticateUserFalse(t *testing.T) {
	users.users = make(map[string]string)
	user := "admin"
	password := "admin"
	CreateUser(&user, &password)
	passwordWrong := "user"
	wahr, cookie := AuthenticateUser(&user, &passwordWrong)
	assert.Equal(t, false, wahr)
	assert.NotEqual(t, nil, cookie)
}

// test that the cookie which is given back is the right one
func TestAuthenticateUserCookieMngmt(t *testing.T) {
	users.users = make(map[string]string)
	user := "admin"
	password := "admin"
	// nutzer erstellen
	CreateUser(&user, &password)
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
func TestCheckCookieTrue(t *testing.T) {
	users.users = make(map[string]string)
	user := "admin"
	password := "admin"
	CreateUser(&user, &password)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin"+users.users["admin"]), 2)
	cookie := "admin|" + string(hashedPassword)
	isAllowed, username := CheckCookie(&cookie)
	assert.Equal(t, true, isAllowed)
	assert.Equal(t, "admin", username)
}

// checks that authenticate user and check cookie work with each other
func TestCheckCookieAndAuthenticateUser(t *testing.T) {
	users.users = make(map[string]string)
	user := "admin"
	password := "admin"
	CreateUser(&user, &password)
	_, cookie := AuthenticateUser(&user, &password)
	isAllowed, _ := CheckCookie(&cookie)
	assert.Equal(t, true, isAllowed)
}
