package operations

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
)

func ReadRandomFile(directory string) (io.Reader, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	if len(files) < 1 {
		return nil, fmt.Errorf("%s is an empty directory", directory)
	}
	file := files[rand.Intn(len(files))]
	path := filepath.Join(directory, file.Name())
	reader, err := os.Open(path)
	return reader, err
}
