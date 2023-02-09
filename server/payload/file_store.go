package payload

import (
	"errors"
	"fmt"
	"os"
)

func createFile(name string) error {
	_, err := os.Stat(name)
	if err == nil {
		return fmt.Errorf("file %v already exists", name)
	} else if errors.Is(err, os.ErrNotExist) {
		err = os.WriteFile("filename.txt", []byte("Hello"), 0644)
	}

	return err
}

func writeFile(name string, b []byte) error {
	f, err := os.OpenFile(name, os.O_APPEND, 0644)
	defer f.Close()
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}

  return nil
}
