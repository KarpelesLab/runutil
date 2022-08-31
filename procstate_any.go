package runutil

import "time"

type ProcState interface {
	IsRunning() bool
	Started() (time.Time, error)
}
