package lents

import (
	"net"
	"sync"
	"time"
)

const (
	generalLock = iota
	readLock
	writeLock
)

const (
	readBuffer = iota
	writeBuffer
)

type Socket struct {
	udp    *net.UDPConn
	lock   [3]sync.Mutex
	buffer [2]*DataBuffer
}

func (socket *Socket) Write(p []byte) (int, error) {
	socket.lock[writeLock].Lock()
	defer socket.lock[writeLock].Unlock()
	return 0, nil
}

func (socket *Socket) Read(p []byte) (int, error) {
	socket.lock[readLock].Lock()
	defer socket.lock[readLock].Unlock()
	return socket.buffer[readBuffer].Pop(p)
}

// Required to satisfy the net.Conn interface
func (socket *Socket) LocalAddr() net.Addr {
	return nil
}
func (socket *Socket) RemoteAddr() net.Addr {
	return nil
}
func (socket *Socket) SetDeadline(t time.Time) error {
	return nil
}
func (socket *Socket) SetReadDeadline(t time.Time) error {
	return nil
}
func (socket *Socket) SetWriteDeadline(t time.Time) error {
	return nil
}

func (socket *Socket) Close() error {
	return nil
}

func (socket *Socket) Dial(network, address string) (net.Conn, error) {
	return nil, nil
}
