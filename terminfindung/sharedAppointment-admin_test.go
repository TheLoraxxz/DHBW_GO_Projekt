package terminfindung

import (
	"DHBW_GO_Projekt/dateisystem"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestCreateSharedTermin
// tests that the shared termin is created correctly without an error and then adds it accordingly
func TestCreateSharedTermin_RightInput(t *testing.T) {
	// reset to zero
	allTermine.shared = make(map[string]TerminFindung)
	assert.Equal(t, 0, len(allTermine.shared))
	user := "test"
	//create termin
	termin := dateisystem.CreateNewTermin("Test", "Test Description", dateisystem.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.UTC),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.UTC),
		user, "test")
	terminId, err := CreateSharedTermin(&termin, &user)
	assert.Equal(t, "test", terminId)
	assert.Equal(t, 1, len(allTermine.shared))

	//check that it creates automaticaally the first appointment
	uuid := user + "|" + terminId
	fmt.Println(allTermine.shared[uuid])
	assert.Equal(t, 1, len(allTermine.shared[uuid].VorschlagTermine))
	assert.Equal(t, nil, err)

}
func TestCreateSharedTerminWrongInput(t *testing.T) {
	allTermine.shared = make(map[string]TerminFindung)
	assert.Equal(t, 0, len(allTermine.shared))
	user := "test"
	//create termin
	termin := dateisystem.CreateNewTermin("Test", "Test Description", dateisystem.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.FixedZone("Berlin", 1)),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.FixedZone("Berlin", 1)),
		user, "test")
	user = ""
	//should return error if user is empty
	terminId, err := CreateSharedTermin(&termin, &user)
	assert.Equal(t, 0, len(allTermine.shared))
	assert.Equal(t, "", terminId)
	assert.NotEqual(t, nil, err)
	//should return error if terminId is not set
	termin.ID = ""
	user = "test"
	terminId, err = CreateSharedTermin(&termin, &user)
	assert.Equal(t, 0, len(allTermine.shared))
	assert.Equal(t, "", terminId)
	assert.NotEqual(t, nil, err)
}

func TestCreateNewProposedDateRight_SameDate(t *testing.T) {
	allTermine.shared = make(map[string]TerminFindung)
	assert.Equal(t, 0, len(allTermine.shared))
	user := "test"
	//create termin
	termin := dateisystem.CreateNewTermin("Test", "Test Description", dateisystem.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.UTC),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.UTC),
		user, "test")
	//should return error if user is empty
	//create an appointment and a new proposed Date
	terminId, _ := CreateSharedTermin(&termin, &user)
	startDate := time.Date(2022, 12, 10, 12, 0, 0, 0, time.UTC)
	endDate := time.Date(2022, 12, 10, 12, 0, 0, 0, time.UTC)
	err := CreateNewProposedDate(startDate, endDate, &terminId, &user, false)
	//it should create it and the proposed date should be added
	assert.Equal(t, nil, err)
	userTerminId := user + "|" + terminId
	assert.Equal(t, 2, len(allTermine.shared[userTerminId].VorschlagTermine))
	proposedTermin := allTermine.shared[userTerminId].VorschlagTermine
	//everything should be eempty exept the start and enddate
	assert.Equal(t, true, proposedTermin[1].Date.Equal(startDate))
	assert.Equal(t, true, proposedTermin[1].EndDate.Equal(endDate))
	assert.Empty(t, proposedTermin[1].Title)
	assert.Empty(t, proposedTermin[1].Description)
	assert.NotEmpty(t, proposedTermin[1].ID)
}

func TestCreateNewProposedDate_StartDateAfterEnddate(t *testing.T) {
	allTermine.shared = make(map[string]TerminFindung)
	assert.Equal(t, 0, len(allTermine.shared))
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
	err := CreateNewProposedDate(startDate, endDate, &user, &terminId, false)
	assert.NotEqual(t, nil, err)
	userTerminId := user + "|" + terminId
	assert.Equal(t, 1, len(allTermine.shared[userTerminId].VorschlagTermine))

}

func TestCreatePerson(t *testing.T) {
	allTermine.shared = make(map[string]TerminFindung)
	assert.Equal(t, 0, len(allTermine.shared))
	user := "test"
	//create termin
	termin := dateisystem.CreateNewTermin("Test", "Test Description", dateisystem.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.UTC),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.UTC),
		user, "test2")
	terminId, _ := CreateSharedTermin(&termin, &user)
	_, err := CreatePerson(&user, &terminId, &user)
	// CreatePerson should return the right values --> the url with all the things and no error
	assert.Equal(t, err, nil)
	//the data blocks should be the right and a new user should be added
	assert.Equal(t, 1, len(allTermine.shared))
	assert.Equal(t, 1, len(allTermine.shared[user+"|"+terminId].Persons))
	assert.Equal(t, 1, len(allTermine.links))
	//assert.Equal(t, allTermine.shared[user+"|"+terminId].Persons[user].Url,)

}

func TestGetAllLinks(t *testing.T) {
	allTermine.shared = make(map[string]TerminFindung)
	assert.Equal(t, 0, len(allTermine.shared))
	user := "test"
	//create termin
	termin := dateisystem.CreateNewTermin("Test", "Test Description", dateisystem.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.UTC),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.UTC),
		user, "test2")
	terminId, _ := CreateSharedTermin(&termin, &user)
	_, err := CreatePerson(&user, &terminId, &user)
	assert.Equal(t, nil, err)
	users, err := GetAllLinks(&user, &terminId)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(users))

}

func TestSelectDate_RightInput(t *testing.T) {
	allTermine.shared = make(map[string]TerminFindung)
	assert.Equal(t, 0, len(allTermine.shared))
	user := "test"
	//create termin
	termin := dateisystem.CreateNewTermin("Test", "Test Description", dateisystem.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.UTC),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.UTC),
		user, "test2")
	//create shared appointment
	terminId, err := CreateSharedTermin(&termin, &user)
	assert.Equal(t, nil, err)
	//create proposed time
	startDate := time.Date(2022, 12, 10, 12, 0, 0, 0, time.UTC)
	endDate := time.Date(2022, 12, 12, 12, 0, 0, 0, time.UTC)
	//create another prop date, should work wihtout proble,s
	err = CreateNewProposedDate(startDate, endDate, &user, &terminId, false)
	assert.Equal(t, nil, err)
	//get termin for select date
	terminObj, _ := GetTerminFromShared(&user, &terminId)
	propDate := terminObj.VorschlagTermine[0].ID
	//select date should have the expected outcome
	err = SelectDate(&propDate, &terminId, &user)
	assert.Equal(t, nil, err)
	assert.Equal(t, propDate, allTermine.shared[user+"|"+terminId].FinalTermin.ID)
}

func TestSelectDate_WrongIDPropDatae(t *testing.T) {

}
