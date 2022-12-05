/*
@author: 2447899 8689159 3000685
*/
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
	req := httptest.NewRequest("GET", "localhost:443", nil)
	//check
	recorder := httptest.NewRecorder()

	RootHandler{}.ServeHTTP(recorder, req)
	resp := recorder.Result()
	//should be reachable
	assert.Equal(t, 200, resp.StatusCode, "the Server doesn't return an answer")
	_, err := io.ReadAll(resp.Body)
	assert.Equal(t, err, nil, "No Error in the Body")
}

// TestRootHandler_ServeHTTP_GET
// tests that it creates the setup message right and the button is shown which is only on login
func TestRootHandler_ServeHTTP_GET(t *testing.T) {
	req := httptest.NewRequest("GET", "localhost:443", nil)
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
	req := httptest.NewRequest("POST", "localhost:443", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	// call function
	RootHandler{}.ServeHTTP(rec, req)
	//cookie should be set to the right response
	cookies := rec.Result().Cookies()
	assert.Equal(t, 1, len(cookies))
	assert.Equal(t, "SessionID-Kalender", cookies[0].Name)
	//uri shoud not be localhost:443/ but on a different --> kalenderansicht/... or anything else but /
	assert.Equal(t, http.StatusFound, rec.Code)
	url, err := rec.Result().Location()
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", url.Path)
}

// TestRootHandler_ServeHTTP_wrongRequest
// test what root does if it posts a wrong request
func TestRootHandler_ServeHTTP_wrongRequest(t *testing.T) {
	//create admin to make sure the user is already created
	user := "admin"
	authentifizierung.CreateUser(&user, &user)
	//setup the caller
	reader := strings.NewReader("user=admin&password=user")
	req := httptest.NewRequest("POST", "localhost:443", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	// call function
	RootHandler{}.ServeHTTP(rec, req)
	//it should redirect to main website and it should not add any cookies
	assert.Equal(t, http.StatusContinue, rec.Code)
	assert.Equal(t, 0, len(rec.Result().Cookies()))
	// should redirect to same website
	url, _ := rec.Result().Location()
	assert.Equal(t, "/error", url.Path)

}

// TestLogoutHandler_ServeHTTP_RightInput
// checks that on logout the cookie is deleted
func TestLogoutHandler_ServeHTTP_RightInput(t *testing.T) {
	// ad user and get the cookie
	user := "admin"
	authentifizierung.CreateUser(&user, &user)
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	// setup request call for logout
	req := httptest.NewRequest("GET", "localhost:443/logout", nil)
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

// TestLogoutHandler_ServeHTTP_WrongInput
// tests what happens to logout if it is the wrong input
func TestLogoutHandler_ServeHTTP_WrongInput(t *testing.T) {
	//setup request
	req := httptest.NewRequest("GET", "localhost:443/logout", nil)
	rec := httptest.NewRecorder()
	LogoutHandler{}.ServeHTTP(rec, req)
	//should return to the original path because an error occured
	url, _ := rec.Result().Location()
	assert.Equal(t, "", url.Path)
	assert.Equal(t, http.StatusContinue, rec.Code)

}

func TestCreatUserHandler_ServeHTTP_GETMetod(t *testing.T) {
	//get cookie and setup
	user := "admin"
	authentifizierung.CreateUser(&user, &user)
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	// setup request call for logout
	req := httptest.NewRequest("GET", "localhost:443/user/create", nil)
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
	//execute handler
	CreatUserHandler{}.ServeHTTP(rec, req)
	//it should contain the h3 header
	assert.Equal(t, strings.Contains(rec.Body.String(), "<h3>Neuen Nutzer erstellen</h3>"), true)
}

// TestCreatUserHandler_ServeHTTP_CorrectPost
// tests that it creates a user on creation
func TestCreatUserHandler_ServeHTTP_CorrectPost(t *testing.T) {
	//setup user and cookie
	user := "admin"
	authentifizierung.CreateUser(&user, &user)
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	//setup the caller
	reader := strings.NewReader("newUsername=user&newPassword=user")
	req := httptest.NewRequest("POST", "localhost:443", reader)
	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	//execute request
	CreatUserHandler{}.ServeHTTP(rec, req)
	//it should be redirect to user website
	url, _ := rec.Result().Location()
	assert.Equal(t, "", url.Path)
	assert.Equal(t, http.StatusContinue, rec.Code)
	//the user should exist and the authentication should return true
	user = "user"
	userExists, _ := authentifizierung.AuthenticateUser(&user, &user)
	assert.Equal(t, true, userExists)
}

// TestCreatUserHandler_ServeHTTP_NoCookie
// when the cookie is not set or wrong it should return to root
func TestCreatUserHandler_ServeHTTP_NoCookie(t *testing.T) {
	//setup the caller
	reader := strings.NewReader("s &=&?! user")
	req := httptest.NewRequest("POST", "localhost:443", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	//execute handler
	CreatUserHandler{}.ServeHTTP(rec, req)
	//it should redirect to the main path
	url, _ := rec.Result().Location()
	assert.Equal(t, "", url.Path)
	assert.Equal(t, http.StatusContinue, rec.Code)
}

// TestCreatUserHandler_ServeHTTP_WrongInput
// tests that it should return to null if the user is empty or wrong formatted in the post request
func TestCreatUserHandler_ServeHTTP_WrongInput(t *testing.T) {
	//create user according to everything
	user := "admin"
	authentifizierung.CreateUser(&user, &user)
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	//setup the caller
	reader := strings.NewReader("s &=&?! user")
	req := httptest.NewRequest("POST", "localhost:443", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)
	//execute handler
	CreatUserHandler{}.ServeHTTP(rec, req)
	//it should redirect to the main path
	url, _ := rec.Result().Location()
	assert.Equal(t, "", url.Path)
	assert.Equal(t, http.StatusContinue, rec.Code)
}

func TestCreateUserHandler_ServeHTTP_wrongCookie(t *testing.T) {
	//setup wrong cookie but with the right value
	user := "admin"
	authentifizierung.CreateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    "cookieValue|test",
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	//setup the caller
	req := httptest.NewRequest("GET", "localhost:443/user", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)
	//execute
	CreatUserHandler{}.ServeHTTP(rec, req)
	//should redirect to home website
	assert.Equal(t, http.StatusContinue, rec.Code)
	url, _ := rec.Result().Location()
	assert.Equal(t, "", url.Path)
}

func TestUserHandler_ServeHTTP_GETRequest(t *testing.T) {
	//create user according to everything
	user := "admin"
	authentifizierung.CreateUser(&user, &user)
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	//setup the caller
	req := httptest.NewRequest("GET", "localhost:443/user", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)
	//execute
	UserHandler{}.ServeHTTP(rec, req)
	//should show up the name of the user and the correct format
	assert.Equal(t, strings.Contains(rec.Body.String(), "Username: "+user), true)
	assert.Equal(t, strings.Contains(rec.Body.String(), "<h3>Konto</h3>"), true)
}

// TestUserHandler_ServeHTTP_NoCookie
// tests if it redirects on no cookie
func TestUserHandler_ServeHTTP_NoCookie(t *testing.T) {
	//setup call without cookie
	req := httptest.NewRequest("GET", "localhost:443/user", nil)
	rec := httptest.NewRecorder()
	UserHandler{}.ServeHTTP(rec, req)
	//shouild redirect to newstatus
	assert.Equal(t, http.StatusContinue, rec.Code)
	url, _ := rec.Result().Location()
	assert.Equal(t, "/error", url.Path)
}

// TestUserHandler_ServeHTTP_wrongCookie
// tests if it redirects if it not the right cookie
func TestUserHandler_ServeHTTP_wrongCookie(t *testing.T) {
	//setup wrong cookie but with the right value
	user := "admin"
	authentifizierung.CreateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    "cookieValue|test",
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	//setup the caller
	req := httptest.NewRequest("GET", "localhost:443/user", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)
	//execute
	UserHandler{}.ServeHTTP(rec, req)
	//should redirect to home website
	assert.Equal(t, http.StatusContinue, rec.Code)
	url, _ := rec.Result().Location()
	assert.Equal(t, "/error", url.Path)
}

// TestChangeUserHandler_ServeHTTP_NoCookie
// tests if it redirects on no cookie
func TestChangeUserHandler_ServeHTTP_NoCookie(t *testing.T) {
	//setup call without cookie
	req := httptest.NewRequest("GET", "localhost:443/user", nil)
	rec := httptest.NewRecorder()
	ChangeUserHandler{}.ServeHTTP(rec, req)
	//shouild redirect to newstatus
	assert.Equal(t, http.StatusContinue, rec.Code)
	url, _ := rec.Result().Location()
	assert.Equal(t, "", url.Path)
}

// TestChangeUserHandler_ServeHTTP_wrongCookie
// tests if it redirects if it not the right cookie
func TestChangeUserHandler_ServeHTTP_wrongCookie(t *testing.T) {
	//setup wrong cookie but with the right value
	user := "admin"
	authentifizierung.CreateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    "cookieValue|test",
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	//setup the caller
	req := httptest.NewRequest("GET", "localhost:443/user", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)
	//execute
	ChangeUserHandler{}.ServeHTTP(rec, req)
	//should redirect to home website
	assert.Equal(t, http.StatusContinue, rec.Code)
	url, _ := rec.Result().Location()
	assert.Equal(t, "", url.Path)
}

// TestChangeUserHandler_ServeHTTP_GET
// should return the right template with altes password and new password in it
func TestChangeUserHandler_ServeHTTP_GET(t *testing.T) {
	//create user according to everything
	user := "admin"
	authentifizierung.CreateUser(&user, &user)
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	//setup the caller
	req := httptest.NewRequest("GET", "localhost:443/user", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)
	//execute
	ChangeUserHandler{}.ServeHTTP(rec, req)
	//should create template with the right input of old and new password
	assert.Equal(t, strings.Contains(rec.Body.String(), "Neues Passwort:"), true)
	assert.Equal(t, strings.Contains(rec.Body.String(), "Altes Passwort:"), true)
}

func TestChangeUserHandler_ServeHTTP_CorrectPost(t *testing.T) {
	//setup user and cookie
	user := "admin"
	authentifizierung.CreateUser(&user, &user)
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	//setup the caller
	reader := strings.NewReader("oldPassword=admin&newPassword=user")
	req := httptest.NewRequest("POST", "localhost:443", reader)
	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	//execute request
	ChangeUserHandler{}.ServeHTTP(rec, req)
	//it should be redirect to user website
	url, _ := rec.Result().Location()
	assert.Equal(t, "/user", url.Path)
	assert.Equal(t, http.StatusContinue, rec.Code)
	//the user should  change with the right authentication
	password := "user"
	userExists, _ := authentifizierung.AuthenticateUser(&user, &password)
	assert.Equal(t, true, userExists)
	//cookie given back should pass the authentication
	isallowed, _ := authentifizierung.CheckCookie(&rec.Result().Cookies()[0].Value)
	assert.Equal(t, true, isallowed)
}

func TestChangeUserHandler_ServeHTTP_wronguser(t *testing.T) {
	user := "admin"
	authentifizierung.CreateUser(&user, &user)
	oldpassw := "user"
	authentifizierung.ChangeUser(&user, &oldpassw, &user)
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	//setup the caller and insert wrong password
	reader := strings.NewReader("oldPassword=wrongPassword&newPassword=user")
	req := httptest.NewRequest("POST", "localhost:443", reader)
	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	//execute request
	ChangeUserHandler{}.ServeHTTP(rec, req)
	//it should be redirect to user website
	url, _ := rec.Result().Location()
	assert.Equal(t, "/error", url.Path)
	assert.Equal(t, http.StatusContinue, rec.Code)
	//the user shoudn't have changed
	password := "user"
	userExists, _ := authentifizierung.AuthenticateUser(&user, &password)
	assert.Equal(t, false, userExists)

}
