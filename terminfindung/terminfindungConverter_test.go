package terminfindung

import (
	"DHBW_GO_Projekt/dateisystem"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestSaveSharedTermineToDisk_RightFunction
// tests that if it saves it saves to the right versio
func TestSaveSharedTermineToDisk_RightFunction(t *testing.T) {
	allTermine.shared = make(map[string]TerminFindung)
	user := "test"
	//create termin
	termin := dateisystem.CreateNewTermin("Test", "Test Description", dateisystem.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.UTC),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.UTC),
		user, "test")
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

func TestLoadDataToSharedTermin_RightShared(t *testing.T) {
	allTermine.shared = make(map[string]TerminFindung)
	allTermine.links = make(map[string]string)
	user := "test"
	//create termin
	termin := dateisystem.CreateNewTermin("Test", "Test Description", dateisystem.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.UTC),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.UTC),
		user, "test")
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
