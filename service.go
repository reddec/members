package members

import (
	"context"
	"net"
)

//go:generate msgp
//msgp:tuple Service Info

type Service struct {
	Name string
	Port uint16
}

type Info struct {
	ID       string
	Services []Service
}

func (n *NodeBuilder) Service(name string, port int) *NodeBuilder {
	n.info.Services = append(n.info.Services, Service{Name: name, Port: uint16(port)})
	return n
}

func (n *NodeBuilder) Listen(serviceName string) net.Listener {
	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		panic(err)
	}
	n.Service(serviceName, listener.Addr().(*net.TCPAddr).Port)
	return listener
}

func (n *NodeBuilder) Simple(name string, callback func(ctx context.Context, conn net.Conn)) *NodeBuilder {
	listener := n.Listen(name)
	go func() {
		<-n.ctx.Done()
		listener.Close()
	}()
	go func() {
		defer listener.Close()

		for {
			conn, err := listener.Accept()
			if err != nil {
				break
			}
			go callback(n.ctx, conn)
		}

	}()
	return n
}
