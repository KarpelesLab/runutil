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
