package runutil

import "errors"

var (
	ErrCommandMissing = errors.New("command is missing")
	ErrNotSupported   = errors.New("operation not supported on this platform")
)
