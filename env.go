package runutil

import (
	"os"
	"os/exec"
	"path"
	"slices"
	"strings"
)

type Env []string

// SysEnv returns a Env object with os.Environ() loaded
func SysEnv() Env {
	return Env(os.Environ())
}

// NewEnv returns an empty env with only HOME, PATH set
func NewEnv(home string) Env {
	usr := "root"
	if home != "/" {
		usr = path.Base(home)
	}
	return Env{"USER=" + usr, "PWD=/", "HOME=" + home, "PATH=/usr/sbin:/usr/bin:/sbin:/bin"}
}

// Run is a very simple invokation of command run, with output forwarded to stdout. This will wait for the command to complete.
func (e Env) Run(arg ...string) error {
	if len(arg) == 0 {
		return ErrCommandMissing
	}

	cmd, err := exec.LookPath(arg[0])
	if err != nil {
		return err
	}

	c := &exec.Cmd{
		Path:   cmd,
		Args:   arg,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    []string(e),
	}

	return c.Run()
}

func (e Env) Join(others ...Env) Env {
	n := e
	for _, x := range others {
		n = append(n, x...)
	}
	return n
}

func (e *Env) Set(k, v string) {
	k2 := k + "="
	for n, s := range *e {
		if strings.HasPrefix(s, k2) {
			(*e)[n] = k2 + v
			return
		}
	}

	// not found, append
	*e = append(*e, k2+v)
}

func (e Env) Get(k string) string {
	k2 := k + "="
	for _, s := range e {
		if strings.HasPrefix(s, k2) {
			return s[len(k2):]
		}
	}
	return ""
}

func (e *Env) Unset(k string) {
	k2 := k + "="
	for n, s := range *e {
		if strings.HasPrefix(s, k2) {
			*e = slices.Delete(*e, n, n)
		}
	}
}

func (e Env) Contains(k string) bool {
	k2 := k + "="
	for _, s := range e {
		if strings.HasPrefix(s, k2) {
			return true
		}
	}
	return false
}
