### H2C Server Shutdown Demo
This code demonstrates that an H2C Server does not gracefully shutdown properly. It will continue to
to service requests so long as a client is connected via an HTTP2 stream.

The demo code starts an HTTP1/HTTP2/H2C and client in order in a go routine. The client will constantly make requests 
to the server in order to keep the connection in an active state. After a few milliseconds of client requests
the code preforms a graceful shutdown of the server. In both the HTTP/1 and HTTP/2 case the server eventually
closes all connections and refuses new connections.

However, the H2C server never terminates active connections and continues to serve requests forever, even after the
server has preformed a graceful shutdown.

Running `cmd/demo/main.go` has the following output

```
--- Start HTTP 1 Server ---
HTTP Listening on localhost:2319 ...
111111111111111111111111111111111111111111111111111111111111111111111111111111111111 [SNIP]
--- HTTP Server Shutdown ---
Should see a client error here....
11
--- Serve() Returned ---
Client Err: Get "http://localhost:2319": dial tcp [::1]:2319: connect: connection refused

--- Shutdown Client ---

--- Done ---
2024/04/18 16:53:16 Generating CA Certificates....
2024/04/18 16:53:16 Generating Server Private Key and Certificate....
2024/04/18 16:53:16 Cert DNS names: (localhost)
2024/04/18 16:53:16 Cert IPs: (127.0.0.1)

--- Start HTTP 2 Server ---
HTTP Listening on localhost:2319 ...
22222222222222222222222222222222222222222222222222222222222222222222222222222222222  [SNIP]
--- HTTP Server Shutdown ---
Should see a client error here....
2
--- Serve() Returned ---
2Client Err: Get "https://localhost:2319": dial tcp [::1]:2319: connect: connection refused

--- Shutdown Client ---

--- Done ---

--- Start H2C Server ---
22222222222222222222222222222222222222222222222222222222222222222222222222222222222 [SNIP] 
--- HTTP Server Shutdown ---
No error here, client just keeps making requests

--- Serve() Returned ---
22222222222222222222222222222222222222222222222222222222222222222222222222222222222 [SNIP]
--- Shutdown Client ---
2
--- Done ---
```

### 2024-04-19 UPDATE!
I discovered when creating a h2c server you must call `http2.ConfigureServer()` in order of the graceful shutdown
message to be registered with `http.Server` such that it is called when `Server.Shutdown()` is called.
```go
// -- SNIP --

s.listener, err = net.Listen("tcp", address)
if err != nil {
	panic(err)
}

h2s := &http2.Server{}
mux := http.NewServeMux()
mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, msg)
})

s.srv = &http.Server{
	Handler: h2c.NewHandler(mux, h2s),
	Addr:    address,
}

// Must add ConfigureServer in order for graceful shutdown to work as expected.
if err := http2.ConfigureServer(s.srv, h2s); err != nil {
	panic(err)
}

// -- SNIP --
```
