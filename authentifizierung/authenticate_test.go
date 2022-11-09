package authentifizierung

import (
	"testing"
)

func TestAuthenticateUserAdmin(t *testing.T) {
	var user string = "admin"
	var password string = "admin"
	AuthenticateUser(&user, &password)
}

func TestCreateUser(t *testing.T) {
	var user string = "admin"
	var password string = "admin"
	CreateUser(&user, &password)
}
