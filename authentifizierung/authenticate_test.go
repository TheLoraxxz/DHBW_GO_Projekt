package authentifizierung

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	user := "admin"
	password := "admin"
	assert.Equal(t, 0, len(users))
	CreateUser(&user, &password)
	assert.Equal(t, 1, len(users))
}

func TestAuthenticateUserTrue(t *testing.T) {
	users = make(map[string]string)
	user := "admin"
	password := "admin"
	CreateUser(&user, &password)
	wahr, _ := AuthenticateUser(&user, &password)
	assert.Equal(t, true, wahr)
}

func TestAuthenticateUserFalse(t *testing.T) {
	users = make(map[string]string)
	user := "admin"
	password := "admin"
	CreateUser(&user, &password)
	passwordWrong := "user"
	wahr, cookie := AuthenticateUser(&user, &passwordWrong)
	assert.Equal(t, false, wahr)
	assert.NotEqual(t, nil, cookie)
}

func TestAuthenticateUserCookieMngmt(t *testing.T) {
	users = make(map[string]string)
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
	//TODO: finish this shit that it works
	isSame := bcrypt.CompareHashAndPassword([]byte(cookie), []byte(users["admin"]))
	assert.Equal(t, nil, isSame)
	assert.Equal(t, user, username)
}

func TestCheckCookie(t *testing.T) {

}
