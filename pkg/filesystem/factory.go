package filesystem

import "github.com/spf13/afero"

type (
	FactoryInterface interface {
		New() afero.Fs
	}
)
