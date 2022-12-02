package main

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestErrorSite_ServeHttp_EmptyGetter(t *testing.T) {
	req := httptest.NewRequest("GET", "localhost:443/error", nil)
	rec := httptest.NewRecorder()

	ErrorSite_ServeHttp(rec, req)
	assert.Equal(t, 200, rec.Result().StatusCode)
	assert.Equal(t, true, strings.Contains(rec.Body.String(), "Interner Error Problem"))
}

func TestErrorSite_ServeHttp_Correct(t *testing.T) {
	req := httptest.NewRequest("GET", "localhost:443/error?type=internal&link="+url.QueryEscape("/"), nil)
	rec := httptest.NewRecorder()
	// execute
	ErrorSite_ServeHttp(rec, req)
	// check that the link is correctly set and the right error is given
	assert.Equal(t, 200, rec.Result().StatusCode)
	assert.Equal(t, true, strings.Contains(rec.Body.String(), "Interner Server error"))
	assert.Equal(t, true, strings.Contains(rec.Body.String(), "<a class=\"btn btn-primary\" href=\"https://"+req.Host+"/\">Zurück</a>"))
}

func TestErrorSite_ServeHttp_EmptyLink(t *testing.T) {
	//link is not set
	req := httptest.NewRequest("GET", "localhost:443/error?type=internal", nil)
	rec := httptest.NewRecorder()
	// execute
	ErrorSite_ServeHttp(rec, req)
	// check that the link is correctly set and the right error is given --> link should be just the host
	assert.Equal(t, 200, rec.Result().StatusCode)
	assert.Equal(t, true, strings.Contains(rec.Body.String(), "Interner Server error"))
	assert.Equal(t, true, strings.Contains(rec.Body.String(), "<a class=\"btn btn-primary\" href=\"https://"+req.Host+"\">Zurück</a>"))
}
