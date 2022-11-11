package dateisystem

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func init() {
	CreateNewTermin("test", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), "mik")
}

func TestGetDirectory(t *testing.T) {
	directory := getDirectory("mik")
	directory = filepath.Dir(directory)
	con, _ := os.Getwd()

	assert.Equal(t, con, directory)
}

func TestGetFilenameByTitle(t *testing.T) {
	file := getFileNameByTitle("test", "mik")
	_, file = filepath.Split(file)

	assert.Equal(t, "test.json", file)
}

func TestGetFilenameByTerminObj(t *testing.T) {
	ter := NewTerminObj("test", "test", WEEKLY, time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC), time.Date(2007, 3, 2, 15, 2, 5, 0, time.UTC))
	file := getFileNameByTerminObj(ter, "mik")
	_, file = filepath.Split(file)

	assert.Equal(t, "test.json", file)
}
