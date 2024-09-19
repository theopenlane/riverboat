package jobs

import "errors"

var (
	// ErrMissingTestDir is the error for missing test directory
	ErrMissingTestDir = errors.New("missing test directory in email config")
	// ErrMissingToken is the error for missing token
	ErrMissingToken = errors.New("missing resend api token, set token or use dev mode")
)
