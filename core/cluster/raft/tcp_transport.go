package raft

import (
	"errors"
	"io"
	"net"
	"time"
)

var (
	errNotAdvertisable = errors.New("local bind address is not advertisable")
	errNotTCP          = errors.New("local address is not a TCP address")
)

// TCPStreamLayer implements StreamLayer interface for plain TCP.
type TCPStreamLayer struct {
	advertise net.Addr
	listener  *net.TCPListener
}

// NewTCPTransport returns a NetworkTransport that is built on top of
// a TCP streaming transport layer.
func NewTCPTransport(
	bindAddr string,
	advertise net.Addr,
	maxPool int,
	timeout time.Duration,
	logOutput io.Writer, listener net.Listener,
) (*NetworkTransport, error) {
	return newTCPTransport(bindAddr, advertise, maxPool, timeout, func(stream StreamLayer) *NetworkTransport {
		return NewNetworkTransport(stream, maxPool, timeout, logOutput)
	}, listener)
}

// NewTCPTransportWithLogger returns a NetworkTransport that is built on top of
// a TCP streaming transport layer, with log output going to the supplied Logger
func NewTCPTransportWithLogger(
	bindAddr string,
	advertise net.Addr,
	maxPool int,
	timeout time.Duration,
	listener net.Listener,
) (*NetworkTransport, error) {
	return newTCPTransport(bindAddr, advertise, maxPool, timeout, func(stream StreamLayer) *NetworkTransport {
		return NewNetworkTransportWithLogger(stream, maxPool, timeout)
	}, listener)
}

func newTCPTransport(bindAddr string,
	advertise net.Addr,
	maxPool int,
	timeout time.Duration,
	transportCreator func(stream StreamLayer) *NetworkTransport, listener net.Listener) (*NetworkTransport, error) {
	// Try to bind
	//list, err := net.Listen("tcp", bindAddr)
	//if err != nil {
	//	return nil, err
	//}

	// Create stream
	stream := &TCPStreamLayer{
		advertise: advertise,
		listener:  listener.(*net.TCPListener),
	}

	// Verify that we have a usable advertise address
	//addr, ok :=
	//	stream.Addr().(*net.TCPAddr)
	//if !ok {
	//	list.Close()
	//	return nil, errNotTCP
	//}
	//if addr.IP.IsUnspecified() {
	//	list.Close()
	//	return nil, errNotAdvertisable
	//}

	// Create the network transport
	trans := transportCreator(stream)
	return trans, nil
}

// Dial implements the StreamLayer interface.
func (t *TCPStreamLayer) Dial(address string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout("tcp", address, timeout)
}

// Accept implements the net.Listener interface.
func (t *TCPStreamLayer) Accept() (c net.Conn, err error) {
	return t.listener.Accept()
}

// Close implements the net.Listener interface.
func (t *TCPStreamLayer) Close() (err error) {
	return nil
	//return t.listener.Close()
}

// Addr implements the net.Listener interface.
func (t *TCPStreamLayer) Addr() net.Addr {
	// Use an advertise addr if provided
	if t.advertise != nil {
		return t.advertise
	}
	return t.listener.Addr()
}
