package fs

import (
	"github.com/spf13/afero"
)

var appFs = afero.NewOsFs()

// Get returns the afero filesystem
func Get() afero.Fs {
	return appFs
}

// SetFs sets the filesystem to an in memory file system instead. Useful for testing.
func SetFs() {
	appFs = afero.NewMemMapFs()
}
