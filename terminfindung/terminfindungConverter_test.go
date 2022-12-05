/*
@author: 2447899 8689159 3000685
*/
package terminfindung

import (
	"DHBW_GO_Projekt/dateisystem"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestSaveSharedTermineToDisk_RightFunction
// tests that if it saves it saves to the right versio
func TestSaveSharedTermineToDisk_RightFunction(t *testing.T) {
	allTermine.shared = make(map[string]TerminFindung)
	allTermine.links = make(map[string]string)
	user := "test"
	//create termin
	termin := dateisystem.CreateNewTermin("Test", "Test Description", dateisystem.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.UTC),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.UTC),
		true, "test")
	//should return error if user is empty
	//create an appointment and a new proposed Date
	terminId, _ := CreateSharedTermin(&termin, &user)
	startDate := time.Date(2022, 12, 10, 12, 0, 0, 1, time.UTC)
	endDate := time.Date(2022, 12, 10, 12, 0, 0, 0, time.UTC)
	CreateNewProposedDate(startDate, endDate, &user, &terminId, false)
	CreatePerson(&user, &terminId, &user)
	path, err := filepath.Abs("../")
	assert.Equal(t, err, nil)
	err = SaveSharedTermineToDisk(&path)
	assert.Equal(t, nil, err)
	//check that it worked and read the path again
	path, err = filepath.Abs("../data/shared-termin/shared-termin-data.json")
	maps := make(map[string]TerminFindung)
	assert.Equal(t, nil, err)
	//read file and see if i have the rights
	file, err := os.ReadFile(path)
	assert.Equal(t, nil, err)
	err = json.Unmarshal(file, &maps)
	//should perfectly unmarshall it and then it should be one
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(maps))

	// ------------------------------------------------
	//trying links
	//it should start it and there should be one link in it because we created one person
	path2, err := filepath.Abs("../data/shared-termin/links.json")
	mapsLinks := make(map[string]string)
	assert.Equal(t, nil, err)
	//read file and see if i have the rights
	file, err = os.ReadFile(path2)
	assert.Equal(t, nil, err)
	err = json.Unmarshal(file, &mapsLinks)
	//should perfectly unmarshall it and then it should be one
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(maps))

	err = os.Remove(path2)
	assert.Equal(t, nil, err)
	err = os.Remove(path)
	assert.Equal(t, nil, err)
	dateisystem.DeleteAll(dateisystem.GetTermine(user), user)
}

// TestLoadDataToSharedTermin_RightShared
// This functions checks that if everything is created it should
// save in the right amount
func TestLoadDataToSharedTermin_RightShared(t *testing.T) {
	allTermine.shared = make(map[string]TerminFindung)
	allTermine.links = make(map[string]string)
	user := "test"
	//create termin
	termin := dateisystem.CreateNewTermin("Test", "Test Description", dateisystem.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.UTC),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.UTC),
		true, "test")
	//should return error if user is empty
	//create an appointment and a new proposed Date
	terminId, _ := CreateSharedTermin(&termin, &user)
	startDate := time.Date(2022, 12, 10, 12, 0, 0, 1, time.UTC)
	endDate := time.Date(2022, 12, 10, 12, 0, 0, 0, time.UTC)
	CreateNewProposedDate(startDate, endDate, &user, &terminId, false)
	CreatePerson(&user, &terminId, &user)
	path, _ := filepath.Abs("../")
	SaveSharedTermineToDisk(&path)
	err := LoadDataToSharedTermin(&path)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(allTermine.shared))
	assert.Equal(t, 1, len(allTermine.links))
	path2, err := filepath.Abs("../data/shared-termin/links.json")
	path, err = filepath.Abs("../data/shared-termin/shared-termin-data.json")
	err = os.Remove(path2)
	assert.Equal(t, nil, err)
	err = os.Remove(path)
	assert.Equal(t, nil, err)
	dateisystem.DeleteAll(dateisystem.GetTermine(user), user)
}

// TestSaveSharedTermineToDisk_EmptyShared
// checks that on saving an empty shared and links both files are empty
func TestSaveSharedTermineToDisk_EmptyShared(t *testing.T) {
	allTermine.shared = make(map[string]TerminFindung)
	allTermine.links = make(map[string]string)
	path, _ := filepath.Abs("../")
	err := SaveSharedTermineToDisk(&path)
	assert.Equal(t, nil, err)
	err = LoadDataToSharedTermin(&path)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(allTermine.shared))
	assert.Equal(t, 0, len(allTermine.links))
	path2, err := filepath.Abs("../data/shared-termin/links.json")
	path, err = filepath.Abs("../data/shared-termin/shared-termin-data.json")
	err = os.Remove(path2)
	assert.Equal(t, nil, err)
	err = os.Remove(path)
	assert.Equal(t, nil, err)

}

