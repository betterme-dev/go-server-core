package s3

type (
	Error struct {
		Message string
	}

	NotImplementedError struct {
		Error
	}

	FileClosedError struct {
		Error
	}
	UndefinedWhenceError struct {
		Error
	}
)

func (e Error) Error() string {
	return e.Message
}
