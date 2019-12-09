package resources

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"path"
)

type Files interface {
	Add(string, string)
	Write(filesystem afero.Fs, dir, subDir string) error
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

func (f *FileMap) Write(appFs afero.Fs, dir, subDir string) error {
	for fileName, content := range f.files {
		if content == "" {
			continue
		}

		f, err := appFs.Create(path.Join(dir, subDir, fileName))

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

func (f *FileMap) GetResources() []string {
	fileNames := make([]string, 0, len(f.files))
	for name := range f.files {
		fileNames = append(fileNames, name)
	}
	return fileNames
}
