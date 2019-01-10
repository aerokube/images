package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aerokube/util"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strconv"
	"strings"
)

var (
	paths = struct {
		Version, Json, JsonList, JsonProxy, Page, Devtools, DevtoolsInspector string
	}{
		Version:           "/json/version",
		Json:              "/json",
		JsonList:          "/json/list",
		JsonProxy:         "/json/",
		Page:              "/devtools/page/",
		Devtools:          "/devtools/",
		DevtoolsInspector: "/devtools/inspector.html",
	}

	devtoolsHost = "127.0.0.1:9222"
)

const (
	webSocketDebuggerURL = "webSocketDebuggerUrl"
)

func mux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", devtools)
	return mux
}

func devtools(w http.ResponseWriter, r *http.Request) {
	sid, _ := splitRequestPath(r.URL.Path)
	w.Header().Add("Content-Type", "application/json")
	r.URL.Path = strings.TrimPrefix(r.URL.Path, paths.Devtools+sid)
	devtoolsMux(sid).ServeHTTP(w, r)
}

func splitRequestPath(p string) (string, string) {
	const slash = "/"
	fragments := strings.Split(p, slash)
	return fragments[2], slash + strings.Join(fragments[3:], slash)
}

func devtoolsMux(sid string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(paths.Version, devtoolsVersion(sid))
	mux.HandleFunc(paths.Json, devtoolsTargets(sid))
	mux.HandleFunc(paths.JsonList, devtoolsTargets(sid))
	mux.HandleFunc(paths.JsonProxy, devtoolsProxy(sid))
	ws := websocket.Server{Handler: page(sid)} //Origin checking is turned off
	mux.Handle(paths.Page, ws)
	mux.HandleFunc(paths.DevtoolsInspector, devtoolsProxy(sid))
	return mux
}

func devtoolsVersion(sid string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		(&httputil.ReverseProxy{
			Director: func(r *http.Request) {
				r.URL.Scheme = "http"
				r.URL.Host = devtoolsHost
				log.Printf("[VERSION] [%s]", sid)
			},
			ModifyResponse: func(resp *http.Response) error {
				data := make(map[string]string)
				err := json.NewDecoder(resp.Body).Decode(&data)
				if err != nil {
					return fmt.Errorf("failed to read response body: %v", err)
				}
				newUrl, err := replaceUrlHost(listenHost(r), sid, data[webSocketDebuggerURL])
				if err != nil {
					return fmt.Errorf("failed to replace debugger URL: %v", err)
				}
				data[webSocketDebuggerURL] = newUrl

				b, err := json.Marshal(data)
				if err != nil {
					return fmt.Errorf("failed to encode modified response body: %v", err)
				}
				body := ioutil.NopCloser(bytes.NewReader(b))
				resp.Body = body
				resp.ContentLength = int64(len(b))
				resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
				return nil
			},
			ErrorHandler: defaultErrorHandler(),
		}).ServeHTTP(w, r)
	}
}

func defaultErrorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		user, remote := util.RequestInfo(r)
		log.Printf("[PROXY_ERROR] [%s] [%s] [%v]", user, remote, err)
		w.WriteHeader(http.StatusBadGateway)
	}
}

func listenHost(r *http.Request) string {
	originalHost := r.Header.Get("Host")
	if originalHost != "" {
		return originalHost
	}
	host, port, _ := net.SplitHostPort(listen)
	if host != "" && port != "" {
		return net.JoinHostPort(host, port)
	}
	return net.JoinHostPort("localhost", port)
}

func replaceUrlHost(listenHost, sid, rawUrl string) (string, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", fmt.Errorf("invalid input URL: %v", err)
	}
	u.Host = listenHost
	u.Path = path.Join(paths.Devtools, sid, u.Path)
	return u.String(), nil
}

func devtoolsTargets(sid string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		(&httputil.ReverseProxy{
			Director: func(r *http.Request) {
				r.URL.Scheme = "http"
				r.URL.Host = devtoolsHost
				log.Printf("[TARGETS] [%s]", sid)
			},
			ModifyResponse: func(resp *http.Response) error {
				var data, newData []map[string]string
				err := json.NewDecoder(resp.Body).Decode(&data)
				if err != nil {
					return fmt.Errorf("failed to read response body: %v", err)
				}
				for _, v := range data {
					newUrl, err := replaceUrlHost(listenHost(r), sid, v[webSocketDebuggerURL])
					if err != nil {
						return fmt.Errorf("failed to replace debugger URL: %v", err)
					}
					v[webSocketDebuggerURL] = newUrl
					newData = append(newData, v)
				}

				b, err := json.Marshal(newData)
				if err != nil {
					return fmt.Errorf("failed to encode modified response body: %v", err)
				}
				body := ioutil.NopCloser(bytes.NewReader(b))
				resp.Body = body
				resp.ContentLength = int64(len(b))
				resp.Header.Set("Content-Length", strconv.Itoa(len(b)))

				return nil
			},
			ErrorHandler: defaultErrorHandler(),
		}).ServeHTTP(w, r)
	}
}

func devtoolsProxy(sid string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		(&httputil.ReverseProxy{
			Director: func(r *http.Request) {
				r.URL.Scheme = "http"
				r.URL.Host = devtoolsHost
				log.Printf("[DEVTOOLS_PROXY] [%s] [%s]", sid, r.URL.Path)
			},
			ErrorHandler: defaultErrorHandler(),
		}).ServeHTTP(w, r)
	}
}

func page(sid string) func(*websocket.Conn) {
	return func(wsconn *websocket.Conn) {
		newUrl := *wsconn.Request().URL
		newUrl.Scheme = "ws"
		newUrl.Host = devtoolsHost
		log.Printf("[WEBSOCKET] [%s] [%s]", sid, newUrl.String())
		conn, err := websocket.Dial(newUrl.String(), "", "http://localhost/")
		if err != nil {
			log.Printf("[WEBSOCKET_ERROR] [%v]", err)
			return
		}
		defer conn.Close()
		wsconn.PayloadType = websocket.BinaryFrame
		go func() {
			io.Copy(wsconn, conn)
			wsconn.Close()
			log.Printf("[WEBSOCKET_CLOSED] [%s]", sid)
		}()
		io.Copy(conn, wsconn)
		log.Printf("[WEBSOCKET_CLIENT_DISCONNECTED] [%s]", sid)
	}
}
