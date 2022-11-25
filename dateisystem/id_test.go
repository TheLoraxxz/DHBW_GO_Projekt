package dateisystem

//Mat-Nr. 8689159
import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetID(t *testing.T) {
	createID()
	n := getID()
	assert.Equal(t, n, getID())
}

func TestIncrementID(t *testing.T) {
	createID()
	n := getID()
	incrementID()
	assert.NotEqual(t, n, getID())
}

func TestDecrementID(t *testing.T) {
	createID()
	incrementID()
	n := getID()
	decrementID()
	assert.Equal(t, getID(), n-1)
}
