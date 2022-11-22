package main

import (
	"crypto/tls"
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
	_, err := tls.Dial("tcp", "localhost:80", &tls.Config{InsecureSkipVerify: true})
	assert.Equal(t, err, nil)

	err = Server.Close()
	if err != nil {
		assert.Failf(t, "Coudn't Close server", "")
	}

}
