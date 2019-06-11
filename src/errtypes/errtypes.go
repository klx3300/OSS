package errtypes

import "errors"

// these represents some basic error types
var (
	ETEST    = errors.New("test error")
	EEXIST   = errors.New("already exists")
	ENEXIST  = errors.New("not exist")
	EAGAIN   = errors.New("try again")
	EBADF    = errors.New("bad descriptor/main object")
	ECORRUPT = errors.New("data corruption")
)

// GenTestError generates a test error to test if these stuff works
func GenTestError() error {
	return ETEST
}
