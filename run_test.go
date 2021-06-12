package runutil

import (
	"bytes"
	"io/ioutil"
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
	res, err := RunRead("/bin/sh", "-c", "echo -n this will echo something but then things will go wrong; exit 1")
	if err != nil {
		t.Errorf("failed to run test: %s", err)
		return
	}

	buf, err := ioutil.ReadAll(res)

	if string(buf) != "this will echo something but then things will go wrong" {
		t.Errorf("failed, buf did not contain the expected stuff")
	}
	if err == nil {
		// TODO check if exit status 1 ?
		t.Errorf("failed, the command was supposed to return an error but didn't")
	}
}
