package gstreamer

/*
#cgo pkg-config: gstreamer-1.0
#include "../../include/gostreamer.h"
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

type Caps interface {
	Unref()
	GetCapsPointer() *C.GstCaps
}

type caps struct {
	GstCaps *C.GstCaps
}

func NewCapsAny() (Caps, error) {
	ccaps := C.gst_caps_new_any()
	if ccaps == nil {
		return nil, fmt.Errorf("failed to create Caps")
	}

	caps := &caps{}
	caps.GstCaps = ccaps

	runtime.SetFinalizer(caps, func(c Caps) {
		c.Unref()
	})

	return caps, nil
}

func NewCapsEmpty() (Caps, error) {
	ccaps := C.gst_caps_new_empty()
	if ccaps == nil {
		return nil, fmt.Errorf("failed to create Caps")
	}

	caps := &caps{}
	caps.GstCaps = ccaps

	runtime.SetFinalizer(caps, func(c Caps) {
		c.Unref()
	})

	return caps, nil
}

func NewCapsEmptySimple(mediaType string) (Caps, error) {
	cmediaType := C.CString(mediaType)
	defer C.free(unsafe.Pointer(cmediaType))

	ccaps := C.gst_caps_new_empty_simple(cmediaType)
	if ccaps == nil {
		return nil, fmt.Errorf("failed to create Caps")
	}

	caps := &caps{}
	caps.GstCaps = ccaps

	runtime.SetFinalizer(caps, func(c Caps) {
		c.Unref()
	})

	return caps, nil
}

func (c *caps) Unref() {
	C.gst_caps_unref(c.GstCaps)
}

func (c *caps) GetCapsPointer() *C.GstCaps {
	return c.GstCaps
}
