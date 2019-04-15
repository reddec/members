package members

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"net"
	"testing"
	"time"
)

func TestNode(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	nodeA := Node().Context(ctx).TTL(1 * time.Second)
	nodeB := Node().Context(ctx).TTL(1 * time.Second)

	nodeA.Simple("echo", func(ctx context.Context, conn net.Conn) {
		done := make(chan struct{})
		go func() {
			select {
			case <-done:
			case <-ctx.Done():
			}
			conn.Close()
		}()
		io.Copy(conn, conn)
		close(done)
	})

	alice, cA := nodeA.Start()
	_, cB := nodeB.Start()

	time.Sleep(1 * time.Second)

	dialer := alice.WithGRPC()

	conn, err := grpc.DialContext(ctx, "echo", grpc.WithInsecure(), dialer)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("resolved!", conn.Target())
		defer conn.Close()
	}
	err = <-cA
	fmt.Println("alice done:", err)
	err = <-cB
	fmt.Println("bob done:", err)

	select {
	case <-ctx.Done():
	default:
		t.Error("complete not by timeout")
	}
}
