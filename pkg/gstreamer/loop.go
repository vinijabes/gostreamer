package gstreamer

/*
#cgo pkg-config: gstreamer-1.0
#include "../../include/gostreamer.h"
*/
import "C"

import (
	"runtime"
)

//MainLoop ...
type MainLoop struct {
	C *C.GMainLoop
}

//NewMainLoop ...
func NewMainLoop() (loop *MainLoop) {
	CLoop := C.g_main_loop_new(nil, C.gboolean(0))
	loop = &MainLoop{C: CLoop}
	runtime.SetFinalizer(loop, func(loop *MainLoop) {
		C.g_main_loop_unref(loop.C)
	})

	return
}

//Run ...
func (l *MainLoop) Run() {
	C.g_main_loop_run(l.C)
}

//Quit ...
func (l *MainLoop) Quit() {
	C.g_main_loop_quit(l.C)
}

//IsRunning ...
func (l *MainLoop) IsRunning() bool {
	Cbool := C.g_main_loop_is_running(l.C)
	if Cbool == 1 {
		return true
	}
	return false
}
