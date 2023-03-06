// Package transform handles transformation of tcp streams.
package transform

import "io"

// ReaderMiddleware is a function of data stream manipulation on the read side.
type ReaderMiddleware func(from io.Reader) (to io.Reader)

// WriterMiddleware is a function of data stream manipulation on the write side. It can be closed.
type WriterMiddleware func(from WritePartialCloser) (to WritePartialCloser)

// ChainReader chains together ReaderMiddlewares to form a new data stream.
func ChainReader(rms ...ReaderMiddleware) ReaderMiddleware {
	return func(r io.Reader) io.Reader {
		for i := 0; i < len(rms); i++ {
			r = rms[i](r)
		}
		return r
	}
}

// ChainWriter chains together WriterMiddleware to form a new data stream.
func ChainWriter(wms ...WriterMiddleware) WriterMiddleware {
	return func(w WritePartialCloser) WritePartialCloser {
		for i := len(wms) - 1; i >= 0; i-- {
			w = wms[i](w)
		}
		return w
	}
}

// BuildReader chains a slice of ReaderMiddlewares.
func BuildReader(rms []ReaderMiddleware) ReaderMiddleware {
	return ChainReader(rms...)
}

// BuildWriter chains a slice of WriterMiddleware.
func BuildWriter(wms []WriterMiddleware) WriterMiddleware {
	return ChainWriter(wms...)
}
