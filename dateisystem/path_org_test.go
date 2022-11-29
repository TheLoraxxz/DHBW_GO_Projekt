package dateisystem

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

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
