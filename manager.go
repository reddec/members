package members

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

type Manager struct {
	lock    sync.RWMutex
	members map[string]*Member
}

func (m *Manager) Discovered(member *Member) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.members == nil {
		m.members = make(map[string]*Member)
	}
	m.members[member.Info.ID] = member
}

func (m *Manager) Updated(member *Member) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.members == nil {
		m.members = make(map[string]*Member)
	}
	m.members[member.Info.ID] = member
}

func (m *Manager) Lost(member *Member) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.members, member.Info.ID)
}

func (m *Manager) Members() []*Member {
	ans := make([]*Member, 0, len(m.members))
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, mem := range m.members {
		ans = append(ans, mem)
	}
	return ans
}

func (m *Manager) IP(id string) (string, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	mem, ok := m.members[id]
	if !ok {
		return "", false
	}
	return mem.Addr.(*net.UDPAddr).IP.String(), true
}

func (m *Manager) Find(serviceName string) []string {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.findUnsafe(serviceName)
}

func (m *Manager) findUnsafe(serviceName string) []string {
	var addresses []string
	for _, mem := range m.members {
		for _, srv := range mem.Info.Services {
			if srv.Name != serviceName {
				continue
			}
			addresses = append(addresses, fmt.Sprint(mem.Addr.(*net.UDPAddr).IP.String(), ":", srv.Port))
		}
	}
	return addresses
}

func (m *Manager) First(serviceName string) (string, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, mem := range m.members {
		for _, srv := range mem.Info.Services {
			if srv.Name != serviceName {
				continue
			}
			return fmt.Sprint(mem.Addr.(*net.UDPAddr).IP.String(), ":", srv.Port), true
		}
	}
	return "", false
}

func (m *Manager) Random(serviceName string) (string, bool) {
	addrs := m.Find(serviceName)
	if len(addrs) == 0 {
		return "", false
	}
	idx := rand.Intn(len(addrs))
	return addrs[idx], true
}

func (m *Manager) WaitContextInterval(ctx context.Context, serviceName string, poll time.Duration) (string, error) {
	for {
		addr, ok := m.Random(serviceName)
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

func (m *Manager) WaitContext(ctx context.Context, serviceName string) (string, error) {
	return m.WaitContextInterval(ctx, serviceName, 1*time.Second)
}

func (m *Manager) DialContext(ctx context.Context, serviceName string) (net.Conn, error) {
	return m.DialContextInterval(ctx, serviceName, 5*time.Second)
}

func (m *Manager) DialContextInterval(ctx context.Context, serviceName string, interval time.Duration) (net.Conn, error) {
	for {
		addresses := m.Find(serviceName)
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

func (m *Manager) Dial(serviceName string) (net.Conn, error) {
	return m.DialContext(context.Background(), serviceName)
}