// TestLoadDataToSharedTermin_FileNotExisting
// tests that if it is not existing it hsould be automatically not existing
func TestLoadDataToSharedTermin_FileNotExisting(t *testing.T) {
	path, _ := filepath.Abs("../")
	err := LoadDataToSharedTermin(&path)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, true, strings.Contains(err.Error(), "coudn't read files"))
	assert.Equal(t, 0, len(allTermine.shared))
}

func TestTerminFindung_ConvertAdminToHTML_right(t *testing.T) {
	allTermine.shared = make(map[string]TerminFindung)
	allTermine.links = make(map[string]string)
	user := "test"
	//create termin
	termin := dateisystem.CreateNewTermin("Test", "Test Description", dateisystem.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.UTC),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.UTC),
		true, "test")
	//should return error if user is empty
	//create an appointment and a new proposed Date
	terminId, _ := CreateSharedTermin(&termin, &user)
	startDate := time.Date(2022, 12, 10, 12, 0, 0, 1, time.UTC)
	endDate := time.Date(2022, 12, 10, 12, 0, 0, 0, time.UTC)
	CreateNewProposedDate(startDate, endDate, &user, &terminId, false)
	CreatePerson(&user, &terminId, &user)
	name := "abcd"
	CreatePerson(&name, &terminId, &user)
	terminFind, _ := GetTerminFromShared(&user, &terminId)
	rightHtml := terminFind.ConvertAdminToHTML()
	assert.Equal(t, 2, len(rightHtml.Persons))
	assert.Equal(t, 1, len(rightHtml.VorschlagTermine))
	assert.Equal(t, false, rightHtml.IsLocked)
}

func TestTerminFindung_ConvertUserSiteToRightHTML(t *testing.T) {
	allTermine.shared = make(map[string]TerminFindung)
	allTermine.links = make(map[string]string)
	user := "test"
	//create termin
	termin := dateisystem.CreateNewTermin("Test", "Test Description", dateisystem.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.UTC),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.UTC),
		true, "test")
	//should return error if user is empty
	//create an appointment and a new proposed Date
	terminId, _ := CreateSharedTermin(&termin, &user)
	startDate := time.Date(2022, 12, 10, 12, 0, 0, 1, time.UTC)
	endDate := time.Date(2022, 12, 10, 12, 0, 0, 0, time.UTC)
	CreateNewProposedDate(startDate, endDate, &user, &terminId, false)
	CreatePerson(&user, &terminId, &user)
	name := "abcd"
	person, _ := CreatePerson(&name, &terminId, &user)
	//get the apikey
	apikey := person[7:]
	terminFind, _ := GetTerminFromShared(&user, &terminId)
	rightHtml := terminFind.ConvertUserSiteToRightHTML(&name, &apikey)
	//should convert it correctly
	assert.Equal(t, apikey, rightHtml.APIKey)
	assert.Equal(t, name, rightHtml.User)
	assert.Equal(t, 0, rightHtml.ToVotes[terminId].Votedfor)
}

func TestTerminFindung_ConvertAdminToHTML_DateSelected(t *testing.T) {
	allTermine.shared = make(map[string]TerminFindung)
	allTermine.links = make(map[string]string)
	user := "test"
	//create termin
	termin := dateisystem.CreateNewTermin("Test", "Test Description", dateisystem.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.UTC),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.UTC),
		true, "test")
	//should return error if user is empty
	//create an appointment and a new proposed Date
	terminId, _ := CreateSharedTermin(&termin, &user)
	startDate := time.Date(2022, 12, 10, 12, 0, 0, 1, time.UTC)
	endDate := time.Date(2022, 12, 10, 12, 0, 0, 0, time.UTC)
	CreateNewProposedDate(startDate, endDate, &user, &terminId, false)
	CreatePerson(&user, &terminId, &user)
	name := "abcd"
	CreatePerson(&name, &terminId, &user)
	VoteForDay(&terminId, &user, &name, &allTermine.shared[user+"|"+terminId].VorschlagTermine[0].ID, true)
	err := SelectDate(&allTermine.shared[user+"|"+terminId].VorschlagTermine[0].ID, &terminId, &user)
	assert.Equal(t, nil, err)
	terminFind, _ := GetTerminFromShared(&user, &terminId)
	rightHtml := terminFind.ConvertAdminToHTML()
	assert.Equal(t, true, rightHtml.IsLocked)
}
