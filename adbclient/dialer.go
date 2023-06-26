package adbclient

import (
	"io"
	"net"
	"runtime"

	"github.com/daimall/goadb/wire"
)

type tcpDialer struct {
	NetConn net.Conn
}

// Dial connects to the adb server on the host and port set on the netDialer.
// The zero-value will connect to the default, localhost:5037.
func (td *tcpDialer) Dial(address string) (*wire.Conn, error) {
	var err error
	td.NetConn, err = net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	// net.Conn can't be closed more than once, but wire.Conn will try to close both sender and scanner
	// so we need to wrap it to make it safe.
	safeConn := wire.MultiCloseable(td.NetConn)

	// Prevent leaking the network connection, not sure if TCPConn does this itself.
	// Note that the network connection may still be in use after the conn isn't (scanners/senders
	// can give their underlying connections to other scanner/sender types), so we can't
	// set the finalizer on conn.
	runtime.SetFinalizer(safeConn, func(conn io.ReadWriteCloser) {
		conn.Close()
	})

	return &wire.Conn{
		Scanner: wire.NewScanner(safeConn),
		Sender:  wire.NewSender(safeConn),
	}, nil
}
