package filesystem

import (
	"github.com/spf13/afero"
)

func NewLocalFs() afero.Fs {
	return afero.NewOsFs()
}
