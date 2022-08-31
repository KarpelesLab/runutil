//go:build !linux

package runutil

func PidOf(name string) (res []int) {
	return nil
}
