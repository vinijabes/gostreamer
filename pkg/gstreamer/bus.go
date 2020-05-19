package gstreamer

/*
#cgo CFLAGS: -I ../../include
#cgo pkg-config: gstreamer-1.0
#include "../../include/pch.h"
*/
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

type Bus interface {
	Object

	HavePending() bool
	Pop() (Message, error)

	GetBusPointer() *C.GstBus
}

type bus struct {
	object
}

var (
	ErrFailedToCreateBus     = errors.New("failed to create bus")
	ErrFailedToCreateMessage = errors.New("failed to create message")
)

func NewBus() (Bus, error) {
	cbus := C.gst_bus_new()
	return newBusFromPointer(cbus)
}

func newBusFromPointer(pointer *C.GstBus) (Bus, error) {
	if pointer == nil {
		return nil, ErrFailedToCreateBus
	}

	bus := &bus{}
	bus.GstObject = convertPointerToObject(unsafe.Pointer(pointer))

	runtime.SetFinalizer(bus, func(b Bus) {
		b.Unref()
	})

	return bus, nil
}

func (b *bus) HavePending() bool {
	return C.gst_bus_have_pending(b.GetBusPointer()) != 0
}

func (b *bus) Pop() (Message, error) {
	cmessage := C.gst_bus_pop(b.GetBusPointer())
	return newMessageFromPointer(cmessage)
}

func (b *bus) GetBusPointer() *C.GstBus {
	return (*C.GstBus)(unsafe.Pointer(b.GstObject))
}
