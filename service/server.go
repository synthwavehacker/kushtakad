package service

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/kushtaka/kushtakas/listener"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("sensors")

func Run() {
	h := &Hub{}
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, os.Interrupt)
		signal.Notify(s, syscall.SIGTERM)

		select {
		case <-s:
			cancel()
		}
	}()

	incoming := make(chan net.Conn)

	l, err := listener.NewSocket()
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				panic(err)
			}

			incoming <- conn

			// in case of goroutine starvation
			// with many connection and single procs
			runtime.Gosched()
		}
	}()

	err = l.Start(ctx)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case conn := <-incoming:
			go h.handle(conn)
		}
	}
}

func (h *Hub) handle(c net.Conn) {
	log.Debug("handle()")

	sm, newConn, err := h.findService(c)
	if sm == nil {
		log.Debug("No suitable handler for %s => %s: %s", c.RemoteAddr(), c.LocalAddr(), err.Error())
		return
	}

	log.Debug("Handling connection for %s => %s %s(%s)", c.RemoteAddr(), c.LocalAddr(), sm.Name, sm.Type)

	newConn = TimeoutConn(newConn, time.Second*30)

	ctx := context.Background()
	if err := sm.Service.Handle(ctx, newConn); err != nil {
		log.Errorf(color.RedString("Error handling service: %s: %s", sm.Name, err.Error()))
	}
}

type Hub struct {
	mu *sync.Mutex

	// Maps a port and a protocol to an array of pointers to services
	ports map[net.Addr][]*ServiceMap
}

// Wraps a Servicer, adding some metadata
type ServiceMap struct {
	Service Servicer

	Name string
	Type string
}

type Servicer interface {
	Handle(context.Context, net.Conn) error
}

type Service struct {
	Port int
}

type Listener interface {
	Start(ctx context.Context) error
	Accept() (net.Conn, error)
}

// Addr, proto, port, error
func ToAddr(input string) (net.Addr, string, int, error) {
	parts := strings.Split(input, "/")

	if len(parts) != 2 {
		return nil, "", 0, fmt.Errorf(`wrong format (needs to be "protocol/(host:)port")`)
	}

	proto := parts[0]

	host, port, err := net.SplitHostPort(parts[1])
	if err != nil {
		port = parts[1]
	}

	portUint16, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return nil, "", 0, fmt.Errorf("error parsing port value: %s", err.Error())
	}

	switch proto {
	case "tcp":
		addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, port))
		return addr, proto, int(portUint16), err
	case "udp":
		addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(host, port))
		return addr, proto, int(portUint16), err
	default:
		return nil, "", 0, fmt.Errorf("unknown protocol %s", proto)
	}
}

type CanHandlerer interface {
	CanHandle([]byte) bool
}

func (hc *Hub) findService(conn net.Conn) (*ServiceMap, net.Conn, error) {
	localAddr := conn.LocalAddr()

	var serviceCandidates []*ServiceMap

	for k, sc := range hc.ports {
		if !compareAddr(k, localAddr) {
			continue
		}

		serviceCandidates = sc
	}

	if len(serviceCandidates) == 0 {
		return nil, nil, fmt.Errorf("No service configured for the given port")
	} else if len(serviceCandidates) == 1 {
		return serviceCandidates[0], conn, nil
	}

	peekUninitialized := true
	var tConn net.Conn
	var pConn *peekConnection
	var n int
	buffer := make([]byte, 1024)
	for _, service := range serviceCandidates {
		ch, ok := service.Service.(CanHandlerer)
		if !ok {
			return service, conn, nil
		}
		if peekUninitialized {
			tConn = TimeoutConn(conn, time.Second*30)
			pConn = PeekConnection(tConn)
			log.Debug("Peeking connection %s => %s", conn.RemoteAddr(), conn.LocalAddr())
			_, err := pConn.Peek(buffer)
			if err != nil {
				return nil, nil, fmt.Errorf("could not peek bytes: %s", err.Error())
			}
			peekUninitialized = false
		}
		if ch.CanHandle(buffer[:n]) {
			return service, pConn, nil
		}
	}
	return nil, nil, fmt.Errorf("No suitable service for the given port")
}

func (h *Hub) heartbeat() {
	beat := time.Tick(30 * time.Second)
	count := 0
	for range beat {
		count++
	}
}

func compareAddr(addr1 net.Addr, addr2 net.Addr) bool {
	if ta1, ok := addr1.(*net.TCPAddr); ok {
		ta2, ok := addr2.(*net.TCPAddr)
		if !ok {
			return false
		}

		if ta1.Port != ta2.Port {
			return false
		}

		if ta1.IP == nil {
		} else if ta2.IP == nil {
		} else if !ta1.IP.Equal(ta2.IP) {
			return false
		}

		return true
	} else if ua1, ok := addr1.(*net.UDPAddr); ok {
		ua2, ok := addr2.(*net.UDPAddr)
		if !ok {
			return false
		}

		if ua1.Port != ua2.Port {
			return false
		}

		if ua1.IP == nil {
		} else if ua2.IP == nil {
		} else if !ua1.IP.Equal(ua2.IP) {
			return false
		}

		return true
	}

	return false
}
