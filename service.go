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

func (node *Node) Listen(serviceName string) net.Listener {
	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		panic(err)
	}
	node.Declare(serviceName, uint16(listener.Addr().(*net.TCPAddr).Port))
	return listener
}

func (node *Node) Simple(name string, callback func(ctx context.Context, conn net.Conn)) *Node {
	listener := node.Listen(name)
	go func() {
		<-node.ctx.Done()
		listener.Close()
	}()
	go func() {
		defer listener.Close()

		for {
			conn, err := listener.Accept()
			if err != nil {
				break
			}
			go callback(node.ctx, conn)
		}

	}()
	return node
}
