package server

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"example.com/common"
)

func Start() {

	fmt.Println("Creating Store")
	// create a `*Store` object
	store := common.NewStore()

	// create a custom RPC server
	fmt.Println("Creating Server")
	server := rpc.NewServer()

	// register `mit` object with `server`
	server.Register(store)

	// create a TCP listener at address : 127.0.0.1:9002
	// https://golang.org/pkg/net/#Listener
	listener, _ := net.Listen("tcp", "127.0.0.1:9002")

	for {
		fmt.Println("Waiting for connections")
		// get connection from the listener when client connects
		conn, _ := listener.Accept() // Accept blocks until next connection is received

		// serve connection in a separate goroutine using JSON codec
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}

}
