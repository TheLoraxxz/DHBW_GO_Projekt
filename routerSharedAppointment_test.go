package main

import (
	"DHBW_GO_Projekt/authentifizierung"
	"DHBW_GO_Projekt/dateisystem"
	"DHBW_GO_Projekt/terminfindung"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

var user = "test"

func TestMain(m *testing.M) {
	//create a user, termin and a shared termin
	authentifizierung.CreateUser(&user, &user)
	termin := dateisystem.CreateNewTermin("Test", "Test Description", dateisystem.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.UTC),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.UTC),
		true, "test")
	terminId, _ := terminfindung.CreateSharedTermin(&termin, &user)
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
	req := httptest.NewRequest("GET", "localhost:80/shared", reader)
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	//execute request
	AdminSiteServeHTTP(rec, req)
}
