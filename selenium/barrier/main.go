package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	listen      string
	gracePeriod time.Duration
	target      string
	queue       = make(chan struct{}, 1)
)

func init() {
	flag.StringVar(&listen, "listen", ":4444", "host and port to listen to")
	flag.DurationVar(&gracePeriod, "grace-period", 300*time.Second, "graceful shutdown period")
	flag.StringVar(&target, "target", "http://localhost:4723", "target proxy url")
	flag.Parse()
}

func signalContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		stop := make(chan os.Signal)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
	}()
	return ctx, cancel
}

func main() {
	ctx, cancel := signalContext()
	defer cancel()
	err := mainCtx(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func mux(ctx context.Context, u *url.URL) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/wd/hub/session", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			return
		case <-ctx.Done():
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		case queue <- struct{}{}:
		}
		(&httputil.ReverseProxy{
			Director: func(r *http.Request) {
				r.URL.Scheme, r.URL.Host = u.Scheme, u.Host
			},
			ModifyResponse: func(resp *http.Response) error {
				if resp.StatusCode != http.StatusOK {
					<-queue
				}
				return nil
			},
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
				<-queue
				log.Printf("proxy to %s: %v", u.String(), err)
				http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)

			},
		}).ServeHTTP(w, r)
	}))
	mux.Handle("/wd/hub/session/", &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			fragments := strings.Split(r.URL.Path, "/")
			if r.Method == http.MethodDelete && len(fragments) == 3 {
				<-queue
			}
			r.URL.Scheme, r.URL.Host = u.Scheme, u.Host
		}})
	return mux
}

func mainCtx(ctx context.Context) error {
	target, err := url.Parse(target)
	if err != nil {
		return fmt.Errorf("parse target url: %v", err)
	}
	server := &http.Server{
		Addr:    listen,
		Handler: mux(ctx, target),
	}
	e := make(chan error)
	go func() {
		e <- server.ListenAndServe()
	}()
	select {
	case err := <-e:
		return err
	case <-ctx.Done():
	}
	log.Printf("starting graceful shutdown in %v]", gracePeriod)
	shCtx, cancel := context.WithTimeout(context.Background(), gracePeriod)
	defer cancel()
	if err := server.Shutdown(shCtx); err != nil {
		return fmt.Errorf("graceful shutdown: %v", err)
	}
	log.Printf("stopped")
	return nil
}
