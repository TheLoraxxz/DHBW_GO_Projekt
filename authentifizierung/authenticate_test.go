package authentifizierung

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthenticateUserTrue(t *testing.T) {
	user := "admin"
	password := "admin"
	CreateUser(&user, &password)
	wahr, _ := AuthenticateUser(&user, &password)
	assert.Equal(t, true, wahr)
}

func TestAuthenticateUserFalse(t *testing.T) {
	user := "admin"
	password := "admin"
	CreateUser(&user, &password)
	passwordWrong := "user"
	wahr, _ := AuthenticateUser(&user, &passwordWrong)
	assert.Equal(t, false, wahr)
}

func TestCreateUser(t *testing.T) {
	user := "admin"
	password := "admin"
	assert.Equal(t, 0, len(users))
	CreateUser(&user, &password)
	assert.Equal(t, 1, len(users))
}
