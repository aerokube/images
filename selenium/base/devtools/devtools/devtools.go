package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aerokube/util"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"net/url"
)

var	devtoolsHost = "127.0.0.1:9222"

func ws() http.Handler {
	return websocket.Server{Handler: page} //Origin checking is turned off
}

func page(wsconn *websocket.Conn){
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
		Host:   devtoolsHost,
		Path:   "/json/list",
	}
	resp, err := http.Get(u.String())

	if err != nil {
		return nil, fmt.Errorf("failed to get debugger url: %v", err)
	}

	var data []map[string]string
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	for _, v := range data {
		if v["type"] == "page" {
			u, err := url.Parse(v["webSocketDebuggerUrl"])
			if err != nil {
				return nil, fmt.Errorf("wrong debugger URL: %v", err)
			}
			return u, nil
		}
	}
	return nil, errors.New("debugger URL information not found")
}
