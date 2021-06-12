[![GoDoc](https://godoc.org/github.com/KarpelesLab/runutil?status.svg)](https://godoc.org/github.com/KarpelesLab/runutil)

# runutil

Various useful tools for running and pipe-ing stuff outside of Go.

This library makes it very easy to execute complex sequences of executables mixing both go-specific filters and commands.

For example it is possible to run the following:

```go
	res, err := RunPipe(input, "gzip", "-9")
	if err != nil {
		t.Errorf("failed to run test: %s", err)
		return
	}

	res, err = gzip.NewReader(res)
	if err != nil {
		t.Errorf("failed to run test: %s", err)
		return
	}
```

Reading from `res` in that example will return the exact same bytes as input, after having been compressed once, then decompressed.

If the command fails, the final Read() call will return the failure code, and allows correctly catching any problem (by default, go `os/exec` will only return the error when calling Wait(), which may result in errors not being catched).
