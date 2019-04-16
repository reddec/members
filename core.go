package members

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"os"
	"sync"
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

type NodeBuilder struct {
	ctx        context.Context
	ttl        time.Duration
	bufferSize int
	ip         string
	port       int
	info       Info
	logger     Logger
}

func New() *NodeBuilder {
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

func (n *NodeBuilder) ID(id string) *NodeBuilder {
	n.info.ID = id
	return n
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

func (n *NodeBuilder) Start() *Node {
	ctx, cancel := context.WithCancel(n.ctx)
	node := &Node{
		info:    n.info,
		update:  make(chan struct{}),
		done:    make(chan struct{}),
		ctx:     ctx,
		members: map[string]*Member{},
		logger:  n.logger,
	}

	multicastAddr, err := net.ResolveUDPAddr("udp", fmt.Sprint(n.ip, ":", n.port))
	if err != nil {
		n.logger.Println("failed list all available interfaces:", err)
		cancel()
		return node.withError(err)
	}
	interfaces, err := net.Interfaces()
	if err != nil {
		n.logger.Println("failed to list all available interfaces:", err)
		cancel()
		return node.withError(err)
	}
	udpConn, err := net.ListenUDP("udp", multicastAddr)
	if err != nil {
		n.logger.Println("failed to listen on udp address", multicastAddr.String(), ":", err)
		cancel()
		return node.withError(err)
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
				cancel()
				return node.withError(err)
			}
		}
	}

	go func() {
		defer udpConn.Close()
		defer cancel()
		n.logger.Println("node loop started")
		node.runLoop(n.ttl, udpConn, n.bufferSize, multicastAddr)
		n.logger.Println("node loop finished")
	}()
	return node
}

type Node struct {
	ctx    context.Context
	logger Logger
	err    error
	done   chan struct{}
	// internal services management
	infoLock sync.RWMutex
	info     Info
	update   chan struct{}
	// members management
	membersLock sync.RWMutex
	members     map[string]*Member
}

func (node *Node) Error() error {
	return node.err
}

func (node *Node) withError(err error) *Node {
	node.err = err
	close(node.done)
	return node
}

func (node *Node) Done() <-chan struct{} {
	return node.done
}

func (node *Node) Notify() {
	select {
	case node.update <- struct{}{}:
	case <-node.ctx.Done():
	}
}

func (node *Node) Declare(serviceName string, port uint16) {
	node.infoLock.Lock()
	updated := false
	for i, srv := range node.info.Services {
		if srv.Name == serviceName {
			node.info.Services[i].Port = port
			updated = true
			break
		}
	}
	if !updated {
		node.info.Services = append(node.info.Services, Service{Port: port, Name: serviceName})
	}
	node.infoLock.Unlock()
}

func (node *Node) Forget(serviceName string) {
	node.infoLock.Lock()
	for i, srv := range node.info.Services {
		if srv.Name == serviceName {
			node.info.Services = append(node.info.Services[:i], node.info.Services[i+1:]...)
			break
		}
	}
	node.infoLock.Unlock()
}

func (node *Node) servicesInfo() []byte {
	node.infoLock.RLock()
	data, err := node.info.MarshalMsg(nil)
	node.infoLock.RUnlock()
	if err != nil {
		panic(err)
	}
	return data
}

type packet struct {
	Addr net.Addr
	Data []byte
}

func (node *Node) runListenerLoop(conn net.PacketConn, bufferSize int, data chan<- packet) error {
	defer close(node.done)
	defer close(data)
	for {
		buffer := make([]byte, bufferSize) // TODO: make a zero-copy buffer
		size, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			return err
		}
		select {
		case data <- packet{Addr: addr, Data: buffer[:size]}:
		case <-node.ctx.Done():
			return node.ctx.Err()
		}
	}
}

func (node *Node) runLoop(ttl time.Duration, conn net.PacketConn, bufferSize int, castAddress net.Addr) error {
	defer conn.Close()

	data := make(chan packet)

	go node.runListenerLoop(conn, bufferSize, data)

	cleanupTicker := time.NewTicker(ttl)
	defer cleanupTicker.Stop()

	keepAliveTicker := time.NewTicker(ttl / 3)
	defer keepAliveTicker.Stop()

	node.logger.Println("cleanup interval:", ttl)
	node.logger.Println("keep-alive interval:", ttl/3)
	// send notification about your self
	conn.WriteTo(node.servicesInfo(), castAddress)
LOOP:
	for {
		select {
		case msg, ok := <-data:
			if !ok {
				break
			}
			var serviceInfo Info
			_, err := serviceInfo.UnmarshalMsg(msg.Data)
			if err != nil {
				// just skip
				node.logger.Println("can't decode message from", msg.Addr, ":", err)
				continue
			}
			mem := &Member{
				Addr: msg.Addr,
				Info: serviceInfo,
				At:   time.Now(),
			}
			node.upsertMember(mem)
		case <-cleanupTicker.C:
			// clean deprecated members
			now := time.Now()
			members := node.Members()
			for _, member := range members {
				if now.Sub(member.At) > ttl {
					node.logger.Println("service", member.Info.ID, "removed due TTL")
					node.removeMember(member)
				}
			}
		case <-keepAliveTicker.C:
			// send notification about your self
			conn.WriteTo(node.servicesInfo(), castAddress)
		case <-node.update:
			// force update
			conn.WriteTo(node.servicesInfo(), castAddress)
		case <-node.ctx.Done():
			break LOOP
		}
	}
	return node.ctx.Err()
}
