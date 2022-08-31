//go:build !linux

package runutil

func PidState(pid uint64) (ProcState, error) {
	return nil, ErrNotSupported
}
