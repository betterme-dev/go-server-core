package s3

type (
	FSError struct {
		Message string
	}

	NotImplementedError struct {
		FSError
	}

	FileClosedError struct {
		FSError
	}
	UndefinedWhenceError struct {
		FSError
	}
)

func (e FSError) Error() string {
	return e.Message
}
