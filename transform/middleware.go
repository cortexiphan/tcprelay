package transform

import (
	"io"
)

// DummyReaderMiddleware returns the input io.Reader directly.
func DummyReaderMiddleware(r io.Reader) io.Reader {
	return r
}

// DummyWriterMiddleware returns the input io.Writer directly.
func DummyWriterMiddleware(w WritePartialCloser) WritePartialCloser {
	return w
}

func readFullBuffer(r io.Reader, buf []byte) ([]byte, error) {
	var (
		i   int
		n   int
		err error
	)
	for {
		n, err = r.Read(buf[i:])
		i += n
		if i == len(buf) {
			return buf, nil
		}
		if err != nil {
			return buf[:i], err
		}
	}
}

type obfuscateFunc func([]byte) []byte

// NewReaderMw takes a buffer size and a obfuscating function to be applied on that buffer,
// returning a ReaderMiddleware.
func NewReaderMw(n int, f obfuscateFunc) ReaderMiddleware {
	return func(r io.Reader) io.Reader {
		pr, pw := io.Pipe()
		go func() {
			var err error
			defer func() {
				_ = pw.CloseWithError(err)
			}()
			bufRead := make([]byte, n)
			for err != io.EOF {
				if bufRead, err = readFullBuffer(r, bufRead); err != nil && err != io.EOF {
					// non EOF error encountered, discard any data received.
					return
				}
				if len(bufRead) == 0 {
					continue
				}
				bufWrite := f(bufRead)
				if _, er := pw.Write(bufWrite); er != nil {
					return
				}
			}
		}()
		return pr
	}
}

// NewWriterMw takes a buffer size and a obfuscating function to be applied on that buffer,
// returning a WriterMiddleware.
func NewWriterMw(n int, f obfuscateFunc) WriterMiddleware {
	return func(w WritePartialCloser) WritePartialCloser {
		pr, pw := io.Pipe()
		go func() {
			var err error
			defer func() {
				_ = w.CloseWrite()
			}()
			bufRead := make([]byte, n)
			for err != io.EOF {
				if bufRead, err = readFullBuffer(pr, bufRead); err != nil && err != io.EOF {
					return
				}
				if len(bufRead) == 0 {
					continue
				}
				bufWrite := f(bufRead)
				if _, er := w.Write(bufWrite); er != nil {
					return
				}
			}
		}()
		return NewClosablePipe(pw)
	}
}
