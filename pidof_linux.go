package runutil

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func PidOf(name string) (res []int) {
	// for all processes (do not care about errors since l will be nil)
	l, _ := os.ReadDir("/proc")
	for _, proc := range l {
		pidstr := proc.Name()
		pid, err := strconv.ParseUint(pidstr, 10, 64)
		if err != nil {
			// not numeric. Don't care.
			continue
		}

		// let's first try to match based on /proc/%d/cmdline
		cmdline, err := ioutil.ReadFile(filepath.Join("/proc", pidstr, "cmdline"))
		if err == nil {
			cmdlineA := bytes.Split(cmdline, []byte{0})
			cmdname := filepath.Base(string(cmdlineA[0]))
			if cmdname == name {
				res = append(res, int(pid))
				continue
			}
		}

		// second, check exe symlink
		exe, err := os.Readlink(filepath.Join("/proc", pidstr, "exe"))
		if err != nil {
			// mmh?
			continue
		}

		if filepath.Base(exe) == name {
			res = append(res, int(pid))
		}
	}
	return
}
