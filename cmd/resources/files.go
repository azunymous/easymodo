package resources

import (
	log "github.com/sirupsen/logrus"
	"os"
)

type Files interface {
	Add(string, string)
	Write() error
	Get() map[string]string
}

type FileMap struct {
	files map[string]string
}

func NewFileMap() *FileMap {
	return &FileMap{make(map[string]string)}
}

func (f *FileMap) Add(fileName, content string) {
	f.files[fileName] = content
}

func (f *FileMap) Write() error {
	for fileName, content := range f.files {
		f, err := os.Create(fileName)

		if err != nil {
			log.Fatalf("Could not create %s file: %v", fileName, err)
		}

		n, err := f.WriteString(content)
		log.Debugf("Wrote %d bytes to %s", n, fileName)

		if err != nil {
			log.Fatalf("Failed writing content to %s", fileName)
		}
	}
	return nil
}

func (f *FileMap) Get() map[string]string {
	return f.files
}
