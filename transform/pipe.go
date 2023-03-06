package transform

import "io"

// WritePartialCloser is like io.WriterCloser, but can only be close from the write side.
type WritePartialCloser interface {
	Write(p []byte) (n int, err error)
	CloseWrite() error
}

// ClosablePipe is a wrapper of io.PipeWriter, making it compatible with CloseWrite.
type ClosablePipe struct {
	pw *io.PipeWriter
}

// Write implements WritePartialCloser
func (c *ClosablePipe) Write(p []byte) (n int, err error) {
	return c.pw.Write(p)
}

// CloseWrite implements WritePartialCloser
func (c *ClosablePipe) CloseWrite() error {
	return c.pw.Close()
}

// NewClosablePipe wraps a io.PipeWriter in a ClosablePipe.
func NewClosablePipe(pw *io.PipeWriter) *ClosablePipe {
	return &ClosablePipe{pw: pw}
}
