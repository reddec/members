package members

import (
	"fmt"
	"math/rand"
	"net"
)

func (node *Node) Members() []*Member {
	var nodes = make([]*Member, 0, len(node.members))
	node.membersLock.RLock()
	defer node.membersLock.RUnlock()
	for _, m := range node.members {
		nodes = append(nodes, m)
	}
	return nodes
}

func (node *Node) upsertMember(member *Member) {
	node.membersLock.Lock()
	defer node.membersLock.Unlock()
	if node.members == nil {
		node.members = make(map[string]*Member)
	}
	node.members[member.Info.ID] = member
}

func (node *Node) removeMember(member *Member) {
	node.membersLock.Lock()
	defer node.membersLock.Unlock()
	delete(node.members, member.Info.ID)
}

func (node *Node) IP(id string) (string, bool) {
	node.membersLock.RLock()
	defer node.membersLock.RUnlock()
	mem, ok := node.members[id]
	if !ok {
		return "", false
	}
	return mem.Addr.(*net.UDPAddr).IP.String(), true
}

func (node *Node) Find(serviceName string) []string {
	node.membersLock.RLock()
	defer node.membersLock.RUnlock()
	return node.findUnsafe(serviceName)
}

func (node *Node) findUnsafe(serviceName string) []string {
	var addresses []string
	for _, mem := range node.members {
		for _, srv := range mem.Info.Services {
			if srv.Name != serviceName {
				continue
			}
			addresses = append(addresses, fmt.Sprint(mem.Addr.(*net.UDPAddr).IP.String(), ":", srv.Port))
		}
	}
	return addresses
}

func (node *Node) First(serviceName string) (string, bool) {
	node.membersLock.RLock()
	defer node.membersLock.RUnlock()
	for _, mem := range node.members {
		for _, srv := range mem.Info.Services {
			if srv.Name != serviceName {
				continue
			}
			return fmt.Sprint(mem.Addr.(*net.UDPAddr).IP.String(), ":", srv.Port), true
		}
	}
	return "", false
}

func (node *Node) Random(serviceName string) (string, bool) {
	addrs := node.Find(serviceName)
	if len(addrs) == 0 {
		return "", false
	}
	idx := rand.Intn(len(addrs))
	return addrs[idx], true
}
