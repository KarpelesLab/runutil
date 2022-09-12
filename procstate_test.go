//go:build linux

package runutil

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestLinuxParse(t *testing.T) {
	v := "3947 (bash test) S 3799 3947 3799 34828 3964 4194304 547 1212 0 2 0 0 2 0 20 0 1 0 806660689 10452992 1039 18446744073709551615 94202653388800 94202653965661 140728351794288 0 0 0 65536 3686404 1266761467 1 0 0 17 11 0 0 0 0 0 94202654164112 94202654186268 94202661220352 140728351795906 140728351795918 140728351795918 140728351801324 0\n"

	s := &LinuxProcState{}

	err := s.parse(v)

	if err != nil {
		t.Fatalf("test failed: %s", err)
	}

	log.Printf("test = %+v", s)
	// {Pid:3947 Comm:bash test State:'S' PPid:3799 PGrp:3947 Session:3799 TtyNr:34828 Tpgid:3964 Flags:4194304 Minflt:547 Cminflt:1212 Majflt:0 Cmajflt:2 Utime:0 Stime:0 Cutime:2 Cstime:0 Priority:20 Nice:0 NumThreads:1 Itrealvalue:0 StartTime:806660689 Vsize:10452992 RSS:1039 RSSlim:18446744073709551615}

	if s.Pid != 3947 {
		t.Fatal("invalid pid in decoded state")
	}
	if s.Comm != "bash test" {
		t.Fatal("invalid comm in decoded state")
	}
	if s.State != 'S' {
		t.Fatal("invalid state in decoded state")
	}
	if s.RSS != 1039 {
		t.Fatal("invalid RSS in decoded state")
	}

	testVal := 30 * time.Millisecond

	time.Sleep(testVal)

	// read own value
	s2, err := PidState(uint64(os.Getpid()))
	if err != nil {
		t.Fatalf("test failed: %s", err)
	}

	tv, err := s2.Started()
	if err != nil {
		t.Fatalf("test failed: %s", err)
	}

	ago := time.Since(tv).Truncate(10 * time.Millisecond)

	if ago != testVal {
		t.Fatalf("expected %s uptime of process, but got %s", testVal, ago)
	}
}
