# Distributed multicast-based discovery
[![GoDoc](https://godoc.org/github.com/reddec/members?status.svg)](https://godoc.org/github.com/reddec/members)

`go get -v github.com/reddec/memebers`

* protocol independent
* supports GRPC


![diag](https://user-images.githubusercontent.com/6597086/56138741-0e551400-5fca-11e9-8f2d-40e7bea49d94.png)

## Use-cases

### Automatic service definition

1. initialize node
2. make simple TCP acceptor on a random port
3. add it to local registry
4. start broadcasting and listening

```go
node := New().Start()
node.Simple("echo", func(ctx context.Context, conn net.Conn){
	defer conn.Close()
    io.Copy(conn, conn)	
})
// do other stuff
```


### GRPC server


1. initialize node
2. create listener and add to registry
3. serve grpc server

```go
node := New().Start()

listener := node.Listen("myService")
defer listener.Close()

grpcServer.Serve(listener)	
```


### Native TCP client


```go
node := New().Start()

conn, err := node.Dial("echo") // iterate over all services in loop until connected
// do some stuff
```


### GRPC client

```go
node := New().Start()

grpcConn, err := grpc.Dial("echo", node.WithGRPC())
```

## Protocol

Almost all parameters can be changed during node initialization, however, default parameters are:

* Multicast group IP: `224.0.0.145`
* Port: `33214`
* Buffer size: `8192` bytes
* TTL: `5s`

Packet data is message pack tuple:

`[ID: string, [[ServiceName: string, Port: uint16 ] ,...]]`