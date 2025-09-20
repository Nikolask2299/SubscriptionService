package sql

type SqlError struct {
	error string
}
func NewSqlError(err error) SqlError {
	if err == nil {
		return SqlError{error: "Unknown error"}
	} else {
		return SqlError{error: err.Error()}
	}
}

func (e SqlError) Error() string {
	return e.error
}

func (e SqlError) Is(err error) bool {
	if err == nil {
		return false
	} else {
		return e.Error() == err.Error()
	}
}