package go_http_shutdown

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net"
	"net/http"
)

func SpawnH2CServer(msg string) *Server {
	address := "localhost:2319"
	s := Server{}

	var err error
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
		Addr:    address,
		Handler: h2c.NewHandler(mux, h2s),
	}

	s.wg.Add(1)
	go func() {
		if err := s.srv.Serve(s.listener); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
		fmt.Printf("\n--- Serve() Returned ---\n")
		s.wg.Done()
	}()

	if err := waitForConnect(context.Background(), s.listener.Addr().String(), nil); err != nil {
		panic(err)
	}
	return &s
}

func RunH2CClient(msg string) func() {
	// Create an H2C client (HTTP/2 over Cleartext)
	client := &http.Client{
		Transport: &http2.Transport{
			// So http2.Transport doesn't complain the URL scheme isn't 'https'
			AllowHTTP: true,
			// Pretend we are dialing a TLS endpoint. (Note, we ignore the passed tls.Config)
			DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
				var d net.Dialer
				return d.DialContext(ctx, network, addr)
			},
		},
	}
	return runClient(client, "http://", msg)
}
