//go:build !linux

package runutil

func PidState(name string) (ProcState, error) {
	return nil, ErrNotSupported
}
