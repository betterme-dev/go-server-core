package filesystem

import (
	"github.com/spf13/afero"
)

type (
	LocalFsFactory struct{}
)

func (lfs LocalFsFactory) New() afero.Fs {
	return afero.NewOsFs()
}
