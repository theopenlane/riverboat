package jobs

import (
	"errors"
	"fmt"
)

var (
	// ErrMissingTestDir is the error for missing test directory
	ErrMissingTestDir = errors.New("missing test directory in email config")
	// ErrMissingToken is the error for missing token
	ErrMissingToken = errors.New("missing resend api token, set token or use dev mode")
	// ErrOpenlaneHostMissing is returned when the Openlane API host is missing
	ErrOpenlaneHostMissing = errors.New("openlane api host is required")
	// ErrOpenlaneTokenMissing is returned when the Openlane API token is missing
	ErrOpenlaneTokenMissing = errors.New("openlane api token is required")
)

// MissingRequiredArgError is returned when a required argument was not provided to a job
type MissingRequiredArgError struct {
	// Arg is the required argument that was not provided
	Arg string `json:"arg"`
	// Job is the job that requires the argument
	Job string `json:"job"`
}

// Error returns the MissingRequiredArgError in string format
func (e *MissingRequiredArgError) Error() string {
	return fmt.Sprintf("%s is required for the %s job", e.Arg, e.Job)
}

// newMissingRequiredArg returns an error for a missing required argument
func newMissingRequiredArg(field, job string) *MissingRequiredArgError {
	return &MissingRequiredArgError{
		Arg: field,
		Job: job,
	}
}
