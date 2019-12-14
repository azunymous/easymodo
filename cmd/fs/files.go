/*
Package fs provides the filesystem (via .Get()) and functions for adding to a set of in memory files
(Files interface) before writing them to the file system (with WriteAll())

For example, you would initially create a file map with NewFileMap(), add files to it via
Add(filename, content) and then write it to the file system with WriteAll(files, directory, subdir).

fs.Get() (returns an afero filesystem) should be used for creating directories and checking files or
directories exist.
*/
package fs

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"path"
)

// Interface Files is a store of files to be written to the file system
type Files interface {
	Add(string, string)
	WriteAll(directory, subDir string)
	GetFilenames() []string
}

type FileMap struct {
	files map[string]string
}

// NewFileMap creates a new file map that can be used for adding and writing files
func NewFileMap() Files {
	return &FileMap{make(map[string]string)}
}

// Add adds a new file to the file map
func (f *FileMap) Add(fileName, content string) {
	f.files[fileName] = content
}

func (f *FileMap) write(appFs afero.Fs, dir, subDir string) error {
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

// GetFilenames returns all filenames of a file map
func (f *FileMap) GetFilenames() []string {
	fileNames := make([]string, 0, len(f.files))
	for name := range f.files {
		fileNames = append(fileNames, name)
	}
	return fileNames
}

// WriteAll creates the given directory and writes all provided files to it.
func (f *FileMap) WriteAll(directory, subDir string) {
	_ = appFs.Mkdir(path.Join(directory, subDir), 0755)

	_ = f.write(appFs, directory, subDir)
}
