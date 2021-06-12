package runutil

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
)

// Run is a very simple invokation of command run, with output forwarded to stdout
func Run(arg ...string) error {
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
	}

	return c.Run()
}

// RunWrite executed the command and passes r as its input
func RunWrite(r io.Reader, arg ...string) error {
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
		Stdin:  r,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	return c.Run()
}

// RunRead executes and returns the command's output as a stream
func RunRead(arg ...string) (io.ReadCloser, error) {
	if len(arg) == 0 {
		return nil, ErrCommandMissing
	}

	cmd, err := exec.LookPath(arg[0])
	if err != nil {
		return nil, err
	}

	c := &exec.Cmd{
		Path:   cmd,
		Args:   arg,
		Dir:    "/",
		Stderr: os.Stderr,
	}

	r, err := c.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err = c.Start(); err != nil {
		return nil, err
	}

	re := newErrorPipe(r, c.Wait)

	return re, nil
}

// RunPipe runs a command, connecting both ends
func RunPipe(r io.Reader, arg ...string) (io.ReadCloser, error) {
	if len(arg) == 0 {
		return nil, ErrCommandMissing
	}

	cmd, err := exec.LookPath(arg[0])
	if err != nil {
		return nil, err
	}

	c := &exec.Cmd{
		Path:   cmd,
		Args:   arg,
		Dir:    "/",
		Stdin:  r,
		Stderr: os.Stderr,
	}

	out, err := c.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err = c.Start(); err != nil {
		return nil, err
	}

	re := newErrorPipe(out, c.Wait)

	return re, nil
}

// RunGet executes the command and returns its output as a buffer
func RunGet(arg ...string) ([]byte, error) {
	if len(arg) == 0 {
		return nil, ErrCommandMissing
	}

	cmd, err := exec.LookPath(arg[0])
	if err != nil {
		return nil, err
	}

	c := &exec.Cmd{
		Path: cmd,
		Args: arg,
	}

	return c.Output()
}

// RunJson executes the command and applies its output to the specified object, parsing json data
func RunJson(obj interface{}, arg ...string) error {
	r, err := RunRead(arg...)
	if err != nil {
		return err
	}

	// parse
	dec := json.NewDecoder(r)
	return dec.Decode(obj)
}
