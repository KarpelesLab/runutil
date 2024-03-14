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
}
