package stuff

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadRandomFile(t *testing.T) {
	// create a temporary file
	dir := t.TempDir()
	path := filepath.Join(dir, "idk.gif")
	if _, err := os.Create(path); err != nil {
		t.Fail()
	} else {
		os.WriteFile(path, make([]byte, 10), fs.ModeAppend)
	}
	// return it
	pic, err := ReadRandomFile(dir)
	assert.Nil(t, err)
	_, err = pic.Read(make([]byte, 1))
	assert.Nil(t, err)
}
