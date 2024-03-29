package runutil

import (
	"context"
	"io"
	"os/exec"
	"sync"
	"time"
)

type Pipe interface {
	io.ReadCloser
	CopyTo(io.Writer) (int64, error)
	CloseWait(ctx context.Context) error
}

type processPipe struct {
	r io.ReadCloser
	e error
	o sync.Once
	p *exec.Cmd
}

func newProcessPipe(r io.ReadCloser, p *exec.Cmd) *processPipe {
	return &processPipe{r: r, p: p}
}

func (r *processPipe) Read(p []byte) (int, error) {
	if r.e != nil {
		return 0, r.e
	}
	n, err := r.r.Read(p)

	if err == io.EOF {
		// check if we received error after waiting for Wait()
		r.o.Do(func() {
			r.e = r.p.Wait()
		})
		if r.e != nil {
			return n, r.e
		}
	}
	return n, err
}

func (r *processPipe) Close() error {
	err := r.r.Close()

	// call CloseWait() in background
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	go func() {
		defer cancel()
		r.CloseWait(ctx)
	}()

	return err
}

func (r *processPipe) CopyTo(w io.Writer) (int64, error) {
	if r.e != nil {
		return 0, r.e
	}

	// read whole pipe & write to writer
	n, err := io.Copy(w, r.r)
	if err != nil {
		return n, err
	}

	// we reached eof
	r.o.Do(func() {
		r.e = r.p.Wait()
	})
	if r.e != nil {
		return n, r.e
	}
	return n, nil
}

func (r *processPipe) CloseWait(ctx context.Context) error {
	err := r.r.Close()
	w := make(chan struct{})

	go func() {
		r.o.Do(func() {
			r.e = r.p.Wait()
		})
		w <- struct{}{}
	}()

	select {
	case <-w:
	case <-ctx.Done():
		r.p.Process.Kill()
		// force wait after kill
		<-w
	}

	if r.e != nil {
		return r.e
	}

	return err
}
