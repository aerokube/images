package main

import (
	"context"
	"encoding/json"
	"fmt"
	. "github.com/aandryashin/matchers"
	. "github.com/aandryashin/matchers/httpresp"
	"github.com/gorilla/websocket"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/rpcc"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	srv         *httptest.Server
	devtoolsSrv *httptest.Server
)

func init() {
	srv = httptest.NewServer(mux())
	listen = srv.Listener.Addr().String()
	devtoolsSrv = httptest.NewServer(mockDevtoolsMux())
	devtoolsHost = devtoolsSrv.Listener.Addr().String()
}

func mockDevtoolsMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/json/version", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(`{
		"Browser": "Chrome/72.0.3601.0",
		"Protocol-Version": "1.3",
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3601.0 Safari/537.36",
		"V8-Version": "7.2.233",
		"WebKit-Version": "537.36 (@cfede9db1d154de0468cb0538479f34c0755a0f4)",
		"webSocketDebuggerUrl": "ws://localhost:9222/devtools/browser/b0b8a4fb-bb17-4359-9533-a8d9f3908bd8"
		}`))
	})
	targets := func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(`[ {
		"description": "",
		"devtoolsFrontendUrl": "/devtools/inspector.html?ws=localhost:9222/devtools/page/DAB7FB6187B554E10B0BD18821265734",
		"id": "DAB7FB6187B554E10B0BD18821265734",
		"title": "Yahoo",
		"type": "page",
		"url": "https://www.yahoo.com/",
		"webSocketDebuggerUrl": "ws://localhost:9222/devtools/page/DAB7FB6187B554E10B0BD18821265734"
		} ]`))
	}
	mux.HandleFunc("/json", targets)
	mux.HandleFunc("/json/list", targets)
	mux.HandleFunc("/json/protocol/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte("{}"))
	})
	mux.HandleFunc("/json/new", func(_ http.ResponseWriter, _ *http.Request) {})
	mux.HandleFunc("/json/activate/", func(_ http.ResponseWriter, _ *http.Request) {})
	mux.HandleFunc("/json/close/", func(_ http.ResponseWriter, _ *http.Request) {})
	mux.HandleFunc("/devtools/inspector.html", func(_ http.ResponseWriter, _ *http.Request) {})
	upgrader := websocket.Upgrader{
		CheckOrigin: func(_ *http.Request) bool {
			return true
		},
	}
	mux.HandleFunc("/devtools/page/", func(w http.ResponseWriter, r *http.Request) {
		//Echo request ID but omit Method
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}
		defer c.Close()
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				break
			}
			type req struct {
				ID uint64 `json:"id"`
			}
			var r req
			err = json.Unmarshal(message, &r)
			if err != nil {
				panic(err)
			}
			output, err := json.Marshal(r)
			if err != nil {
				panic(err)
			}
			err = c.WriteMessage(mt, output)
			if err != nil {
				break
			}
		}
	})
	return mux
}

func withUrl(path string) string {
	return srv.URL + path
}

func TestVersion(t *testing.T) {
	resp, err := http.Get(withUrl("/devtools/test-session/json/version"))
	AssertThat(t, err, Is{nil})
	data := make(map[string]string)
	err = json.NewDecoder(resp.Body).Decode(&data)
	AssertThat(t, err, Is{nil})
	wsUrl, ok := data["webSocketDebuggerUrl"]
	AssertThat(t, ok, Is{true})
	AssertThat(t, wsUrl, EqualTo{fmt.Sprintf("ws://%s/devtools/test-session/devtools/browser/b0b8a4fb-bb17-4359-9533-a8d9f3908bd8", srv.Listener.Addr().String())})
}

func TestStatic(t *testing.T) {
	resp, err := http.Get(withUrl("/devtools/test-session/json/protocol/"))
	AssertThat(t, err, Is{nil})
	AssertThat(t, resp, Code{200})

	resp, err = http.Get(withUrl("/devtools/test-session/devtools/inspector.html"))
	AssertThat(t, err, Is{nil})
	AssertThat(t, resp, Code{200})
}

func TestDevtools(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	devt := devtool.New(withUrl("/devtools/test-session"))
	pt, err := devt.Get(ctx, devtool.Page)
	AssertThat(t, err, Is{nil})
	conn, err := rpcc.DialContext(ctx, pt.WebSocketDebuggerURL)
	AssertThat(t, err, Is{nil})
	defer conn.Close()

	c := cdp.NewClient(conn)
	err = c.Page.Enable(ctx)
	AssertThat(t, err, Is{nil})
}

func TestReplaceUrlHost(t *testing.T) {
	newUrl, err := replaceUrlHost("example.com:4444", "test-session-id", "ws://localhost:9222/devtools/browser/b0b8a4fb-bb17-4359-9533-a8d9f3908bd8")
	AssertThat(t, err, Is{nil})
	AssertThat(t, newUrl, EqualTo{"ws://example.com:4444/devtools/test-session-id/devtools/browser/b0b8a4fb-bb17-4359-9533-a8d9f3908bd8"})
}
