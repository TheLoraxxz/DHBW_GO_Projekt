package main

import (
	"DHBW_GO_Projekt/authentifizierung"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

/*
*
Tests whether web server works even on Port 80
*/

// test routine um ssl zertifikat zuu testen

func TestCertificateWorks(t *testing.T) {
	go main()
	time.Sleep(1 * time.Second)
	//check that the server has been created
	assert.NotEmpty(t, Server)
	// check that the admin user is created
	user := "admin"
	createdUser, _ := authentifizierung.AuthenticateUser(&user, &user)
	assert.Equal(t, true, createdUser)
}
