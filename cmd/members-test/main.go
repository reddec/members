package main

import (
	"flag"
	"fmt"
	"members"
	"os"
	"time"
)

func main() {
	name, _ := os.Hostname()
	var (
		id       = flag.String("id", name, "unique node id")
		interval = flag.Duration("interval", 1*time.Second, "interval to print")
	)
	flag.Parse()
	node, done := members.Node().ID(*id).Start()
	c := time.Tick(*interval)
	for {
		for _, m := range node.Members() {
			fmt.Println("NODE", *id, "knows", m.Info.ID, "at", m.Addr)
		}
		fmt.Println("-------------")
		select {
		case <-c:
		case err := <-done:
			if err != nil {
				fmt.Println("NODE", *id, "finished with error:", err)
			} else {
				fmt.Println("NODE", *id, "finished without error")
			}
		}
	}
}
