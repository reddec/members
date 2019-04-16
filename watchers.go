package members

import (
	"context"
	"math/rand"
	"net"
	"time"
)

func (node *Node) WaitContextInterval(ctx context.Context, serviceName string, poll time.Duration) (string, error) {
	for {
		addr, ok := node.Random(serviceName)
		if ok {
			return addr, nil
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(poll):
		}
	}
}

func (node *Node) WaitContext(ctx context.Context, serviceName string) (string, error) {
	return node.WaitContextInterval(ctx, serviceName, 1*time.Second)
}

func (node *Node) DialContext(ctx context.Context, serviceName string) (net.Conn, error) {
	return node.DialContextInterval(ctx, serviceName, 5*time.Second)
}

func (node *Node) DialContextInterval(ctx context.Context, serviceName string, interval time.Duration) (net.Conn, error) {
	for {
		addresses := node.Find(serviceName)
		rand.Shuffle(len(addresses), func(i, j int) {
			addresses[i], addresses[j] = addresses[j], addresses[i]
		})
		for _, addr := range addresses {
			dialer := &net.Dialer{}
			conn, err := dialer.DialContext(ctx, "tcp", addr)
			if err == nil {
				return conn, nil
			}
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:

			}
		}
		select {
		case <-time.After(interval):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (node *Node) Dial(serviceName string) (net.Conn, error) {
	return node.DialContext(context.Background(), serviceName)
}
