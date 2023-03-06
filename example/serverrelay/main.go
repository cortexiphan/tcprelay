package main

import (
	"github.com/cortexiphan/tcprelay/relay"
	mw "github.com/cortexiphan/tcprelay/transform"
)

var (
	localHostPort  = "localhost:8081"
	remoteHostPort = "localhost:8080"
	fromClientMw   = mw.ChainReader(mw.ReaderXorBit, mw.ReaderSwitchBytes, mw.ReaderDeleteEvery6Bytes)
	fromServerMw   = mw.ChainReader(mw.DummyReaderMiddleware)
	toServerMw     = mw.ChainWriter(mw.DummyWriterMiddleware)
	toClientMw     = mw.ChainWriter(mw.WriterInsertEvery5Bytes, mw.WriterSwitchBytes, mw.WriterXorBit)
)

func main() {
	serverRelay := relay.NewRelay(
		relay.WithLocalAddr(localHostPort),
		relay.WithRemoteAddr(remoteHostPort),
		relay.WithFromClientMiddleware(fromClientMw),
		relay.WithFromServerMiddleware(fromServerMw),
		relay.WithToClientMiddleware(toClientMw),
		relay.WithToServerMiddleware(toServerMw),
	)
	serverRelay.Run()
}
