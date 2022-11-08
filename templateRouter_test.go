package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"testing"
)

func TestReachAbleRoot(t *testing.T) {
	req := httptest.NewRequest("GET", "localhost:80", nil)
	recorder := httptest.NewRecorder()

	RootHandler{}.ServeHTTP(recorder, req)
	resp := recorder.Result()
	assert.Equal(t, 200, resp.StatusCode, "the Server doesn't return an answer")
	_, err := io.ReadAll(resp.Body)
	assert.Equal(t, err, nil, "No Error in the Body")
}
