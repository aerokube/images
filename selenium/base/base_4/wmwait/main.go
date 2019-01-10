package main

// To compile this you have to install libx11-dev package in addition to Golang

// #cgo pkg-config: x11
/*

#include <X11/Xlib.h>
#include <stdio.h>
#include <stdlib.h>

int retcode = -1;

int OnWMDetected(Display* display, XErrorEvent* e) {
    retcode = 0;
    return 0;
}

int check() {
    Display *display = XOpenDisplay(NULL);
    if (display == NULL) {
        return 1;
    }
    XSetErrorHandler(&OnWMDetected);
    XSelectInput(display, DefaultRootWindow(display), SubstructureRedirectMask | SubstructureNotifyMask);
    XSync(display, 1);
    return retcode;
}
*/
import "C"

import (
	"fmt"
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
	fmt.Println("Waiting for WM to start")
	ch := make(chan struct{}, 1)
	go func(ch chan struct{}) {
		for {
			code := C.check()
			if code == 0 {
				close(ch)
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}(ch)
	select {
	case <-ch:
		return nil
	case <-time.After(60 * time.Second):
		return fmt.Errorf("timed out waiting for WM to start")
	}
}
