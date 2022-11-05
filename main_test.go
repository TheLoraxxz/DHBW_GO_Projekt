package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"testing"
)

/*
*
Tests whether web server works even on Port 80
*/
func TestWebServerStart(t *testing.T) {
	req := httptest.NewRequest("GET", "localhost:80", nil)
	recorder := httptest.NewRecorder()

	RootHandler{}.ServeHTTP(recorder, req)
	resp := recorder.Result()
	assert.Equal(t, 200, resp.StatusCode, "the Server doesn't return an answer")
	_, err := io.ReadAll(resp.Body)
	assert.Equal(t, err, nil, "No Error in the Body")
}

/*
//test routine um ssl zertifikat zuu testen
func TestCertificateWorks(t *testing.T) {
	go main()

	time.Sleep(1 * time.Second)
	_, err := tls.Dial("tcp", "localhost:80", &tls.Config{InsecureSkipVerify: true})
	assert.Equal(t, err, nil)

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	if err := Server.Shutdown(ctx); err != nil {
		t.Error("coudnt6 stop server")
		t.Fail()
	}
	assert.Equal(t, true, true)

}
*/
