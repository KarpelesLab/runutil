package runutil

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	// CLK_TCK is a constant on Linux for all architectures except alpha and ia64.
	// See e.g.
	// https://git.musl-libc.org/cgit/musl/tree/src/conf/sysconf.c#n30
	// https://github.com/containerd/cgroups/pull/12
	// https://lore.kernel.org/lkml/agtlq6$iht$1@penguin.transmeta.com/
	_SYSTEM_CLK_TCK = 100
)

type LinuxProcState struct {
	Pid         int
	Comm        string // "bash", etc
	State       byte   // 'R', 'S', 'D', 'Z', 'T', 't', 'W', etc
	PPid        int    // parent pid
	PGrp        int    // group id
	Session     int    // session
	TtyNr       int    // tty number
	Tpgid       int
	Flags       uint   // PF_* flags
	Minflt      uint64 // faults
	Cminflt     uint64
	Majflt      uint64
	Cmajflt     uint64
	Utime       uint64
	Stime       uint64
	Cutime      uint64
	Cstime      uint64
	Priority    int64
	Nice        int64 // range 19 (low priority) to -20 (high priority).
	NumThreads  int64
	Itrealvalue int64  // The time in jiffies before the next SIGALRM
	StartTime   uint64 // The time the process started after system boot in clock ticks (_SYSTEM_CLK_TCK)
	Vsize       uint64
	RSS         int64
	RSSlim      uint64
}

func (s *LinuxProcState) IsRunning() bool {
	return s.State == 'R'
}

func (s *LinuxProcState) Started() (time.Time, error) {
	// we can find out the start time of a process by reading /proc/<pid>/stat and /proc/uptime
	uptime, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		return time.Time{}, err
	}
	uptimeA := strings.Fields(strings.TrimSpace(string(uptime)))
	// uptimeA[0] is current uptime in system time, in seconds (as a float)
	uptSec, err := strconv.ParseFloat(uptimeA[0], 64)
	if err != nil {
		return time.Time{}, err
	}

	// read process launch time
	startTime := float64(s.StartTime) / _SYSTEM_CLK_TCK

	// how long ago was this?
	procUptime := uptSec - startTime

	if procUptime < 0 {
		return time.Time{}, fmt.Errorf("unexpected negative uptime value startTime=%f uptime=%f", startTime, uptSec)
	}

	// compute response
	now := time.Now()
	res := now.Add(time.Duration(procUptime * -1 * float64(time.Millisecond)))

	return res, nil
}

func eatValue(in *[]string, out any) error {
	if len(*in) < 1 {
		return io.ErrUnexpectedEOF
	}
	// take first value from in, and try to put it into out
	v := (*in)[0]
	(*in) = (*in)[1:]

	var err error

	switch x := out.(type) {
	case *byte:
		// expect only one single ASCII char
		if len(v) != 1 {
			return errors.New("bad length for byte")
		}
		(*x) = v[0]
		return nil
	case *int64:
		(*x), err = strconv.ParseInt(v, 0, 64)
		return err
	case *int:
		t, err := strconv.ParseInt(v, 0, 64)
		(*x) = int(t)
		return err
	case *uint64:
		(*x), err = strconv.ParseUint(v, 0, 64)
		return err
	case *uint:
		t, err := strconv.ParseUint(v, 0, 64)
		(*x) = uint(t)
		return err
	default:
		return fmt.Errorf("unsupported write type %T", out)
	}
}

func (s *LinuxProcState) parse(data string) error {
	// see: https://man7.org/linux/man-pages/man5/procfs.5.html
	// 3947 (bash test) S 3799 3947 3799 34828 3964 4194304 547 1212 0 2 0 0 2 0 20 0 1 0 806660689 10452992 1039 18446744073709551615 94202653388800 94202653965661 140728351794288 0 0 0 65536 3686404 1266761467 1 0 0 17 11 0 0 0 0 0 94202654164112 94202654186268 94202661220352 140728351795906 140728351795918 140728351795918 140728351801324 0

	dec := strings.Fields(data)
	if len(dec) < 52 {
		return errors.New("not enough fields in proc state (invalid format)")
	}

	err := eatValue(&dec, &s.Pid) // 1
	if err != nil {
		return err
	}

	// read Comm (2)
	comm := dec[0]
	if len(comm) < 1 || comm[0] != '(' {
		return errors.New("invalid proc state format at comm")
	}
	comm = comm[1:] // strip '('
	dec = dec[1:]
	for comm[len(comm)-1] != ')' {
		comm += " " + dec[0]
		dec = dec[1:]
	}
	comm = comm[:len(comm)-1] // strip ')'
	s.Comm = comm

	eatValue(&dec, &s.State)       // 3
	eatValue(&dec, &s.PPid)        // 4
	eatValue(&dec, &s.PGrp)        // 5
	eatValue(&dec, &s.Session)     // 6
	eatValue(&dec, &s.TtyNr)       // 7
	eatValue(&dec, &s.Tpgid)       // 8
	eatValue(&dec, &s.Flags)       // 9
	eatValue(&dec, &s.Minflt)      // 10
	eatValue(&dec, &s.Cminflt)     // 11
	eatValue(&dec, &s.Majflt)      // 12
	eatValue(&dec, &s.Cmajflt)     // 13
	eatValue(&dec, &s.Utime)       // 14
	eatValue(&dec, &s.Stime)       // 15
	eatValue(&dec, &s.Cutime)      // 16
	eatValue(&dec, &s.Cstime)      // 17
	eatValue(&dec, &s.Priority)    // 18
	eatValue(&dec, &s.Nice)        // 19
	eatValue(&dec, &s.NumThreads)  // 20
	eatValue(&dec, &s.Itrealvalue) // 21
	eatValue(&dec, &s.StartTime)   // 22
	eatValue(&dec, &s.Vsize)       // 23
	eatValue(&dec, &s.RSS)         // 24
	eatValue(&dec, &s.RSSlim)      // 25

	return err
}

func PidState(pid uint64) (ProcState, error) {
	return LinuxPidState(pid)
}

func LinuxPidState(pid uint64) (*LinuxProcState, error) {
	pStat, err := ioutil.ReadFile(filepath.Join("/proc", strconv.Itoa(int(pid)), "stat"))
	if err != nil {
		// process not running?
		return nil, err
	}

	state := &LinuxProcState{}
	return state, state.parse(string(pStat))
}
