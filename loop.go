package gostreamer

/*
#cgo pkg-config: gstreamer-1.0
#include "include/gostreamer.h"
*/
import "C"

import (
	"runtime"
)

//GMainLoop ...
type GMainLoop struct {
	C *C.GMainLoop
}

//NewMainLoop ...
func NewMainLoop() (loop *GMainLoop) {
	CLoop := C.g_main_loop_new(nil, C.gboolean(0))
	loop = &GMainLoop{C: CLoop}
	runtime.SetFinalizer(loop, func(loop *GMainLoop) {
		C.g_main_loop_unref(loop.C)
	})

	return
}

//Run ...
func (l *GMainLoop) Run() {
	C.g_main_loop_run(l.C)
}

//Quit ...
func (l *GMainLoop) Quit() {
	C.g_main_loop_quit(l.C)
}

//IsRunning ...
func (l *GMainLoop) IsRunning() bool {
	Cbool := C.g_main_loop_is_running(l.C)
	if Cbool == 1 {
		return true
	}
	return false
}
