package dateisystem

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func init() {
	CreateNewTermin("test", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik", "0")
}

func TestGetDirectory(t *testing.T) {
	directory := GetDirectory("mik")
	directory = filepath.Dir(directory)
	con, _ := os.Getwd()

	assert.Equal(t, con, directory)
}

func TestGetFile(t *testing.T) {
	file := getFile("test", "mik")
	_, file = filepath.Split(file)

	assert.Equal(t, "test.json", file)
}
