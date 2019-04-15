package members

import (
	"google.golang.org/grpc"
)

func (mg *Manager) WithGRPC() grpc.DialOption {
	return grpc.WithContextDialer(mg.DialContext)
}
