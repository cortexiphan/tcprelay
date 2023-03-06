package transform

import "math/rand"

// ReaderFlipBit returns ReaderMiddleware that flips every bits of the buffer.
var ReaderFlipBit = NewReaderMw(8, func(b []byte) []byte {
	for i := 0; i < len(b); i++ {
		b[i] = ^b[i]
	}
	return b
})

// WriterFlipBit returns WriterMiddleware that flips every bits of the buffer.
var WriterFlipBit = NewWriterMw(8, func(b []byte) []byte {
	for i := 0; i < len(b); i++ {
		b[i] = ^b[i]
	}
	return b
})

var mask = []byte{0x3, 0xb, 0x7, 0xa, 0x1, 0x3, 0xf, 0x0}

func xorBits(b []byte) []byte {
	for i := 0; i < len(b) && i < len(mask); i++ {
		b[i] = b[i] ^ mask[i]
	}
	return b
}

// ReaderXorBit returns a ReaderMiddleware that xor every bits of the buffer.
var ReaderXorBit = NewReaderMw(8, xorBits)

// ReaderXorBit returns a WriterMiddleware that xor every bits of the buffer.
var WriterXorBit = NewWriterMw(8, xorBits)

func swithBytes(b []byte) []byte {
	n := len(b) - 1
	for i := 0; i < n; i, n = i+1, n-1 {
		b[i], b[n] = b[n], b[i]
	}
	return b
}

// ReaderSwitchBytes returns a ReaderMiddleware that switch bytes inside the buffer.
var ReaderSwitchBytes = NewReaderMw(7, swithBytes)

// WriterSwitchBytes returns a WriterMiddleware that switch bytes inside the buffer.
var WriterSwitchBytes = NewWriterMw(7, swithBytes)

// ReaderDeleteEvery6Bytes deletes one byte every six bytes.
var ReaderDeleteEvery6Bytes = NewReaderMw(6, func(b []byte) []byte {
	if len(b) == 6 {
		return b[:len(b)-1]
	}
	return b
})

func getInsertFunc(size int) func(b []byte) []byte {
	randomBytes := make([]byte, size)
	for i := 0; i < size; i++ {
		randomBytes[i] = byte(rand.Int() % 127)
	}
	i := 0
	return func(b []byte) []byte {
		if len(b) == 5 {
			i = (i + 1) % size
			b = append(b, randomBytes[i])
		}
		return b
	}
}

// WriterInsertEvery5Bytes insert one random byte every five bytes.
var WriterInsertEvery5Bytes = NewWriterMw(5, getInsertFunc(100))
