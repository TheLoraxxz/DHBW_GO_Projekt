package terminfindung

import (
	"DHBW_GO_Projekt/dateisystem"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestCreateSharedTermin
// tests that the shared termin is created correctly without an error and then adds it accordingly
func TestCreateSharedTermin(t *testing.T) {
	// reset to zero
	allTermine.shared = make(map[string]TerminFindung)
	assert.Equal(t, 0, len(allTermine.shared))
	user := "admin"
	//create termin
	termin := dateisystem.CreateNewTermin("Test", "Test Description", dateisystem.Never,
		time.Date(2022, 12, 12, 12, 12, 0, 0, time.FixedZone("Berlin", 1)),
		time.Date(2022, 12, 13, 12, 12, 0, 0, time.FixedZone("Berlin", 1)),
		user, "test")
	id := CreateSharedTermin(&termin, &user)
	assert.Equal(t, "admin|test", id)
	assert.Equal(t, 1, len(allTermine.shared))

}
