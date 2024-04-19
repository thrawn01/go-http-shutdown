package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	shutdown "github.com/thrawn01/go-http-shutdown"
)

func main() {
	var (
		srv    *shutdown.Server
		cancel func()
	)
	// --------------------------------------
	//  HTTP 1
	// --------------------------------------

	fmt.Printf("\n--- Start HTTP 1 Server ---\n")
	srv = shutdown.SpawnHTTP1Server("Hello, HTTP1")
	cancel = shutdown.RunHTTP1Client("Hello, HTTP1")

	time.Sleep(time.Millisecond * 300)

	fmt.Printf("\n--- HTTP Server Shutdown ---\n")
	fmt.Printf("Should see a client error here....\n")
	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}

	// Allow the client to continue making requests
	time.Sleep(time.Second * 2)

	// Cancel the client
	fmt.Printf("\n--- Shutdown Client ---\n")
	cancel()
	fmt.Printf("\n--- Done ---\n")

	// --------------------------------------
	//  HTTP 2
	// --------------------------------------

	var tls shutdown.TLSConfig
	if err := shutdown.SetupTLS(&tls); err != nil {
		panic(err)
	}

	fmt.Printf("\n--- Start HTTP 2 Server ---\n")
	srv = shutdown.SpawnHTTP2Server(tls, "Hello, HTTP2")
	cancel = shutdown.RunHTTP2Client(tls, "Hello, HTTP2")

	time.Sleep(time.Millisecond * 300)

	fmt.Printf("\n--- HTTP Server Shutdown ---\n")
	fmt.Printf("Should see a client error here....\n")
	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}

	// Allow the client to continue making requests
	time.Sleep(time.Second * 2)

	// Cancel the client
	fmt.Printf("\n--- Shutdown Client ---\n")
	cancel()
	fmt.Printf("\n--- Done ---\n")

	// --------------------------------------
	//  H2C
	// --------------------------------------

	fmt.Printf("\n--- Start H2C Server ---\n")
	srv = shutdown.SpawnH2CServer("Hello, H2C")
	cancel = shutdown.RunH2CClient("Hello, H2C")

	time.Sleep(time.Millisecond * 300)

	fmt.Printf("\n--- HTTP Server Shutdown ---\n")
	fmt.Printf("No error here, client just keeps making requests\n")
	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}

	atomic.StoreInt64(&shutdown.StopTheWorld, 1)

	// Allow the client to continue making requests
	time.Sleep(time.Second * 2)

	// Cancel the client
	fmt.Printf("\n--- Shutdown Client ---\n")
	cancel()
	fmt.Printf("\n--- Done ---\n")

}
