package terminfindung

import (
	"DHBW_GO_Projekt/dateisystem"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
	"time"
)

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
}
