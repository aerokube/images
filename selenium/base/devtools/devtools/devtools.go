package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aerokube/util"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var devtoolsHost = "127.0.0.1:9222"

func ws() http.Handler {
	return websocket.Server{Handler: page} //Origin checking is turned off
}

func page(wsconn *websocket.Conn) {
	_, remote := util.RequestInfo(wsconn.Request())
	u, err := getDebuggerUrl()
	if err != nil {
		log.Printf("[WEBSOCKET_URL_ERROR] [%v]", err)
		return
	}
	conn, err := websocket.Dial(u.String(), "", "http://localhost/")
	if err != nil {
		log.Printf("[WEBSOCKET_CONNECTION_ERROR] [%v]", err)
		return
	}
	log.Printf("[WEBSOCKET] [%s] [%s]", remote, u.String())
	defer conn.Close()
	wsconn.PayloadType = websocket.BinaryFrame
	go func() {
		io.Copy(wsconn, conn)
		wsconn.Close()
		log.Printf("[WEBSOCKET_CLOSED] [%s]", remote)
	}()
	io.Copy(conn, wsconn)
	log.Printf("[WEBSOCKET_CLIENT_DISCONNECTED] [%s]", remote)
}

func getDebuggerUrl() (*url.URL, error) {
	u := url.URL{
		Scheme: "http",
		Host:   detectDevtoolsHost("/tmp"),
		Path:   "/json/version",
	}
	resp, err := http.Get(u.String())

	if err != nil {
		return nil, fmt.Errorf("failed to get debugger url: %v", err)
	}

	var version map[string]string
	err = json.NewDecoder(resp.Body).Decode(&version)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	wsUrl, err := url.Parse(version["webSocketDebuggerUrl"])
	if err == nil {
		return wsUrl, nil
	}
	return nil, errors.New("debugger URL information not found")
}

func detectDevtoolsHost(baseDir string) string {
	candidates, err := filepath.Glob(filepath.Join(baseDir, ".org.chromium.Chromium*"))
	if err == nil {
		for _, c := range candidates {
			f, err := os.Stat(c)
			if err != nil {
				continue
			}
			if !f.IsDir() {
				continue
			}
			portFile := filepath.Join(c, "DevToolsActivePort")
			data, err := ioutil.ReadFile(portFile)
			if err != nil {
				continue
			}
			lines := strings.Split(string(data), "\n")
			if len(lines) == 0 {
				continue
			}
			port, err := strconv.Atoi(lines[0])
			if err != nil {
				continue
			}
			return fmt.Sprintf("127.0.0.1:%d", port)
		}
	}
	return devtoolsHost
}
