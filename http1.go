package go_http_shutdown

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
)

func SpawnHTTP1Server(msg string) *Server {
	address := "localhost:2319"
	s := Server{}

	var err error
	s.listener, err = net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, msg)
	})

	s.srv = &http.Server{
		Addr:    address,
		Handler: mux,
	}

	s.wg.Add(1)
	go func() {
		fmt.Printf("HTTP Listening on %s ...\n", address)
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

func RunHTTP1Client(msg string) func() {
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:     false,
			MaxIdleConns:          0,
			MaxIdleConnsPerHost:   0,
			MaxConnsPerHost:       0,
			IdleConnTimeout:       0,
			ResponseHeaderTimeout: 0,
		},
	}
	return runClient(client, "http://", msg)
}
