package terminfindung

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateSharedTermin(t *testing.T) {
	allTermine.shared = []TerminFindung{}
	assert.Equal(t, 0, len(allTermine.shared))

}
