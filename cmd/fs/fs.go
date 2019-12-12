package fs

import (
	"github.com/spf13/afero"
)

var appFs = afero.NewOsFs()

func Get() afero.Fs {
	return appFs
}

func SetFs() {
	appFs = afero.NewMemMapFs()
}
