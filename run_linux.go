package runutil

import "strings"

// Sh runs a shell command (linux only)
func Sh(cmd string) error {
	return Run("/bin/sh", "-c", cmd)
}

func ShQuote(s string) string {
	if len(s) == 0 {
		return "''"
	}
	return "'" + strings.Replace(s, "'", "'\"'\"'", -1) + "'"
}
