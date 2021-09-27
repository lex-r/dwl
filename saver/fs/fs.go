package fs

import (
	"io"
	"log"
	"os"
)

// Saver is implementation of downloader.Saver interface.
type Saver struct {
	dir string
}

// NewSaver creates new Saver.
func NewSaver(dir string) *Saver {
	return &Saver{dir: dir}
}

// Save saves given data into file.
func (s *Saver) Save(fileName string, r io.Reader) error {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("error creating file for fileName %s: %s", fileName, err)
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		log.Printf("error saving file %s: %s", fileName, err)
		return err
	}

	err = f.Sync()
	if err != nil {
		log.Printf("error syncing file: %s", err)
		return err
	}

	return nil
}
