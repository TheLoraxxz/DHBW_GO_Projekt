package terminfindung

import (
	"DHBW_GO_Projekt/dateisystem"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestVoteForDayRightInput(t *testing.T) {
	allTermine.shared = make(map[string]TerminFindung)
	allTermine.links = make(map[string]string)
	user := "test"
	//create termin
	termin := dateisystem.CreateNewTermin("Test", "Test Description", dateisystem.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.UTC),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.UTC),
		user, "test2")
	//create shared appointment
	terminId, _ := CreateSharedTermin(&termin, &user)
	//create a person so it can be entered
	CreatePerson(&user, &terminId, &user)
	termine, _ := GetTerminFromShared(&user, &terminId)
	err := VoteForDay(&terminId, &user, &user, &termine.VorschlagTermine[0].ID, true)
	//should run without problem and the user should have the according entry
	assert.Equal(t, nil, err)
	assert.Equal(t, true, allTermine.shared[user+"|"+terminId].Persons[user].Votes[termine.VorschlagTermine[0].ID])

}

func TestGetTerminViaApiKey_RightInput(t *testing.T) {
	allTermine.shared = make(map[string]TerminFindung)
	allTermine.links = make(map[string]string)
	user := "test"
	//create termin
	termin := dateisystem.CreateNewTermin("Test", "Test Description", dateisystem.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.UTC),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.UTC),
		user, "test2")
	//create shared appointment
	terminId, _ := CreateSharedTermin(&termin, &user)
	//create a person so it can be entered
	person, err := CreatePerson(&user, &terminId, &user)
	assert.Equal(t, nil, err)
	assert.Equal(t, err, nil)
	key := person[7:]
	terminFinal, userFinal, err := GetTerminViaApiKey(&key)
	assert.Equal(t, nil, err)
	termine, _ := GetTerminFromShared(&user, &terminId)
	assert.Equal(t, user, userFinal)
	assert.Equal(t, termine.Info.ID, terminFinal.Info.ID)
	dateisystem.DeleteAll(dateisystem.GetTermine(user), user)

}

func TestGetTerminViaApiKey_WrongAPIKey(t *testing.T) {
	apikey := ""
	termin, user, err := GetTerminViaApiKey(&apikey)
	assert.Empty(t, termin)
	assert.Equal(t, user, "")
	assert.NotEqual(t, err, nil)
	assert.Equal(t, true, strings.Contains(err.Error(), "API Key is not defined"))
}

// TestGetTerminViaApiKey_EmptyObject checks if it returns an error
// if the object is empty (can happen if it will be deleted)
func TestGetTerminViaApiKey_EmptyObject(t *testing.T) {
	allTermine.shared = make(map[string]TerminFindung)
	apikey := "test"
	allTermine.shared[apikey] = TerminFindung{}
	allTermine.links[apikey] = apikey
	termin, user, err := GetTerminViaApiKey(&apikey)
	assert.NotEqual(t, nil, err)
	assert.Empty(t, termin)
	assert.Equal(t, user, "")
}

func TestVoteForDay(t *testing.T) {
	allTermine.shared = make(map[string]TerminFindung)
	allTermine.links = make(map[string]string)
	user := "test"
	//create termin
	termin := dateisystem.CreateNewTermin("Test", "Test Description", dateisystem.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.UTC),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.UTC),
		user, "test2")
	//create shared appointment
	terminId, _ := CreateSharedTermin(&termin, &user)
	//create a person so it can be tested
	CreatePerson(&user, &terminId, &user)
	person := "nichtPerson"
	//testing voted day if it is voting on the wrong day
	err := VoteForDay(&terminId, &user, &person, &terminId, true)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "cound't find username", err.Error())
}
