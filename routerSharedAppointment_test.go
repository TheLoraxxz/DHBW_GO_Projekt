package main

import (
	"DHBW_GO_Projekt/authentifizierung"
	"DHBW_GO_Projekt/dateisystem"
	"DHBW_GO_Projekt/terminfindung"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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
	reader := strings.NewReader("newUsername=user&newPassword=user")
	req := httptest.NewRequest("GET", "localhost:443/shared?terminID="+terminId, reader)
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
	reader := strings.NewReader("newUsername=user&newPassword=user")
	req := httptest.NewRequest("GET", "localhost:443/shared?terminID="+terminId+"&selected="+termin.VorschlagTermine[0].ID, reader)
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
	reader := strings.NewReader("newUsername=user&newPassword=user")
	req := httptest.NewRequest("GET", "localhost:443/shared?terminID="+terminId, reader)
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

func TestCreateLinkServeHTTP_CreateTerminID(t *testing.T) {
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
	reader := strings.NewReader("newUsername=user&newPassword=user")
	req := httptest.NewRequest("GET", "localhost:443/shared?terminID="+terminId, reader)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	//execute link
	CreateLinkServeHTTP(rec, req)

}
