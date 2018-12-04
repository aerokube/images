package main

import (
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
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
		fmt.Printf("Failed to wait for WM: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Starting command: %v\n", os.Args[1:])
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Failed to start command: %v\n", err)
		os.Exit(1)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	cmd.Process.Signal(syscall.SIGINT)
	fmt.Println("Stopping command")
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Failed to stop command: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)

}

func waitWM() error {
	x, err := xgbutil.NewConn()
	if err != nil {
		return fmt.Errorf("failed to connect to display: %v", err)
	}
	
	if wm, err := ewmh.GetEwmhWM(x); err != nil {
		fmt.Printf("Detected running WM %s\n", wm)
		return nil
	}

	fmt.Println("Waiting for WM to start")
	ch := make(chan struct{}, 1)
	err = xwindow.New(x, x.RootWin()).Listen(xproto.EventMaskSubstructureNotify, xproto.EventMaskSubstructureRedirect)
	if err != nil {
		return fmt.Errorf("failed to listen for X events: %v", err)
	}
	h := func(_ *xgbutil.XUtil, _ xevent.ClientMessageEvent) {
		close(ch)
	}
	xevent.ClientMessageFun(h).Connect(x, x.RootWin())
	go xevent.Main(x)
	select {
	case <-ch:
		return nil
	case <-time.After(60*time.Second):
		return fmt.Errorf("timed out waiting for WM to start")
	}
}
