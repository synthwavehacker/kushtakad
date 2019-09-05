package service

import (
	"net"
	"sync"
)

func PeekConnection(conn net.Conn) *peekConnection {
	return &peekConnection{
		conn,
		[]byte{},
		sync.Mutex{},
	}
}

type peekConnection struct {
	net.Conn

	buffer []byte
	m      sync.Mutex
}

func (pc *peekConnection) Peek(p []byte) (int, error) {
	pc.m.Lock()
	defer pc.m.Unlock()

	n, err := pc.Conn.Read(p)

	pc.buffer = append(pc.buffer, p[:n]...)
	return n, err
}

func (pc *peekConnection) Read(p []byte) (n int, err error) {
	pc.m.Lock()
	defer pc.m.Unlock()

	// first serve from peek buffer
	if len(pc.buffer) > 0 {
		bn := copy(p, pc.buffer)
		pc.buffer = pc.buffer[bn:]
		return bn, nil
	}

	return pc.Conn.Read(p)
}
