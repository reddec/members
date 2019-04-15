package members

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

const (
	Scheme = "mqst"
)

type grpcWatcher struct {
	w *Watcher
}

func (gw *grpcWatcher) ResolveNow(resolver.ResolveNowOption) {
	gw.w.ForceUpdate()
}

func (gw *grpcWatcher) Close() {
	gw.w.Close()
}

type grpcBuilder struct {
	scheme  string
	manager *Manager
}

func (gb *grpcBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	watcher := &grpcWatcher{
		w: gb.manager.Watch(target.Endpoint),
	}
	go func() {
		for upd := range watcher.w.Watch() {
			var addresses = make([]resolver.Address, len(upd))
			for i, addr := range upd {
				addresses[i] = resolver.Address{
					Addr:       addr,
					Type:       resolver.Backend,
					ServerName: addr,
				}
			}
			cc.UpdateState(resolver.State{Addresses: addresses})
		}
	}()
	return watcher, nil
}

func (gb *grpcBuilder) Scheme() string {
	return gb.scheme
}

func (mg *Manager) GRPCResolverSchema(scheme string) {
	resolver.Register(&grpcBuilder{
		scheme:  scheme,
		manager: mg,
	})
}

func (mg *Manager) GRPCResolver() {
	mg.GRPCResolverSchema(Scheme)
}

func (mg *Manager) WithGRPC() grpc.DialOption {
	return grpc.WithContextDialer(mg.DialContext)
}
