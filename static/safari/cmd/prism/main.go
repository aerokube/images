package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

var (
	listen         = ":4444"
	target         = "http://localhost:5555"
	waitTimeout    = 30 * time.Second
	gracePeriod    = 30 * time.Second
	browserName    = "safari"
	browserVersion = "15.0"
)

func wait(ctx context.Context, target string) (*url.URL, error) {
	for {
		r, err := http.NewRequest(http.MethodHead, target, http.NoBody)
		if err != nil {
			return nil, fmt.Errorf("new %s request to %s: %v", http.MethodHead, target, err)
		}
		resp, err := http.DefaultClient.Do(r.WithContext(ctx))
		if resp != nil {
			resp.Body.Close()
		}
		if err != nil {
			if err, ok := err.(*url.Error); ok {
				switch err.Err {
				case context.Canceled, context.DeadlineExceeded:
					return nil, err
				default:
					<-time.After(100 * time.Millisecond)
					continue
				}
			}
		}
		return r.URL, nil
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	e := make(chan error)
	go func() {
		stop := make(chan os.Signal)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		select {
		case err := <-e:
			log.Fatalf("server: %v", err)
		case <-stop:
			cancel()
		}
	}()
	waitCtx, waitCancel := context.WithTimeout(ctx, waitTimeout)
	defer waitCancel()
	u, err := wait(waitCtx, target)
	if err != nil {
		log.Fatal(fmt.Errorf("wait target: %v", err))
	}
	server := &http.Server{
		Addr: listen,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var value map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&value)
			if err != nil && err != io.EOF {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err == nil {
				if _, ok := value["desiredCapabilities"]; ok {
					delete(value, "desiredCapabilities")
				}
				if o, ok := value["capabilities"]; ok {
					if w3cCapabilities, ok := o.(map[string]interface{}); ok {
						for _, match := range []string{"alwaysMatch", "firstMatch"} {
							delete(w3cCapabilities, match)
						}
					}
				}
				body, err := json.Marshal(value)
				if err != nil {
					log.Printf("[ERROR] marshalling capabilities: %v", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				r.Body = io.NopCloser(bytes.NewReader(body))
				r.ContentLength = int64(len(body))
			}
			(&httputil.ReverseProxy{
				Director: func(r *http.Request) {
					r.URL.Scheme, r.URL.Host, r.URL.Path = u.Scheme, u.Host, path.Join(u.Path, r.URL.Path)
				},
				ModifyResponse: func(resp *http.Response) error {
					if resp.StatusCode != http.StatusOK {
						return nil
					}
					var values map[string]interface{}
					defer resp.Body.Close()
					err := json.NewDecoder(resp.Body).Decode(&values)
					if err != nil {
						return fmt.Errorf("decode json response: %v", err)
					}
					if o, ok := values["value"]; ok {
						if value, ok := o.(map[string]interface{}); ok {
							if o, ok := value["capabilities"]; ok {
								if capabilities, ok := o.(map[string]interface{}); ok {
									capabilities["browserName"] = browserName
									capabilities["browserVersion"] = browserVersion
									delete(capabilities, "platformName")
								}
							}
						}
					}
					buf, err := json.Marshal(&values)
					if err != nil {
						return fmt.Errorf("encode json response: %v", err)
					}
					resp.Header.Del("Server")
					resp.Header.Del("Content-Length")
					resp.ContentLength = int64(len(buf))
					resp.Body = io.NopCloser(bytes.NewReader(buf))
					return nil
				},
			}).ServeHTTP(w, r)
		}),
	}
	go func() {
		e <- server.ListenAndServe()
	}()
	<-ctx.Done()
	shCtx, shCancel := context.WithTimeout(context.Background(), gracePeriod)
	defer shCancel()
	if err := server.Shutdown(shCtx); err != nil {
		log.Fatalf("graceful shutdown: %v]", err)
	}
}
