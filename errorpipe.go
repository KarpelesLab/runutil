package runutil

import (
	"io"
	"sync"
)

type errorPipe struct {
	r io.ReadCloser
	e error
	s bool // is error set
	l sync.RWMutex
	c *sync.Cond
}

func newErrorPipe(r io.ReadCloser) *errorPipe {
	re := &errorPipe{r: r}
	re.c = sync.NewCond(re.l.RLocker())
	return re
}

func (r *errorPipe) Read(p []byte) (int, error) {
	if e := r.getError(); e != nil {
		return 0, e
	}
	n, err := r.r.Read(p)

	if err == io.EOF {
		// check if we received error after waiting for Wait()
		if e := r.getErrorWait(); e != nil {
			return n, e
		}
	}
	return n, err
}

func (r *errorPipe) Close() error {
	return r.r.Close()
}

func (r *errorPipe) getError() error {
	r.l.RLock()
	defer r.l.RUnlock()
	return r.e
}

func (r *errorPipe) getErrorWait() error {
	// return error, wait for setError() to have been called if it hasn't yet
	r.l.RLock()
	defer r.l.RUnlock()

	for {
		if r.s {
			return r.e
		}

		r.c.Wait()
	}
}

func (r *errorPipe) setError(e error) {
	r.l.Lock()
	defer r.l.Unlock()
	r.e = e
	r.s = true

	r.c.Broadcast()
}
