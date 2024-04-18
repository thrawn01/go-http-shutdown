package go_http_shutdown

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Server struct {
	listener net.Listener
	wg       sync.WaitGroup
	srv      *http.Server
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	// Blindly close all connections.
	//if err := s.srv.Close(); err != nil {
	//	return err
	//}

	s.wg.Wait()
	return nil
}

func waitForConnect(ctx context.Context, address string, cfg *tls.Config) error {
	if address == "" {
		return fmt.Errorf("WaitForConnect() requires a valid address")
	}

	var errs []string
	for {
		var d proxy.ContextDialer
		if cfg != nil {
			d = &tls.Dialer{Config: cfg}
		} else {
			d = &net.Dialer{}
		}
		conn, err := d.DialContext(ctx, "tcp", address)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		errs = append(errs, err.Error())
		if ctx.Err() != nil {
			errs = append(errs, ctx.Err().Error())
			return errors.New(strings.Join(errs, "\n"))
		}
		time.Sleep(time.Millisecond * 100)
		continue
	}
}

func runClient(client *http.Client, scheme, msg string) func() {
	done := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%slocalhost:2319", scheme), nil)
		if err != nil {
			panic(err)
		}
		defer wg.Done()

		for {
			//fmt.Printf("R")
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Client Err: %s\n", err)
				return
			}
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Client Err: %s\n", err)
				return
			}

			if string(b) != msg {
				panic(fmt.Sprintf("'%s' != '%s'", b, msg))
			}
			fmt.Printf("%d", resp.ProtoMajor)

			select {
			case <-done:
				return
			default:
			}
		}
	}()
	return func() {
		close(done)
		wg.Wait()
	}
}
