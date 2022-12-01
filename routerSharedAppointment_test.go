package main

import (
	"DHBW_GO_Projekt/authentifizierung"
	"DHBW_GO_Projekt/dateisystem"
	"DHBW_GO_Projekt/terminfindung"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

var user = "test"
var terminId = ""

func TestMain(m *testing.M) {
	//create a user, termin and a shared termin
	setErrorconfigs()
	authentifizierung.CreateUser(&user, &user)
	termin := dateisystem.CreateNewTermin("Test", "Test Description", dateisystem.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.UTC),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.UTC),
		true, "test")
	terminId, _ = terminfindung.CreateSharedTermin(&termin, &user)
	startDate := time.Date(2022, 12, 10, 12, 0, 0, 0, time.UTC)
	endDate := time.Date(2022, 12, 10, 12, 0, 0, 0, time.UTC)
	terminfindung.CreateNewProposedDate(startDate, endDate, &user, &terminId, false)
	name := "test2"
	terminfindung.CreatePerson(&name, &terminId, &user)
	//run test
	code := m.Run()
	//delete data
	dateisystem.DeleteAll(dateisystem.GetTermine(user), user)
	os.Exit(code)

}

func TestAdminSiteServeHTTP_normalCookie(t *testing.T) {
	//setup the caller
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	req := httptest.NewRequest("GET", "localhost:443/shared?terminID="+terminId, nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	//execute request
	AdminSiteServeHTTP(rec, req)
	//should accept it
	assert.Equal(t, 200, rec.Result().StatusCode)
	termin, _ := terminfindung.GetTerminFromShared(&user, &terminId)
	//check that normal works
	assert.Equal(t, true, strings.Contains(rec.Body.String(), termin.Info.Title))

}
func TestAdminSiteServeHTTP_selectDate(t *testing.T) {
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	termin, _ := terminfindung.GetTerminFromShared(&user, &terminId)
	startDate := time.Date(2022, 12, 12, 12, 12, 12, 12, time.UTC)
	enddate := time.Date(2022, 12, 13, 12, 12, 12, 12, time.UTC)
	terminfindung.CreateNewProposedDate(startDate, enddate, &user, &terminId, false)
	req := httptest.NewRequest("GET", "localhost:443/shared?terminID="+terminId+"&selected="+termin.VorschlagTermine[0].ID, nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	//execute request
	AdminSiteServeHTTP(rec, req)
	termin, _ = terminfindung.GetTerminFromShared(&user, &terminId)
	//check that the feinaltermin has been really checked
	assert.NotEmpty(t, termin.FinalTermin)
	assert.Equal(t, termin.VorschlagTermine[0].ID, termin.FinalTermin.ID)

}

func TestAdminSiteServeHTTP_wrongCookie(t *testing.T) {
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    "sad",
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	req := httptest.NewRequest("GET", "localhost:443/shared?terminID="+terminId, nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	//execute request
	AdminSiteServeHTTP(rec, req)
	assert.Equal(t, 100, rec.Result().StatusCode)
	urls, err := rec.Result().Location()
	assert.Equal(t, err, nil)
	termin, _ := terminfindung.GetTerminFromShared(&user, &terminId)
	assert.Equal(t, "wrongAuthentication", urls.Query().Get("type"))
	assert.Equal(t, false, strings.Contains(termin.Info.Title, rec.Body.String()))
}

func TestCreateLinkServeHTTP_GETRequest(t *testing.T) {
	//setup
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	req := httptest.NewRequest("GET", "localhost:443/shared/create/link?terminID="+terminId, nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	//execute link
	CreateLinkServeHTTP(rec, req)
	//should return statuscode 200
	assert.Equal(t, 200, rec.Result().StatusCode)
	assert.Equal(t, true, strings.Contains(rec.Body.String(), "<h2>Neue Person einladen</h2>"))

}

func TestCreateLinkServeHTTP_noCoockie(t *testing.T) {
	req := httptest.NewRequest("GET", "localhost:443/shared/create/link?terminID="+terminId, nil)
	rec := httptest.NewRecorder()
	CreateLinkServeHTTP(rec, req)

	assert.Equal(t, 100, rec.Result().StatusCode)
	assert.Equal(t, false, strings.Contains(rec.Body.String(), "<h2>Neue Person einladen</h2>"))
}

func TestCreateLinkServeHTTP_RightPostRequest(t *testing.T) {
	//setup
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	reader := strings.NewReader("name=testuser")

	req := httptest.NewRequest("POST", "localhost:443/shared/create/link?terminID="+terminId, reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	//execute link
	CreateLinkServeHTTP(rec, req)
	//tests that it should give back the right link and it should be code 200
	assert.Equal(t, 200, rec.Result().StatusCode)
	termin, _ := terminfindung.GetTerminFromShared(&user, &terminId)
	assert.Equal(t, true, strings.Contains(rec.Body.String(), "https://example.com/shared/public?apiKey="+termin.Persons["testuser"].Url))
}

func TestCreateLinkServeHTTP_WrongPostRequest(t *testing.T) {
	//setup
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	reader := strings.NewReader("/&?sas=?")

	req := httptest.NewRequest("POST", "localhost:443/shared/create/link?terminID="+terminId, reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	//execute link
	CreateLinkServeHTTP(rec, req)
	//tests that it should give back the right link and it should be code 200
	assert.Equal(t, 100, rec.Result().StatusCode)
}

func TestServeHTTPSharedAppCreateDate_GetRequestNormal(t *testing.T) {
	// setup statement
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	req := httptest.NewRequest("GET", "localhost:443/shared/create/app?terminID="+terminId, nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	//execute statement
	ServeHTTPSharedAppCreateDate(rec, req)
	//should have the form in the response
	assert.Equal(t, 200, rec.Result().StatusCode)
	assert.Equal(t, true, strings.Contains(rec.Body.String(), "<input required id=\"enddate\" type=\"date\" class=\"form-control\" name=\"enddate\">"))
}

func TestServeHTTPSharedAppCreateDate_WrongCookie(t *testing.T) {
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    "cookieValue",
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	req := httptest.NewRequest("GET", "localhost:443/shared/create/app?terminID="+terminId, nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	//execute statement
	ServeHTTPSharedAppCreateDate(rec, req)
	//should be redirected
	assert.Equal(t, 100, rec.Result().StatusCode)
	assert.Equal(t, false, strings.Contains(rec.Body.String(), "<input required id=\"enddate\" type=\"date\" class=\"form-control\" name=\"enddate\">"))
}

func TestServeHTTPSharedAppCreateDate_FalseTerminId(t *testing.T) {
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	req := httptest.NewRequest("GET", "localhost:443/shared/create/app?terminID=asdasd", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	//execute statement
	ServeHTTPSharedAppCreateDate(rec, req)
	//should be redirected
	assert.Equal(t, 100, rec.Result().StatusCode)
	assert.Equal(t, false, strings.Contains(rec.Body.String(), "<input required id=\"enddate\" type=\"date\" class=\"form-control\" name=\"enddate\">"))
	//should redirect to right url should give authentication error
	urls, err := rec.Result().Location()
	assert.Equal(t, nil, err)
	assert.Equal(t, "https://"+req.Host+"/error?type=wrongAuthentication&link="+url.QueryEscape("/shared?terminID=asdasd"), urls.String())
}

func TestServeHTTPSharedAppCreateDate_POSTRequest(t *testing.T) {
	//setup (right cookie in post)
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	reader := strings.NewReader("startdate=2022-12-12&enddate=2022-12-13")
	req := httptest.NewRequest("POST", "localhost:443/shared/create/app?terminID="+terminId, reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	//execute
	termin, _ := terminfindung.GetTerminFromShared(&user, &terminId)
	ServeHTTPSharedAppCreateDate(rec, req)
	urls, err := rec.Result().Location()
	assert.Equal(t, nil, err)
	assert.Equal(t, 100, rec.Result().StatusCode)
	termin2, _ := terminfindung.GetTerminFromShared(&user, &terminId)
	//should have created one proposed termin
	assert.Equal(t, 1, len(termin2.VorschlagTermine)-len(termin.VorschlagTermine))
	assert.Equal(t, "https://"+req.Host+"/shared?terminID="+terminId, urls.String())
}

func TestServeHTTPSharedAppCreateDate_WrongDataFormat(t *testing.T) {
	//setup (right cookie in post)
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	reader := strings.NewReader("startdate=23-332-32&enddate=23-34-66")
	req := httptest.NewRequest("POST", "localhost:443/shared/create/app?terminID="+terminId, reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	//execute
	ServeHTTPSharedAppCreateDate(rec, req)
	//should reject it
	assert.Equal(t, 100, rec.Result().StatusCode)
	urls, err := rec.Result().Location()
	assert.Equal(t, nil, err)
	//the new link should be error and it should give back the continue site
	assert.Equal(t, true, strings.Contains(rec.Body.String(), "\">Continue</a>."))
	assert.Equal(t, "https://"+req.Host+"/error?type=wrong_date_format&link="+url.QueryEscape("/shared?terminID="+terminId), urls.String())
}

func TestServeHTTPSharedAppCreateDate_StartDateAfterEnddate(t *testing.T) {
	//setup (right cookie in post)
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	reader := strings.NewReader("startdate=2022-12-14&enddate=2022-12-12")
	req := httptest.NewRequest("POST", "localhost:443/shared/create/app?terminID="+terminId, reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	//execute
	ServeHTTPSharedAppCreateDate(rec, req)
	//should reject it
	assert.Equal(t, 100, rec.Result().StatusCode)
	urls, err := rec.Result().Location()
	assert.Equal(t, nil, err)
	//the new link should be error and it should give back the continue site
	assert.Equal(t, true, strings.Contains(rec.Body.String(), "\">Continue</a>."))
	assert.Equal(t, "https://"+req.Host+"/error?type=dateIsAfter&link="+url.QueryEscape("/shared/create/app?terminID="+terminId), urls.String())
}

func TestShowAllLinksServeHttp_GETRequest(t *testing.T) {
	//create some persons
	name := "test"
	terminfindung.CreatePerson(&name, &terminId, &user)
	name = "user"
	terminfindung.CreatePerson(&name, &terminId, &user)
	//setup (right cookie in post)
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	req := httptest.NewRequest("GET", "localhost:443/shared/showAllLink?terminID="+terminId, nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	//execute statement
	ShowAllLinksServeHttp(rec, req)

	assert.Equal(t, 200, rec.Result().StatusCode)
	assert.Equal(t, true, strings.Contains(rec.Body.String(), "</a>"))
}

func TestShowAllLinksServeHttp_no_cookie(t *testing.T) {
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    "cookieValue",
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	req := httptest.NewRequest("GET", "localhost:443/shared/showAllLink?terminID="+terminId, nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	//execute statement
	ShowAllLinksServeHttp(rec, req)
	// should automatically return to base website and error
	assert.Equal(t, 100, rec.Result().StatusCode)
	urls, err := rec.Result().Location()

	assert.Equal(t, nil, err)
	assert.Equal(t, "https://"+req.Host+"/error?type=wrongAuthentication&link="+url.QueryEscape("/"), urls.String())
}

func TestShowAllLinksServeHttp_TerminIDNotSet(t *testing.T) {
	//setup (right cookie in post)
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	cookie := &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	req := httptest.NewRequest("GET", "localhost:443/shared/showAllLink?terminID=asd", nil)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	//execute statement
	ShowAllLinksServeHttp(rec, req)
	// check that if not set it says authentication wrong
	assert.Equal(t, 100, rec.Result().StatusCode)
	urls, err := rec.Result().Location()
	assert.Equal(t, nil, err)
	assert.Equal(t, "https://"+req.Host+"/error?type=wrongAuthentication&link="+url.QueryEscape("/"), urls.String())
}

func TestPublicSharedWebsite(t *testing.T) {

}
