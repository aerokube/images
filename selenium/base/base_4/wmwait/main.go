package main

import (
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := waitWM()
	if err != nil {
		fmt.Printf("Failed to wait for WM: %v", err)
		os.Exit(1)
	}
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Failed to start command: %v", err)
		os.Exit(1)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		cmd.Process.Signal(syscall.SIGINT)
		err := cmd.Wait()
		if err != nil {
			fmt.Printf("Failed to wait for command to stop: %v", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

}

func waitWM() error {
	x, err := xgbutil.NewConn()
	if err != nil {
		return fmt.Errorf("failed to connect to display: %v", err)
	}
	ch := make(chan struct{}, 1)
	err = xwindow.New(x, x.RootWin()).Listen(xproto.EventMaskSubstructureNotify, xproto.EventMaskSubstructureRedirect)
	if err != nil {
		return fmt.Errorf("failed to listen for X events: %v", err)
	}
	h := func(_ *xgbutil.XUtil, _ xevent.ClientMessageEvent) {
		close(ch)
	}
	xevent.ClientMessageFun(h).Connect(x, x.RootWin())
	xevent.Main(x)
	select {
	case <-ch:
		return nil
	case <-time.After(60*time.Second):
		return fmt.Errorf("timed out waiting for WM to start")
	}
}
