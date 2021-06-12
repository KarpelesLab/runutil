package runutil

import (
	"io"
	"sync"
)

type errorPipe struct {
	r io.ReadCloser
	e error
	o sync.Once
	c func() error
}

func newErrorPipe(r io.ReadCloser, c func() error) *errorPipe {
	return &errorPipe{r: r, c: c}
}

func (r *errorPipe) Read(p []byte) (int, error) {
	if r.e != nil {
		return 0, r.e
	}
	n, err := r.r.Read(p)

	if err == io.EOF {
		// check if we received error after waiting for Wait()
		r.o.Do(func() {
			r.e = r.c()
		})
		if r.e != nil {
			return n, r.e
		}
	}
	return n, err
}

func (r *errorPipe) Close() error {
	return r.r.Close()
}
