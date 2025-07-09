//go:build !linux

package runutil

func PidOf(name string) (res []int) {
	return nil
}

func ArgsOf(pid int) ([]string, error) {
	return nil, nil
}
