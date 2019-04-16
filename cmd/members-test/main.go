package main

import (
	"flag"
	"fmt"
	"members"
	"os"
)

func main() {
	name, _ := os.Hostname()
	var (
		id = flag.String("id", name, "unique node id")
	)
	flag.Parse()
	node := members.New().ID(*id).Start()

	node.Watch(func(m *members.Member, change members.MemberChange) {
		switch change {
		case members.MemberAdded:
			fmt.Println("[", *id, "]", "discovered", m.Info.ID, "at", m.Addr)
		case members.MemberRemoved:
			fmt.Println("[", *id, "]", "lost", m.Info.ID, "at", m.Addr)
		}
	})
	<-node.Done()
	err := node.Error()
	if err != nil {
		fmt.Println("NODE", *id, "finished with error:", err)
		os.Exit(1)
	} else {
		fmt.Println("NODE", *id, "finished without error")
	}
}
