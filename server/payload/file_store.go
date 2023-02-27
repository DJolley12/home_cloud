package payload

import (
	"errors"
	"fmt"
	"os"
)

func createFile(name string) (*os.File, error) {
	_, err := os.Stat(name)
	if err == nil {
		return nil, fmt.Errorf("file %v already exists", name)
	} else if errors.Is(err, os.ErrNotExist) {
		return os.Create(name)
	}

	return nil, err
}

func writeFile(f *os.File, b []byte) (int, error) {
	return f.Write(b)
}
