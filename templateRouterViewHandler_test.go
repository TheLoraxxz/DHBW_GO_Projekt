package main

import (
	"DHBW_GO_Projekt/authentifizierung"
	ds "DHBW_GO_Projekt/dateisystem"
	"DHBW_GO_Projekt/terminfindung"
	"github.com/stretchr/testify/assert"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

/*
**************************************************************************************************************
Setup-Funktionen, die zu Begin von vielen Tests benötigt werden, um diese Tests aufzusetzen (=> Setups für Tests).
Teardown-Funktion, um nach einem Test erstellte Testtermine wieder zu löschen,
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
var cookie *http.Cookie
var vmHandler ViewManagerHandler

func setupCookie(correctCookie bool) {
	var value string
	user := "admin"
	authentifizierung.CreateUser(&user, &user)
	_, cookieValue := authentifizierung.AuthenticateUser(&user, &user)
	if correctCookie {
		value = cookieValue
	} else {
		value = "cookieValue|test"
	}
	cookie = &http.Cookie{
		Name:     "SessionID-Kalender",
		Value:    value,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}

func setupHandler() {
	vmHandler = ViewManagerHandler{}
	path, _ := os.Getwd()
	vmHandler.viewManagerTpl = template.Must(template.New("tbl.html").ParseFiles(path+"/assets/sites/tbl.html", path+"/assets/templates/header.html", path+"/assets/templates/footer.html", path+"/assets/templates/creator.html", path+"/assets/templates/listing.html"))
	template.Must(vmHandler.viewManagerTpl.New("liste.html").ParseFiles(path+"/assets/sites/liste.html", path+"/assets/templates/header.html", path+"/assets/templates/footer.html", path+"/assets/templates/creator.html"))
	template.Must(vmHandler.viewManagerTpl.New("editor.html").ParseFiles(path+"/assets/sites/editor.html", path+"/assets/templates/header.html", path+"/assets/templates/footer.html", path+"/assets/templates/listing.html"))
	template.Must(vmHandler.viewManagerTpl.New("filterTermins.html").ParseFiles(path + "/assets/sites/filterTermins.html"))
}

// teardownDeleteTermin
// zur Probe erstellten Termine löschen
func teardownDeleteTermin() {
	ds.DeleteAll(vmHandler.vm.TerminCache, vmHandler.vm.Username)
}

/*
**************************************************************************************************************
Hilfsfunktionen, die in den Tests öfters benötigt werden
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
// createTerminRequest
// Parameter: Die Werte für die Request
// Rückgabewert: Request, um einen Termin zu erstellen.
// createReaderData
// Parameter: Die Werte für die Request
// Rückgabewert: *Reader, um eine Request zu erstellen.
func createReaderData(shared, title, description, repeat, date, endDate, mode, id string) string {
	data := url.Values{}
	data.Add("ID", id)
	data.Add("shared", shared)
	data.Add("editing", mode)
	data.Add("title", title)
	data.Add("description", description)
	data.Add("rep", repeat)
	data.Add("date", date)
	data.Add("endDate", endDate)

	//Erstellen der Request
	return data.Encode()
}

// initTestTermin
// sorgt dafür, dass ein Testtermin erstellt wird, der im Dateisystem liegt.
func initTestTermin(shared bool) ds.Termin {
	termin := ds.CreateNewTermin("Spezifischer-Test-Termin", "Test Description", ds.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.UTC),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.UTC),
		shared, "admin")
	if shared {
		user = "admin"
		terminfindung.CreateSharedTermin(&termin, &user)
	}
	return termin
}

/*
**************************************************************************************************************
Tests für ViewManagerHandler
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
//testViewManager_WrongCookie
//wenn die Authentifizierung fehlschlägt, soll ein Redirect erfolgen
func testViewManager_WrongCookie(t *testing.T) {
	//Test-Setups
	setupCookie(false)
	setupHandler()

	//setup the caller
	req := httptest.NewRequest("GET", "localhost:443/user/view/table", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//execute
	handler := ViewManagerHandler{}
	handler.ServeHTTP(rec, req)

	//Should redirect
	assert.Equal(t, http.StatusContinue, rec.Code, "Da der die Authentifizierung fehlgeschlagen ist, sollte der Nutzer umgeleitet werden.")
}

// testViewManager_NewViewManagerObject
// falls das View-manger-Objekt im ViewManagerHandler-Struct noch nil ist, soll eines erstellt werden
func testViewManager_NewViewManagerObject(t *testing.T) {
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	req := httptest.NewRequest("GET", "localhost:443/user/view/table", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	assert.True(t, vmHandler.vm == nil, "Das ViewManager-Objekt sollte zu Begin nil sein.")

	//Execute
	vmHandler.ServeHTTP(rec, req)

	assert.NotEqual(t, nil, vmHandler.vm, "Das ViewManager-Objekt sollte zu Begin nicht mehr nil sein.")
}

// testViewManager_EditTerminGetRequest
// falls eine Request zum Bearbeiten eines Termins kommt, soll dass Bearbeitungsformular angezeigt werden
func testViewManager_EditTerminGetRequest(t *testing.T) {
	// testtermin zum Bearbeiten ins Dateisystem legen
	testTermin := initTestTermin(false)

	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	reader := createReaderData("", "", "", "", "", "", "", testTermin.ID)
	req := httptest.NewRequest("GET", "localhost:443/user/view/table?edit=true&"+reader, nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	assert.True(t, strings.Contains(rec.Body.String(), "Termin bearbeiten"), "Es sollte ein Bearbeitungsmodus-Feld geben. ")

	//delete Test Termin
	teardownDeleteTermin()
}

// testViewManager_EditTerminGetRequestError
// ERROR bei Request zum Bearbeiten eines Termins kommt, soll dass Bearbeitungsformular angezeigt werden
func testViewManager_EditTerminGetRequestError(t *testing.T) {
	// testtermin zum Bearbeiten ins Dateisystem legen
	initTestTermin(false)

	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	reader := createReaderData("", "", "", "", "", "", "", "banananananan")
	req := httptest.NewRequest("GET", "localhost:443/user/view/table?edit=true&"+reader, nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)

	//Should redirect
	assert.Equal(t, http.StatusContinue, rec.Code, "Da der die Authentifizierung fehlgeschlagen ist, sollte der Nutzer umgeleitet werden.")
	urls, err := rec.Result().Location()
	assert.Equal(t, nil, err)
	assert.Equal(t, "https://"+req.Host+"/error?type=shared_wrong_terminId&link="+url.QueryEscape("/user/view/table"), urls.String())
	//delete Test Termin
	teardownDeleteTermin()
}

// testViewManager_EditTerminGetRequest
// falls eine Post-Request zum Bearbeiten eines Termins kommt, soll der bearbeitete Termin in der Ansicht angezeigt werden
func testViewManager_EditTerminPostRequest(t *testing.T) {
	// testtermin zum Bearbeiten ins Dateisystem legen
	testTermin := initTestTermin(false)

	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	reader := strings.NewReader(createReaderData("", "Bearbeiteter-Titel", "beschreibung", "täglich", "2022-01-12", "2023-01-30", "editing", testTermin.ID))
	req := httptest.NewRequest("POST", "localhost:443/user/view/list?edit=true&", reader)
	rec := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	assert.True(t, strings.Contains(rec.Body.String(), "Bearbeiteter-Titel"), "Es sollte ein neuen Termin mit dem Titel: Bearbeiteter-Titel-Feld geben. ")

	//delete Test Termin
	teardownDeleteTermin()
}

// testViewManager_EditTerminGetRequestError
// ERROR im Formular zum Bearbeiten eines Termins => redirect  (Error ist hier ein fehlerhaftes Datum, bei dem Tag und mOnat vertauscht ist)
func testViewManager_EditTerminPostRequestError(t *testing.T) {
	// testtermin zum Bearbeiten ins Dateisystem legen
	testTermin := initTestTermin(false)

	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	reader := strings.NewReader(createReaderData("", "Bearbeiteter-Titel", "beschreibung", "täglich", "2022-01-12", "2023-30-01", "editing", testTermin.ID))
	req := httptest.NewRequest("POST", "localhost:443/user/view/list?edit=true&", reader)
	rec := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)

	//Should redirect
	assert.Equal(t, http.StatusContinue, rec.Code, "Da der Termin 30.01.2023 nicht existiert, kommt es zu einem Error: Nutzer muss umgeleitet werden.")
	urls, err := rec.Result().Location()
	assert.Equal(t, nil, err)
	assert.Equal(t, "https://"+req.Host+"/error?type=wrong_date_format&link="+url.QueryEscape("/user/view/list"), urls.String(), "Hierhin sollte umgeleitet werden.")

	//delete Test Termin
	teardownDeleteTermin()
}

// testViewManager_DeleteSharedTermin
func testViewManager_DeleteSharedTermin(t *testing.T) {
	// testtermin zum Bearbeiten ins Dateisystem legen
	testTermin := initTestTermin(true)

	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller for looking if shared termin exits
	req := httptest.NewRequest("GET", "localhost:443/user/view/list", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)
	vmHandler.ServeHTTP(rec, req)

	//Vor der Ausführung sollte es einen geteilten Termin geben
	assert.True(t, strings.Contains(rec.Body.String(), testTermin.Title), "Es sollte ein neuen Termin mit dem Titel: Test-Termin geben. ")
	assert.True(t, strings.Contains(rec.Body.String(), "Geteilter-Terminvorschlag"), "Es sollte ein Feld mit der Info: Geteilter-Terminvorschlag geben. ")

	//Setup the caller for deleting
	reader := "deleteSharedTermin=" + testTermin.ID
	req = httptest.NewRequest("GET", "localhost:443/user/view/list?"+reader, nil)
	rec = httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)

	assert.False(t, strings.Contains(rec.Body.String(), testTermin.Title), "Es sollte keinen Termin mit dem Titel: Test-Termin mehr geben. ")

	//delete Test Termin
	teardownDeleteTermin()
}

// testViewManager_DeleteSharedTerminError
func testViewManager_DeleteSharedTerminError(t *testing.T) {
	// testtermin zum Bearbeiten ins Dateisystem legen
	testTermin := initTestTermin(true)

	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller for looking if shared termin exits
	req := httptest.NewRequest("GET", "localhost:443/user/view/list", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)
	vmHandler.ServeHTTP(rec, req)

	//Vor der Ausführung sollte es einen geteilten Termin geben
	assert.True(t, strings.Contains(rec.Body.String(), testTermin.Title), "Es sollte ein neuen Termin mit dem Titel: Test-Termin geben. ")
	assert.True(t, strings.Contains(rec.Body.String(), "Geteilter-Terminvorschlag"), "Es sollte ein Feld mit der Info: Geteilter-Terminvorschlag geben. ")

	//Setup the caller for deleting
	reader := "deleteSharedTermin=" + "banana"
	req = httptest.NewRequest("GET", "localhost:443/user/view/list?"+reader, nil)
	rec = httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)

	//Should redirect
	assert.Equal(t, http.StatusContinue, rec.Code, "Da die ID nicht existiert, kommt es zu einem Error: Nutzer muss umgeleitet werden.")
	urls, err := rec.Result().Location()
	assert.Equal(t, nil, err)
	assert.Equal(t, "https://"+req.Host+"/error?type=shared_wrong_terminId&link="+url.QueryEscape("/user/view/list"), urls.String(), "Hierhin sollte umgeleitet werden.")
	//delete Test Termin
	teardownDeleteTermin()
}

// testViewManager_NewCreateTerminRequest
// falls eine Request zum erstellen eines Termins kommt, soll dieser in der Ansichten angezeigt werden
func testViewManager_NewCreateTerminRequest(t *testing.T) {
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	reader := strings.NewReader(createReaderData("false", "Termin-Titel", "Beschreibung-termin", "niemals", "2022-12-02", "2022-12-02", "", ""))
	req := httptest.NewRequest("POST", "localhost:443/user/view/table?create=true", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	assert.True(t, strings.Contains(rec.Body.String(), "Termine: 1"), "Es sollte ein Tabellen-Feld geben, indem ein Termin eingetragen ist. ")

	//delete Test Termin
	teardownDeleteTermin()
}

// testViewManager_NewCreateTerminRequest
// falls eine Request zum erstellen eines Termins kommt, soll dieser in der Ansichten angezeigt werden
func testViewManager_NewCreateTerminRequestError(t *testing.T) {
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	reader := strings.NewReader(createReaderData("false", "Termin-Titel", "Beschreibung-termin", "niemals", "2022-12-02", "2022-12-02", "", ""))
	req := httptest.NewRequest("POST", "localhost:443/user/view/table?create=true", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	assert.True(t, strings.Contains(rec.Body.String(), "Termine: 1"), "Es sollte ein Tabellen-Feld geben, indem ein Termin eingetragen ist. ")

	//delete Test Termin
	teardownDeleteTermin()
}

/*
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
Ab hier folgen Tests, die in an den Table-Handler weitergeleitet werden
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/

func testViewManager_TvJumpMonthBack(t *testing.T) {
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	req := httptest.NewRequest("GET", "localhost:443/user/view/table?suche=minusMonat", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	newDate := time.Now().AddDate(0, -1, 0)
	newMonth := newDate.Month().String()
	newYear := strconv.Itoa(newDate.Year())
	assert.True(t, strings.Contains(rec.Body.String(), newMonth+" "+newYear), "Es sollte der vorherige Monat angezeigt werden (von heute aus gesehen). ")
}
func testViewManager_TvJumpMonthFor(t *testing.T) {
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	req := httptest.NewRequest("GET", "localhost:443/user/view/table?suche=plusMonat", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	newDate := time.Now().AddDate(0, +1, 0)
	newMonth := newDate.Month().String()
	newYear := strconv.Itoa(newDate.Year())
	assert.True(t, strings.Contains(rec.Body.String(), newMonth+" "+newYear), "Es sollte der nächste Monat angezeigt werden (von heute aus gesehen). ")
}
func testViewManager_TvSelectMonth(t *testing.T) {
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	req := httptest.NewRequest("GET", "localhost:443/user/view/table?monat=01", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	assert.True(t, strings.Contains(rec.Body.String(), "January"), "Es sollte der Monat Januar angezeigt werden (von heute aus gesehen). ")
}
func testViewManager_TvSelectMonthError(t *testing.T) {
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	req := httptest.NewRequest("GET", "localhost:443/user/view/table?monat=20", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)

	//Should redirect
	assert.Equal(t, http.StatusContinue, rec.Code, "Error, da es keinen 20. Monat gibt => umleiten ")
	urls, err := rec.Result().Location()
	assert.Equal(t, nil, err)
	assert.Equal(t, "https://"+req.Host+"/error?type=NowValidMonth&link="+url.QueryEscape("/user/view/table"), urls.String(), "Hierhin sollte umgeleitet werden.")
	//delete Test Termin
	teardownDeleteTermin()
}
func testViewManager_TvJumpYearBack(t *testing.T) {
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	req := httptest.NewRequest("GET", "localhost:443/user/view/table?jahr=Zurueck", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	newDate := time.Now().AddDate(-1, 0, 0)
	newYear := strconv.Itoa(newDate.Year())
	assert.True(t, strings.Contains(rec.Body.String(), newYear), "Es sollte das vorherige Jahr angezeigt werden (von heute aus gesehen). ")
}
func testViewManager_TvJumpYearFor(t *testing.T) {
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	req := httptest.NewRequest("GET", "localhost:443/user/view/table?jahr=Vor", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	newDate := time.Now().AddDate(1, 0, 0)
	newYear := strconv.Itoa(newDate.Year())
	assert.True(t, strings.Contains(rec.Body.String(), newYear), "Es sollte das nächste Jahr angezeigt werden (von heute aus gesehen). ")
}
func testViewManager_TvJumpToToday(t *testing.T) {
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller, um zunächst zu einem anderen Datum zu springen
	req := httptest.NewRequest("GET", "localhost:443/user/view/table?suche=plusMonat", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	//Kontrolle, ob Terminanzeige wirklich nicht mehr dem heutigen Datum entspricht
	newDate := time.Now().AddDate(0, +1, 0)
	newMonth := newDate.Month().String()
	newYear := strconv.Itoa(newDate.Year())
	assert.True(t, strings.Contains(rec.Body.String(), newMonth+" "+newYear), "Es sollte der nächste Monat angezeigt werden (von heute aus gesehen). ")

	//Setup the caller, um zum heutigen Datum zu springen
	req = httptest.NewRequest("GET", "localhost:443/user/view/table?datum=heute", nil)
	rec = httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)

	newDate = time.Now()
	newMonth = newDate.Month().String()
	newYear = strconv.Itoa(newDate.Year())
	assert.True(t, strings.Contains(rec.Body.String(), newMonth+" "+newYear), "Es sollte heutiger Monat & Jahr angezeigt werden (von heute aus gesehen). ")
}
func testViewManager_TvReload(t *testing.T) {
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller um ein anderes Datum zu setzen
	//Test-Setups
	req := httptest.NewRequest("GET", "localhost:443/user/view/table?jahr=Zurueck", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute: überprüfen, ob Datum wirklich nicht heute
	vmHandler.ServeHTTP(rec, req)
	newDate := time.Now().AddDate(-1, 0, 0)
	newYear := strconv.Itoa(newDate.Year())
	assert.True(t, strings.Contains(rec.Body.String(), newYear), "Es sollte das vorherige Jahr angezeigt werden (von heute aus gesehen). ")

	//Setup the caller, wenn Seite neu geladen wird (ohne Parameter)
	req = httptest.NewRequest("GET", "localhost:443/user/view/table", nil)
	rec = httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute um zu kontrollieren, ob Datum gesetzt wurde
	vmHandler.ServeHTTP(rec, req)
	newDate = time.Now()
	newMonth := newDate.Month().String()
	newYear = strconv.Itoa(newDate.Year())
	assert.True(t, strings.Contains(rec.Body.String(), newMonth+" "+newYear), "Es sollte das heutige Datum angezeigt werden. ")
}

/*
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
Ab hier folgen Tests, die in an den List-Handler weitergeleitet werden
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/
func testViewManager_LvSelectDate(t *testing.T) {
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	reader := "2023-01-30"
	req := httptest.NewRequest("GET", "localhost:443/user/view/list?selDate="+reader, nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	newDate := "2023-01-30"
	assert.True(t, strings.Contains(rec.Body.String(), newDate), "Es sollte das gewählte Datum angezeigt werden (von heute aus gesehen). ")
}
func testViewManager_LvSelectDateError(t *testing.T) {
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	reader := "2023-30-01"
	req := httptest.NewRequest("GET", "localhost:443/user/view/list?selDate="+reader, nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)

	//Should redirect
	assert.Equal(t, http.StatusContinue, rec.Code, "Error, da das Datum falsch ist => umleiten ")
	urls, err := rec.Result().Location()
	assert.Equal(t, nil, err)
	assert.Equal(t, "https://"+req.Host+"/error?type=wrong_date_format&link="+url.QueryEscape("/user/view/list"), urls.String(), "Hierhin sollte umgeleitet werden.")
}

func testViewManager_LvSelectEntriesPerPage(t *testing.T) {
	//Create 15 Termins
	for i := 0; i < 15; i++ {
		initTestTermin(false)
	}
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	req := httptest.NewRequest("GET", "localhost:443/user/view/list?Eintraege=10", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	assert.True(t, strings.Contains(rec.Body.String(), "1 von 2"), "Es sollte nun zwei Seiten zur Darstellung benötigt werden. ")

	//delete Test Termins
	teardownDeleteTermin()
}

func testViewManager_LvSelectEntriesPerPageError(t *testing.T) {
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	req := httptest.NewRequest("GET", "localhost:443/user/view/list?Eintraege=100", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)

	//Should redirect
	assert.Equal(t, http.StatusContinue, rec.Code, "Error, da ungültige Anzahl an Einträgen => umleiten ")
	urls, err := rec.Result().Location()
	assert.Equal(t, nil, err)
	assert.Equal(t, "https://"+req.Host+"/error?type=Unvalid_Entries_Per_Page&link="+url.QueryEscape("/user/view/list"), urls.String(), "Hierhin sollte umgeleitet werden.")
}

func testViewManager_LvJumpPageForward(t *testing.T) {
	//Create 15 Termins
	for i := 0; i < 15; i++ {
		initTestTermin(false)
	}
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	req := httptest.NewRequest("GET", "localhost:443/user/view/list?Seite=Vor", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	assert.True(t, strings.Contains(rec.Body.String(), "2 von 3"), "Man befindet sich auf Seite 2. Es sollte drei Seiten benötigt werden. ")

	//delete Test Termins
	teardownDeleteTermin()
}

func testViewManager_LvJumpPageBack(t *testing.T) {
	//Create 15 Termins
	for i := 0; i < 15; i++ {
		initTestTermin(false)
	}
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller for jump forward
	req := httptest.NewRequest("GET", "localhost:443/user/view/list?Seite=Vor", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	assert.True(t, strings.Contains(rec.Body.String(), "2 von 3"), "Man befindet sich auf Seite 2. Es sollte drei Seiten benötigt werden. ")

	//Setup the caller for jump back
	req = httptest.NewRequest("GET", "localhost:443/user/view/list?Seite=Zurueck", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec = httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	assert.True(t, strings.Contains(rec.Body.String(), "1 von 3"), "Man befindet sich auf Seite 1. Es sollte drei Seiten benötigt werden. ")

	//delete Test Termins
	teardownDeleteTermin()
}
func testViewManager_LvReload(t *testing.T) {
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller um ein anderes Datum zu setzen
	reader := "2023-01-30"
	req := httptest.NewRequest("GET", "localhost:443/user/view/list?selDate="+reader, nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute um zu kontrollieren, ob Datum gesetzt wurde
	vmHandler.ServeHTTP(rec, req)
	newDate := "2023-01-30"
	assert.True(t, strings.Contains(rec.Body.String(), newDate), "Es sollte das gewählte Datum angezeigt werden (von heute aus gesehen). ")

	//Setup the caller, wenn Seite neu geladen wird (ohne Parameter)
	req = httptest.NewRequest("GET", "localhost:443/user/view/list", nil)
	rec = httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute um zu kontrollieren, ob Datum gesetzt wurde
	vmHandler.ServeHTTP(rec, req)
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	assert.True(t, strings.Contains(rec.Body.String(), today.Format("2006-01-02")), "Es sollte das heutige Datum angezeigt werden. ")
}

/*
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
Ab hier folgen Tests, die in an den FilterView-Handler weitergeleitet werden
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
*/

func testViewManager_FvSelectEntriesPerPage(t *testing.T) {
	//Create 15 Termins
	for i := 0; i < 15; i++ {
		initTestTermin(false)
	}
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	req := httptest.NewRequest("GET", "localhost:443/user/view/filterTermins?Eintraege=10", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	assert.True(t, strings.Contains(rec.Body.String(), "1 von 2"), "Es sollte nun zwei Seiten zur Darstellung benötigt werden. ")

	//delete Test Termins
	teardownDeleteTermin()
}

func testViewManager_FvSelectEntriesPerPageError(t *testing.T) {
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	req := httptest.NewRequest("GET", "localhost:443/user/view/filterTermins?Eintraege=100", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)

	//Should redirect
	assert.Equal(t, http.StatusContinue, rec.Code, "Error, da ungültige Anzahl an Einträgen => umleiten ")
	urls, err := rec.Result().Location()
	assert.Equal(t, nil, err)
	assert.Equal(t, "https://"+req.Host+"/error?type=Unvalid_Entries_Per_Page&link="+url.QueryEscape("/user/view/filterTermins"), urls.String(), "Hierhin sollte umgeleitet werden.")
}

func testViewManager_FvJumpPageForward(t *testing.T) {
	//Create 15 Termins
	for i := 0; i < 15; i++ {
		initTestTermin(false)
	}
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller
	req := httptest.NewRequest("GET", "localhost:443/user/view/filterTermins?Seite=Vor", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	assert.True(t, strings.Contains(rec.Body.String(), "2 von 3"), "Man befindet sich auf Seite 2. Es sollte drei Seiten benötigt werden. ")

	//delete Test Termins
	teardownDeleteTermin()
}

func testViewManager_FvJumpPageBack(t *testing.T) {
	//Create 15 Termins
	for i := 0; i < 15; i++ {
		initTestTermin(false)
	}
	//Test-Setups
	setupCookie(true)
	setupHandler()

	//Setup the caller for jump forward
	req := httptest.NewRequest("GET", "localhost:443/user/view/filterTermins?Seite=Vor", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	assert.True(t, strings.Contains(rec.Body.String(), "2 von 3"), "Man befindet sich auf Seite 2. Es sollte drei Seiten benötigt werden. ")

	//Setup the caller for jump back
	req = httptest.NewRequest("GET", "localhost:443/user/view/filterTermins?Seite=Zurueck", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec = httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	assert.True(t, strings.Contains(rec.Body.String(), "1 von 3"), "Man befindet sich auf Seite 1. Es sollte drei Seiten benötigt werden. ")

	//delete Test Termins
	teardownDeleteTermin()
}
func testViewManager_FvFilter(t *testing.T) {
	//Test Termin
	testTermin := initTestTermin(false)

	//Test-Setups
	setupCookie(true)
	setupHandler()

	reader := createReaderData("", testTermin.Title, testTermin.Description, "", "", "", "", "")
	//Setup the caller for jump forward
	req := httptest.NewRequest("GET", "localhost:443/user/view/filterTermins?"+reader, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	assert.True(t, strings.Contains(rec.Body.String(), ""), "Man befindet sich auf Seite 2. Es sollte drei Seiten benötigt werden. ")

	//Setup the caller for jump back
	req = httptest.NewRequest("GET", "localhost:443/user/view/filterTermins?Seite=Zurueck", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec = httptest.NewRecorder()
	req.AddCookie(cookie)

	//Execute
	vmHandler.ServeHTTP(rec, req)
	assert.True(t, strings.Contains(rec.Body.String(), "Spezifischer-Test-Termin"), "Ein Termin mit dem Titel sollte gefiltert worden sein")
	assert.True(t, strings.Contains(rec.Body.String(), "Test Description"), "Ein Termin mit der Beschreibung sollte gefiltert worden sein")

	//delete Test Termins
	teardownDeleteTermin()
}
func TestViewManagerHandler_ServeHTTP(t *testing.T) {

	t.Run("testRuns ViewManager_WrongCookie", testViewManager_WrongCookie)
	t.Run("testRuns ViewManager_NewViewManagerObject", testViewManager_NewViewManagerObject)

	//Testfälle bei Get-Request, um Termin zum Bearbeiten zu erhalten
	t.Run("testRuns ViewManager_EditTerminGetRequest", testViewManager_EditTerminGetRequest)
	t.Run("testRuns ViewManager_EditTerminGetRequestError", testViewManager_EditTerminGetRequestError)
	//Testfälle bei POST-request des bearbeitenden Termins
	t.Run("testRuns ViewManager_EditTerminPostRequest", testViewManager_EditTerminPostRequest)
	t.Run("testRuns ViewManager_EditTerminPostRequestError", testViewManager_EditTerminPostRequestError)

	//Testfälle um Terminvorschläge (shared = true ) zu löschen
	t.Run("testRuns ViewManager_DeleteSharedTermin", testViewManager_DeleteSharedTermin)
	t.Run("testRuns ViewManager_DeleteSharedTerminError", testViewManager_DeleteSharedTerminError)

	//Testfälle bei Erstellung eines Termins
	t.Run("testRuns ViewManager_NewCreateTerminRequest", testViewManager_NewCreateTerminRequest)
	t.Run("testRuns ViewManager_NewCreateTerminRequestError", testViewManager_NewCreateTerminRequestError)

	/*
		+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
		Tests, die in an den Table-Handler weitergeleitet werden
		+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	*/
	//Tests, um in der Tabellen-Ansicht einen Monta vor- oder zurück zu springen
	t.Run("testRuns ViewManager_TvJumpMonthBack", testViewManager_TvJumpMonthBack)
	t.Run("testRuns ViewManager_TvJumpMonthFor", testViewManager_TvJumpMonthFor)

	//Tests, um in der Tabellens-Ansicht einen Monat auszuwählen
	t.Run("testRuns ViewManager_TvSelectMonth", testViewManager_TvSelectMonth)
	t.Run("testRuns ViewManager_TvSelectMonthError", testViewManager_TvSelectMonthError)

	//Tests, um in der Tabellen-Ansicht ein Jahr vor- oder zurück zu springen
	t.Run("testRuns ViewManager_TvJumpYearBack", testViewManager_TvJumpYearBack)
	t.Run("testRuns ViewManager_TvJumpYearFor", testViewManager_TvJumpYearFor)

	//Tests, um in der Tabellen-Ansicht zum heutigen Datum zu springen
	t.Run("testRuns ViewManager_TvJumpToToday", testViewManager_TvJumpToToday)

	//Tests, falls Tabellen-Ansicht neu geladen wurde
	t.Run("testRuns ViewManager_TvReload", testViewManager_TvReload)

	/*
		+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
		Tests, die in an den List-Handler weitergeleitet werden
		+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	*/
	//Tests, um in der Listen-Ansicht Anzahl der Einträge und das angezeigte Datum auszuwählen
	t.Run("testRuns ViewManager_LvSelectDate", testViewManager_LvSelectDate)
	t.Run("testRuns ViewManager_LvSelectDateError", testViewManager_LvSelectDateError)
	t.Run("testRuns ViewManager_LvSelectEntriesPerPage", testViewManager_LvSelectEntriesPerPage)
	t.Run("testRuns ViewManager_LvSelectEntriesPerPageError", testViewManager_LvSelectEntriesPerPageError)

	//Tests, um in der Listen-Ansicht eine Seite Vor- oder Zurück zu springen
	t.Run("testRuns ViewManager_LvJumpPageForward", testViewManager_LvJumpPageForward)
	t.Run("testRuns ViewManager_LvJumpPageBack", testViewManager_LvJumpPageBack)

	//Tests, falls Listen-Ansicht neu geladen wurde
	t.Run("testRuns ViewManager_TvReload", testViewManager_LvReload)

	/*
		+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
		Tests, die in an den FilterView-Handler weitergeleitet werden
		+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	*/
	//Tests, um in der Filter-Ansicht die Anzahl an einträgen auszuwählen
	t.Run("testRuns ViewManager_FvSelectEntriesPerPage", testViewManager_FvSelectEntriesPerPage)
	t.Run("testRuns ViewManager_FvSelectEntriesPerPageError", testViewManager_FvSelectEntriesPerPageError)

	//Tests, um in der Filter-Ansicht eine Seite Vor- oder Zurück zu springen
	t.Run("testRuns ViewManager_FvJumpPageForward", testViewManager_FvJumpPageForward)
	t.Run("testRuns ViewManager_FvJumpPageBack", testViewManager_FvJumpPageBack)

	//Tests, um in der Filter-Ansicht zu Filtern
	t.Run("testRuns ViewManager_FvFilter", testViewManager_FvFilter)

}
