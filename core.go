package members

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"os"
	"time"
)

type Member struct {
	Info Info
	Addr net.Addr
	At   time.Time
}

type Logger interface {
	Println(...interface{})
}

type mockLogger struct {
}

func (ml *mockLogger) Println(...interface{}) {}

type Handler interface {
	Discovered(member *Member)
	Lost(member *Member)
	Updated(member *Member)
}

type NodeBuilder struct {
	ctx        context.Context
	ttl        time.Duration
	bufferSize int
	ip         string
	port       int
	info       Info
	logger     Logger
}

func Node() *NodeBuilder {
	return &NodeBuilder{
		info: Info{
			ID: uuid.New().String(),
		},
		ttl:        5 * time.Second,
		bufferSize: 8192,
		ctx:        context.Background(),
		ip:         "224.0.0.145",
		port:       33214,
		logger:     &mockLogger{},
	}
}

func (n *NodeBuilder) Port(num int) *NodeBuilder {
	n.port = num
	return n
}

func (n *NodeBuilder) Logger(logger Logger) *NodeBuilder {
	n.logger = logger
	return n
}

func (n *NodeBuilder) StdLogger(prefix string) *NodeBuilder {
	return n.Logger(log.New(os.Stderr, prefix, log.LstdFlags))
}

func (n *NodeBuilder) Multicast(address string) *NodeBuilder {
	n.ip = address
	return n
}

func (n *NodeBuilder) Buffer(size int) *NodeBuilder {
	n.bufferSize = size
	return n
}

func (n *NodeBuilder) TTL(interval time.Duration) *NodeBuilder {
	n.ttl = interval
	return n
}

func (n *NodeBuilder) Context(ctx context.Context) *NodeBuilder {
	n.ctx = ctx
	return n
}

func (n *NodeBuilder) Start() (*Manager, <-chan error) {
	mgr := &Manager{}
	return mgr, n.StartHandler(mgr)
}

func (n *NodeBuilder) StartHandler(handler Handler) <-chan error {
	done := make(chan error, 1)
	if handler == nil {
		done <- errors.New("handler not defined")
		close(done)
		return done
	}
	dataMsg, err := n.info.MarshalMsg(nil)
	if err != nil {
		done <- err
		close(done)
		return done
	}
	multicastAddr, err := net.ResolveUDPAddr("udp", fmt.Sprint(n.ip, ":", n.port))
	if err != nil {
		n.logger.Println("failed list all available interfaces:", err)
		done <- err
		close(done)
		return done
	}
	interfaces, err := net.Interfaces()
	if err != nil {
		n.logger.Println("failed to list all available interfaces:", err)
		done <- err
		close(done)
		return done
	}
	udpConn, err := net.ListenUDP("udp", multicastAddr)
	if err != nil {
		n.logger.Println("failed to listen on udp address", multicastAddr.String(), ":", err)
		done <- err
		close(done)
		return done
	}

	// join group
	wrapped := ipv4.NewPacketConn(udpConn)
	for _, ifi := range interfaces {
		if ifi.Flags&net.FlagMulticast != 0 {
			n.logger.Println("multicast interface:", ifi.Name)
			err = wrapped.JoinGroup(&ifi, multicastAddr)
			if err != nil {
				n.logger.Println("failed to join multicast group", multicastAddr.String(), "on interface", ifi.Name, ":", err)
				udpConn.Close()
				done <- err
				close(done)
				return done
			}
		}
	}

	go func() {
		n.logger.Println("node loop started")
		done <- nodeLoop(n.ctx, n.ttl, dataMsg, udpConn, n.bufferSize, handler, multicastAddr, n.logger)
		n.logger.Println("node loop finished")
		close(done)
		udpConn.Close()
	}()
	return done
}

func nodeLoop(ctx context.Context, ttl time.Duration, id []byte, conn net.PacketConn, bufferSize int, handler Handler, castAddress net.Addr, logger Logger) error {
	type packet struct {
		Addr net.Addr
		Data []byte
	}

	done := make(chan error, 1)
	data := make(chan packet)

	members := map[string]*Member{}

	go func() {
		defer close(done)
		for {
			buffer := make([]byte, bufferSize) // TODO: make a zero-copy buffer
			n, addr, err := conn.ReadFrom(buffer)
			if err != nil {
				done <- err
				return
			}
			select {
			case data <- packet{Addr: addr, Data: buffer[:n]}:
			case <-ctx.Done():
				return
			}
		}
	}()

	cleanupTicker := time.NewTicker(ttl)
	defer cleanupTicker.Stop()

	keepAliveTicker := time.NewTicker(ttl / 3)
	defer keepAliveTicker.Stop()

	logger.Println("cleanup interval:", ttl)
	logger.Println("keep-alive interval:", ttl/3)
	// send notification about your self
	conn.WriteTo(id, castAddress)
LOOP:
	for {
		select {
		case msg := <-data:
			var serviceInfo Info
			_, err := serviceInfo.UnmarshalMsg(msg.Data)
			if err != nil {
				// just skip
				logger.Println("can't decode message from", msg.Addr, ":", err)
				continue
			}

			_, exists := members[serviceInfo.ID]

			mem := &Member{
				Addr: msg.Addr,
				Info: serviceInfo,
				At:   time.Now(),
			}
			members[serviceInfo.ID] = mem
			if !exists {
				logger.Println("new service discovered:", serviceInfo.ID)
				handler.Discovered(mem)
			} else {
				logger.Println("service", serviceInfo.ID, "keep-alive")
				handler.Updated(mem)
			}
		case <-cleanupTicker.C:
			// clean deprecated members
			now := time.Now()
			var toDel []string
			for key, member := range members {
				if now.Sub(member.At) > ttl {
					handler.Lost(member)
					toDel = append(toDel, key)
				}
			}
			for _, k := range toDel {
				logger.Println("service", k, "removed due TTL")
				delete(members, k)
			}
		case <-keepAliveTicker.C:
			// send notification about your self
			conn.WriteTo(id, castAddress)
		case <-done:
			break LOOP
		case <-ctx.Done():
			break LOOP
		}
	}
	conn.Close()
	return <-done
}
