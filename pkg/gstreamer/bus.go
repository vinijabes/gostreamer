package gstreamer

/*
#cgo CFLAGS: -I ../../include
#cgo pkg-config: gstreamer-1.0
#include "../../include/pch.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Bus interface{
	Object

	HavePending() bool
	GetBusPointer() *C.GstBus
}

type bus struct {
	object
}

func NewBus() (Bus, error) {
	cbus := C.gst_bus_new()

	if cbus == nil {
		return nil, fmt.Errorf("failed to create bus")
	}

	bus := &bus{}
	bus.GstObject = convertPointerToObject(unsafe.Pointer(cbus))
	return bus, nil
}

func (b *bus) HavePending() bool {
	return C.gst_bus_have_pending(b.GetBusPointer()) != 0
}

func (b *bus) GetBusPointer() *C.GstBus {
	return (*C.GstBus)(unsafe.Pointer(b.GstObject))
}