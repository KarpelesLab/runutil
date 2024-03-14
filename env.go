package runutil

import (
	"os"
	"path"
	"slices"
	"strings"
)

type Env []string

// SysEnv returns a nil Env object that would mean pointing to the OS's environ
func SysEnv() Env {
	return nil
}

func sysEnv() Env {
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

// Join appends in one shot multiple environments together. No check for duplicates is done
func (e Env) Join(others ...Env) Env {
	n := e
	if n == nil {
		n = sysEnv()
	}

	for _, x := range others {
		n = append(n, x...)
	}
	return n
}

// Set sets the given variable in the env
func (e *Env) Set(k, v string) {
	if *e == nil {
		*e = sysEnv()
	}

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

// Get returns the requested value, or empty if none was found
func (e Env) Get(k string) string {
	if e == nil {
		return os.Getenv(k)
	}

	k2 := k + "="
	for _, s := range e {
		if strings.HasPrefix(s, k2) {
			return s[len(k2):]
		}
	}
	return ""
}

// Unset removes any instances of a given value from the env
func (e *Env) Unset(k string) {
	if *e == nil {
		*e = sysEnv()
	}

	k2 := k + "="
	for n, s := range *e {
		if strings.HasPrefix(s, k2) {
			*e = slices.Delete(*e, n, n)
		}
	}
}

// Contains checks if the given env contains any value k, and confirm if the value exists or not
func (e Env) Contains(k string) bool {
	if e == nil {
		e = sysEnv()
	}

	k2 := k + "="
	for _, s := range e {
		if strings.HasPrefix(s, k2) {
			return true
		}
	}
	return false
}

// Dedup returns a copy of e with any duplicate value removed
func (e Env) Dedup() Env {
	v := make(map[string]bool)
	ln := len(e)
	ne := make(Env, ln)
	p := ln
	for x := range e {
		s := e[ln-x-1] // start from the end
		// extract key part of the string
		eqp := strings.IndexByte(s, '=')
		k := s
		if eqp != -1 {
			k = k[:eqp]
		}
		// check if we already saw this key
		if _, found := v[k]; found {
			// duplicate
			continue
		}
		// mark it
		v[k] = true
		// store in new env (from the end)
		ne[p-1] = s
		p -= 1
	}

	return ne[p:]
}
