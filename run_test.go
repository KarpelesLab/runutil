package runutil

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"testing"
)

func TestRun(t *testing.T) {
	res, err := RunGet("echo", "-n", "hello world")
	if err != nil {
		t.Errorf("failed to run test: %s", err)
		return
	}

	if string(res) != "hello world" {
		t.Errorf("invlid output, expected hello world, got %+v", res)
	}
}

func TestPipe(t *testing.T) {
	buf := []byte("hello this is a small step for go and a smaller string once compressed")

	res, err := RunPipe(bytes.NewReader(buf), "gzip", "-9")
	if err != nil {
		t.Errorf("failed to run test: %s", err)
		return
	}

	cBuf, err := ioutil.ReadAll(res)
	if err != nil {
		t.Errorf("failed to read: %s", err)
		return
	}

	// try to decompress now
	res, err = RunPipe(bytes.NewReader(cBuf), "gunzip")
	if err != nil {
		t.Errorf("failed to run test: %s", err)
		return
	}

	dBuf, err := ioutil.ReadAll(res)
	if err != nil {
		t.Errorf("failed to read: %s", err)
		return
	}

	if !bytes.Equal(buf, dBuf) {
		t.Errorf("invalid, both buffers should be equal")
	}
}

func TestPipeError(t *testing.T) {
	res, err := RunRead("/bin/sh", "-c", "echo -n this will echo something but then things will go wrong; exit 42")
	if err != nil {
		t.Errorf("failed to run test: %s", err)
		return
	}

	buf, err := ioutil.ReadAll(res)

	if string(buf) != "this will echo something but then things will go wrong" {
		t.Errorf("failed, buf did not contain the expected stuff")
	}
	if err == nil {
		t.Errorf("failed, the command was supposed to return an error but didn't")
		return
	}

	var e *exec.ExitError
	if !errors.As(err, &e) {
		t.Errorf("failed, the command was supposed to return an error of type exec.ExitError, got %T (%s)", err, err)
		return
	}
	if e.ProcessState.ExitCode() != 42 {
		t.Errorf("failed, the command was supposed to return error 42")
	}
}

func TestPipeErrorCascade(t *testing.T) {
	res, err := RunRead("/bin/sh", "-c", "echo -n this will echo something but then things will go wrong; exit 42")
	if err != nil {
		t.Errorf("failed to run test: %s", err)
		return
	}

	res2, err := RunPipe(res, "cat")
	if err != nil {
		t.Errorf("failed to run test: %s", err)
		return
	}

	buf, err := ioutil.ReadAll(res2)

	if string(buf) != "this will echo something but then things will go wrong" {
		t.Errorf("failed, buf did not contain the expected stuff")
	}
	if err == nil {
		// TODO check if exit status 1 ?
		t.Errorf("failed, the command was supposed to return an error but didn't")
	}

	var e *exec.ExitError
	if !errors.As(err, &e) {
		t.Errorf("failed, the command was supposed to return an error of type exec.ExitError, got %T (%s)", err, err)
		return
	}
	if e.ProcessState.ExitCode() != 42 {
		t.Errorf("failed, the command was supposed to return error 42")
	}
}

func TestComplex(t *testing.T) {
	buf := []byte("this string is going to be going through so many things, I'm afraid for it... haha!")

	var res io.Reader

	res, err := RunPipe(bytes.NewReader(buf), "gzip", "-9")
	if err != nil {
		t.Errorf("failed to run test: %s", err)
		return
	}

	res, err = gzip.NewReader(res)
	if err != nil {
		t.Errorf("failed to run test: %s", err)
		return
	}

	res, err = RunPipe(res, "base64")
	if err != nil {
		t.Errorf("failed to run test: %s", err)
		return
	}

	res = base64.NewDecoder(base64.StdEncoding, res)

	final, err := ioutil.ReadAll(res)
	if err != nil {
		t.Errorf("failed to run test: %s", err)
		return
	}

	if !bytes.Equal(final, buf) {
		t.Errorf("failed to run test: should have been equal, instead got %s", final)
	}
}

func TestRemote(t *testing.T) {
	// https://cdn.kernel.org/pub/linux/kernel/v5.x/linux-5.12.10.tar.xz
	resp, err := http.Get("https://cdn.kernel.org/pub/linux/kernel/v5.x/linux-5.12.10.tar.xz")
	if err != nil {
		t.Errorf("failed to run test: %s", err)
		return
	}
	defer resp.Body.Close()

	// pipe this to xz
	res, err := RunPipe(resp.Body, "xz", "--decompress")
	if err != nil {
		t.Errorf("failed to run test: %s", err)
		return
	}

	// read some bytes
	buf := make([]byte, 17)
	_, err = res.Read(buf)
	if err != nil {
		t.Errorf("failed to run test: %s", err)
		return
	}

	// force close of res so we stop download
	err = res.Close()
	// err = signal: broken pipe
	if err.Error() != "signal: broken pipe" {
		t.Errorf("error: was expecting broken pipe error, got %s", err)
	}

	if string(buf) != "pax_global_header" {
		t.Errorf("failed to run test: invalid output (unexpected result)")
	}
}
