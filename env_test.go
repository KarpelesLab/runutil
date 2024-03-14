package runutil_test

import (
	"strings"
	"testing"

	"github.com/KarpelesLab/runutil"
)

func dumpEnv(e runutil.Env) string {
	// return env as a simple string
	return strings.Join([]string(e), " ")
}

func checkEnv(t *testing.T, e runutil.Env, expect string) {
	x := dumpEnv(e)
	if x != expect {
		t.Errorf("env failure, expected %s but got %s", expect, x)
	}
}

func TestEnv(t *testing.T) {
	checkEnv(t, runutil.NewEnv("/"), "USER=root PWD=/ HOME=/ PATH=/usr/sbin:/usr/bin:/sbin:/bin")
	checkEnv(t, runutil.NewEnv("/home/linux"), "USER=linux PWD=/ HOME=/home/linux PATH=/usr/sbin:/usr/bin:/sbin:/bin")
	checkEnv(t, runutil.NewEnv("/", "FOO=bar"), "USER=root PWD=/ HOME=/ PATH=/usr/sbin:/usr/bin:/sbin:/bin FOO=bar")
	checkEnv(t, runutil.NewEnv("/", "FOO=bar", "FOO=baz"), "USER=root PWD=/ HOME=/ PATH=/usr/sbin:/usr/bin:/sbin:/bin FOO=baz")

	n := runutil.NewEnv("/")
	n.Set("PWD", "/tmp")
	checkEnv(t, n, "USER=root PWD=/tmp HOME=/ PATH=/usr/sbin:/usr/bin:/sbin:/bin")

	n = n.Join(runutil.Env{"HOME=/bar"})
	checkEnv(t, n, "USER=root PWD=/tmp HOME=/ PATH=/usr/sbin:/usr/bin:/sbin:/bin HOME=/bar")

	n = n.Dedup()
	checkEnv(t, n, "USER=root PWD=/tmp PATH=/usr/sbin:/usr/bin:/sbin:/bin HOME=/bar")
}
