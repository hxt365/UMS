package utils

type ValidationError struct {
	Err string
}

func (e ValidationError) Error() string {
	return e.Err
}

type AuthError struct {
	Err string
}

func (e AuthError) Error() string {
	return e.Err
}