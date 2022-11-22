package main

import (
	"DHBW_GO_Projekt/authentifizierung"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootHandler_ServeHTTP_reachable(t *testing.T) {
	// create a new request
	req := httptest.NewRequest("GET", "localhost:80", nil)
	//check
	recorder := httptest.NewRecorder()

	RootHandler{}.ServeHTTP(recorder, req)
	resp := recorder.Result()
	assert.Equal(t, 200, resp.StatusCode, "the Server doesn't return an answer")
	_, err := io.ReadAll(resp.Body)
	assert.Equal(t, err, nil, "No Error in the Body")
}

// TestRootHandler_ServeHTTP_GET
// tests that it creates the setup message right and the button is shown which is only on login
func TestRootHandler_ServeHTTP_GET(t *testing.T) {
	req := httptest.NewRequest("GET", "localhost:80", nil)
	rec := httptest.NewRecorder()
	//checks that
	RootHandler{}.ServeHTTP(rec, req)
	assert.Equal(t, strings.Contains(rec.Body.String(), "Einloggen"), true)
}

// TestRootHandler_ServeHTTP_POST_rightRequest
// checks if it submits a message that it automatically changes it to the cookie
func TestRootHandler_ServeHTTP_POST_rightRequest(t *testing.T) {
	//create admin to make sure the user is already created
	user := "admin"
	authentifizierung.CreateUser(&user, &user)
	// create request with the right body
	reader := strings.NewReader("user=admin&password=admin")
	req := httptest.NewRequest("POST", "localhost:80", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	// call function
	RootHandler{}.ServeHTTP(rec, req)
	//cookie should be set to the right response
	cookies := rec.Result().Cookies()
	assert.Equal(t, 1, len(cookies))
	assert.Equal(t, "SessionID-Kalender", cookies[0].Name)
	//uri shoud not be localhost:80/ but on a different --> kalenderansicht/... or anything else but /
	assert.Equal(t, http.StatusFound, rec.Code)
	url, err := rec.Result().Location()
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", url.Path)
}

// TestRootHandler_ServeHTTP_wrongRequest
// test that
func TestRootHandler_ServeHTTP_wrongRequest(t *testing.T) {
	//create admin to make sure the user is already created
	user := "admin"
	authentifizierung.CreateUser(&user, &user)
	//setup the caller
	reader := strings.NewReader("user=admin&password=user")
	req := httptest.NewRequest("POST", "localhost:80", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	// call function
	RootHandler{}.ServeHTTP(rec, req)
	//it should redirect to main website and it should not add any cookies
	assert.Equal(t, http.StatusContinue, rec.Code)
	assert.Equal(t, 0, len(rec.Result().Cookies()))
	// should redirect to same website
	url, _ := rec.Result().Location()
	assert.Equal(t, "/", url.Path)

}

// TestLogoutHandler_ServeHTTP_AllPossibilites
// checks that on logout the cookie is deleted
func TestLogoutHandler_ServeHTTP_RightInput(t *testing.T) {
	// ad user and get the cookie
	user := "admin"
	authentifizierung.CreateUser(&user, &user)
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	// setup request call for logout
	req := httptest.NewRequest("GET", "localhost:80/logout", nil)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	//execute handling
	LogoutHandler{}.ServeHTTP(rec, req)
	//check that the value is zero and not existing anymore
	assert.Equal(t, "", rec.Result().Cookies()[0].Value)
}
func TestLogoutHandler_ServeHTTP_WrongInput(t *testing.T) {
	req := httptest.NewRequest("GET", "localhost:80/logout", nil)
	rec := httptest.NewRecorder()
	LogoutHandler{}.ServeHTTP(rec, req)
	assert.Equal(t)

}
