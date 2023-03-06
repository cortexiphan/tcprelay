package relay

import "github.com/cortexiphan/tcprelay/transform"

type RelayOption func(*Relay)

func WithLocalAddr(hostPort string) RelayOption {
	return func(proxy *Relay) {
		proxy.localHostPort = hostPort
	}
}

func WithRemoteAddr(hostPort string) RelayOption {
	return func(proxy *Relay) {
		proxy.remoteHostPort = hostPort
	}
}

func WithFromClientMiddleware(mw transform.ReaderMiddleware) RelayOption {
	return func(proxy *Relay) {
		proxy.fromClientMw = mw
	}
}

func WithFromServerMiddleware(mw transform.ReaderMiddleware) RelayOption {
	return func(proxy *Relay) {
		proxy.fromServerMw = mw
	}
}

func WithToClientMiddleware(mw transform.WriterMiddleware) RelayOption {
	return func(proxy *Relay) {
		proxy.toClientMw = mw
	}
}

func WithToServerMiddleware(mw transform.WriterMiddleware) RelayOption {
	return func(proxy *Relay) {
		proxy.toServerMw = mw
	}
}
