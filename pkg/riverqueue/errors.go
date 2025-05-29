package riverqueue

import "errors"

// ErrConnectionURIRequired is returned when the connection URI is not provided
var ErrConnectionURIRequired = errors.New("connection URI is required to initialize the client")
