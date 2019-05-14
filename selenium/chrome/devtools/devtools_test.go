package main

import (
	"context"
	"encoding/json"
	"fmt"
	. "github.com/aandryashin/matchers"
	"github.com/gorilla/websocket"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/rpcc"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	srv         *httptest.Server
	devtoolsSrv *httptest.Server
)

func init() {
	srv = httptest.NewServer(ws())
	listen = srv.Listener.Addr().String()
	devtoolsSrv = httptest.NewServer(mockDevtoolsMux())
	devtoolsHost = devtoolsSrv.Listener.Addr().String()
}

func mockDevtoolsMux() http.Handler {
	mux := http.NewServeMux()
	version := func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{
		"Browser": "Chrome/72.0.3601.0",
		"Protocol-Version": "1.3",
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3601.0 Safari/537.36",
		"V8-Version": "7.2.233",
		"WebKit-Version": "537.36 (@cfede9db1d154de0468cb0538479f34c0755a0f4)",
		"webSocketDebuggerUrl": "ws://%s/devtools/browser/b0b8a4fb-bb17-4359-9533-a8d9f3908bd8"
	}`, devtoolsHost)))
	}
	mux.HandleFunc("/json/version", version)
	upgrader := websocket.Upgrader{
		CheckOrigin: func(_ *http.Request) bool {
			return true
		},
	}

	mux.HandleFunc("/devtools/browser/", func(w http.ResponseWriter, r *http.Request) {
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

func TestDevtools(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	u := fmt.Sprintf("ws://%s/", srv.Listener.Addr().String())
	conn, err := rpcc.DialContext(ctx, u)
	AssertThat(t, err, Is{nil})
	defer conn.Close()

	c := cdp.NewClient(conn)
	err = c.Page.Enable(ctx)
	AssertThat(t, err, Is{nil})
}

func TestDetectDevtoolsHost(t *testing.T) {
	name, _ := ioutil.TempDir("", "devtools")
	defer os.RemoveAll(name)
	profilePath := filepath.Join(name, ".org.chromium.Chromium.deadbee")
	os.MkdirAll(profilePath, os.ModePerm)
	portFile := filepath.Join(profilePath, "DevToolsActivePort")
	ioutil.WriteFile(portFile, []byte("12345\n/devtools/browser/6f37c7fe-a0a6-4346-a6e2-80c5da0687f0"), 0644)

	AssertThat(t, detectDevtoolsHost(name), EqualTo{"127.0.0.1:12345"})
}
