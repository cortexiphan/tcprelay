package main

import (
	"github.com/cortexiphan/tcprelay/relay"
	mw "github.com/cortexiphan/tcprelay/transform"
)

var (
	localHostPort  = "localhost:8082"
	remoteHostPort = "localhost:8081"
	fromClientMw   = mw.ChainReader(mw.DummyReaderMiddleware)
	fromServerMw   = mw.ChainReader(mw.ReaderXorBit, mw.ReaderSwitchBytes, mw.ReaderDeleteEvery6Bytes)
	toServerMw     = mw.ChainWriter(mw.WriterInsertEvery5Bytes, mw.WriterSwitchBytes, mw.WriterXorBit)
	toClientMw     = mw.ChainWriter(mw.DummyWriterMiddleware)
)

func main() {
	clientRelay := relay.NewRelay(
		relay.WithLocalAddr(localHostPort),
		relay.WithRemoteAddr(remoteHostPort),
		relay.WithFromClientMiddleware(fromClientMw),
		relay.WithFromServerMiddleware(fromServerMw),
		relay.WithToClientMiddleware(toClientMw),
		relay.WithToServerMiddleware(toServerMw),
	)
	clientRelay.Run()
}
