package members

import (
	"google.golang.org/grpc"
)

func (node *Node) WithGRPC() grpc.DialOption {
	return grpc.WithContextDialer(node.DialContext)
}
