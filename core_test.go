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

	nodeA := New().ID("alice").Context(ctx).TTL(1 * time.Second).Start()
	nodeB := New().ID("bob").Context(ctx).TTL(1 * time.Second).Start()

	if nodeA.Error() != nil {
		t.Error(nodeA.Error())
		return
	}

	if nodeB.Error() != nil {
		t.Error(nodeB.Error())
		return
	}

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

	time.Sleep(1 * time.Second)

	dialer := nodeA.WithGRPC()

	conn, err := grpc.DialContext(ctx, "echo", grpc.WithInsecure(), dialer)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("resolved!", conn.Target())
		defer conn.Close()
	}
	<-nodeA.Done()
	fmt.Println("alice done")
	<-nodeB.Done()
	fmt.Println("bob done")

	select {
	case <-ctx.Done():
	default:
		t.Error("complete not by timeout")
	}
}
