package relay

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/cortexiphan/tcprelay/transform"

	"golang.org/x/sync/errgroup"
)

type Relay struct {
	localHostPort  string
	remoteHostPort string
	localAddr      *net.TCPAddr
	remoteAddr     *net.TCPAddr
	fromClientMw   transform.ReaderMiddleware
	fromServerMw   transform.ReaderMiddleware
	toServerMw     transform.WriterMiddleware
	toClientMw     transform.WriterMiddleware
}

func NewRelay(opts ...RelayOption) *Relay {
	Relay := &Relay{
		localHostPort:  "localhost:3916",
		remoteHostPort: "localhost:3917",
		fromClientMw:   transform.DummyReaderMiddleware,
		fromServerMw:   transform.DummyReaderMiddleware,
		toClientMw:     transform.DummyWriterMiddleware,
		toServerMw:     transform.DummyWriterMiddleware,
	}
	for _, opt := range opts {
		opt(Relay)
	}
	return Relay
}

func (p *Relay) Run() {
	var err error
	if p.remoteAddr, err = net.ResolveTCPAddr("tcp4", p.remoteHostPort); err != nil {
		panic(fmt.Sprintf("resolve remote addr fail:%v", err))
	}
	if p.localAddr, err = net.ResolveTCPAddr("tcp4", p.localHostPort); err != nil {
		panic(fmt.Sprintf("resolve local addr fail:%v", err))
	}
	l, err := net.ListenTCP("tcp4", p.localAddr)
	if err != nil {
		panic(fmt.Sprintf("listen tcp fail:%v", err))
	}
	for {
		cliConn, err := l.AcceptTCP()
		if err != nil {
			continue
		}
		srvConn, err := net.DialTCP("tcp4", nil, p.remoteAddr)
		if err != nil {
			_ = cliConn.Close()
			continue
		}
		go p.handleConnection(cliConn, srvConn)
	}
}

func (p *Relay) handleConnection(cliConn, srvConn *net.TCPConn) {
	once := &sync.Once{}
	cleanup := func() {
		_ = cliConn.Close()
		_ = srvConn.Close()
	}

	g := &errgroup.Group{}
	g.Go(func() error {
		dst := p.toServerMw(srvConn)
		src := p.fromClientMw(cliConn)
		if _, err := io.Copy(dst, src); err != nil {
			once.Do(cleanup)
			return err
		}
		// EOF
		_ = dst.CloseWrite()
		fmt.Printf("close write client to server\n")
		return nil
	})
	g.Go(func() error {
		dst := p.toClientMw(cliConn)
		src := p.fromServerMw(srvConn)
		if _, err := io.Copy(dst, src); err != nil {
			once.Do(cleanup)
			return err
		}
		// EOF
		_ = dst.CloseWrite()
		fmt.Printf("close write server to client\n")
		return nil
	})
	if err := g.Wait(); err != nil {
		fmt.Printf("forward fail:%v\n", err)
	}
}
