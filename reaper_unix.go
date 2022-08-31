//go:build !windows

package runutil

import (
	"errors"
	"syscall"
)

func Reap() error {
	for {
		wpid, err := syscall.Wait4(-1, nil, syscall.WNOHANG, nil)
		if err != nil {
			if errors.Is(err, syscall.ECHILD) {
				// no more children to reap
				return nil
			}
			return err
		}
		if wpid == 0 {
			// probably returning ECHILD
			return nil
		}
		//log.Printf("main: clearing zombie process with pid %d", wpid)
	}
}
