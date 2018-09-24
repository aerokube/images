#include <X11/Xlib.h>
#include <stdio.h>
#include <stdlib.h>

int OnWMDetected(Display* display, XErrorEvent* e) {
    exit(0);
}

int main() {
    Display *display = XOpenDisplay(NULL);
    if (display == NULL) {
        exit(1);
    }
    XSetErrorHandler(&OnWMDetected);
    XSelectInput(display, DefaultRootWindow(display), SubstructureRedirectMask | SubstructureNotifyMask);
    XSync(display, 1);
    exit(-1);
}