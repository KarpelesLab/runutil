package runutil

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
)

// Run is a very simple invokation of command run, with output forwarded to stdout. This will wait for the command to complete.
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

// RunWrite executes the command and passes r as its input, waiting for it to complete.
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

// RunRead executes the command in background and returns its output as a stream.
// Close the stream to kill the command and release its resources.
func RunRead(arg ...string) (Pipe, error) {
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

	re := newProcessPipe(r, c)

	return re, nil
}

// RunPipe runs a command in background, connecting both ends
func RunPipe(r io.Reader, arg ...string) (Pipe, error) {
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

	re := newProcessPipe(out, c)

	return re, nil
}

// RunGet executes the command and returns its output as a buffer after it completes.
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
	defer r.Close() // close pipe after we finish reading

	// parse
	dec := json.NewDecoder(r)
	return dec.Decode(obj)
}
