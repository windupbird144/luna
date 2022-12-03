package operations

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test that a random file from a directory is read
func TestReadRandomFile(t *testing.T) {
	// Create a temporary directory and save the file "test.gif"
	dir := t.TempDir()
	path := filepath.Join(dir, "test.gif")
	os.Create(path)
	os.WriteFile(path, make([]byte, 10), fs.ModeAppend)

	// Get a random file from the directory and verify that you can read it
	pic, err := ReadRandomFile(dir)
	assert.Nil(t, err, "did not return a file from the directory")
	_, err = pic.Read(make([]byte, 1))
	assert.Nil(t, err, "cannot read the file that was returned")
}
