# Distributed multicast-based discovery

`go get -v github.com/reddec/memebers`

* protocol independent
* supports GRPC

![diag](https://user-images.githubusercontent.com/6597086/56138741-0e551400-5fca-11e9-8f2d-40e7bea49d94.png)


## Automatic service definition

1. initialize node
2. make simple TCP acceptor on a random port
3. add it to local registry
4. start broadcasting and listening

```go
node := Node()
node.Simple("echo", func(ctx context.Context, conn net.Conn){
	defer conn.Close()
    io.Copy(conn, conn)	
})
node.Start()
// do other stuff
```


## GRPC server


1. initialize node
2. create listener and add to registry
3. serve grpc server

```go
node := Node()

listener := node.Listen("myService")
defer listener.Close()

node.Start()

grpcServer.Serve(listener)	
```


## Native TCP client


```go
node := Node()
manager := node.Start()

conn, err := manager.Dial("echo") // iterate over all services in loop until connected
// do some stuff
```


## GRPC client

```go
node := Node()
manager := node.Start()

grpcConn, err := grpc.Dial("echo", manager.WithGRPC())
```