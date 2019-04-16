package members

import (
	"context"
	"fmt"
	"time"
)

func testCtx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	return ctx
}

func ExampleNew() {
	alice := New().Context(testCtx()).Start()
	<-alice.Done()
}

func ExampleNew_WithID() {
	alice := New().ID("alice").Context(testCtx()).Start()
	<-alice.Done()
}

func ExampleNode_Declare() {
	alice := New().ID("alice").Context(testCtx()).Start()
	bob := New().ID("bob").Context(testCtx()).Start()

	alice.Declare("nginx", 80)
	alice.Declare("wordpress", 8000)
	alice.Notify() // force notify to not wait for update tick

	time.Sleep(100 * time.Millisecond) // let bob settle data

	for _, member := range bob.Members() {
		if member.Info.ID == "alice" {
			for _, serviceInfo := range member.Info.Services {
				fmt.Println(serviceInfo.Name, "->", serviceInfo.Port)
			}
			break
		}
	}
	// output:
	// nginx -> 80
	// wordpress -> 8000
}
