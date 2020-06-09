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
	GetStructure(int) Structure

	GetCapsPointer() *C.GstCaps
}

type caps struct {
	GstCaps *C.GstCaps
}

func NewCapsAny() (Caps, error) {
	ccaps := C.gst_caps_new_any()
	return newCapsFromPointer(ccaps)
}

func NewCapsEmpty() (Caps, error) {
	ccaps := C.gst_caps_new_empty()
	return newCapsFromPointer(ccaps)
}

func NewCapsEmptySimple(mediaType string) (Caps, error) {
	cmediaType := C.CString(mediaType)
	defer C.free(unsafe.Pointer(cmediaType))

	ccaps := C.gst_caps_new_empty_simple(cmediaType)
	return newCapsFromPointer(ccaps)
}

func NewCapsFromString(caps string) (Caps, error) {
	c := (*C.gchar)(unsafe.Pointer(C.CString(caps)))
	defer C.g_free(C.gpointer(unsafe.Pointer(c)))
	ccaps := C.gst_caps_from_string(c)
	return newCapsFromPointer(ccaps)
}

func newCapsFromPointer(pointer *C.GstCaps) (Caps, error) {
	if pointer == nil {
		return nil, fmt.Errorf("failed to create Caps")
	}

	caps := &caps{}
	caps.GstCaps = pointer

	runtime.SetFinalizer(caps, func(c Caps) {
		c.Unref()
	})

	return caps, nil
}

func (c *caps) Unref() {
	C.gst_caps_unref(c.GstCaps)
}

func (c *caps) GetStructure(index int) Structure {
	cstructure := C.gst_caps_get_structure(c.GetCapsPointer(), C.guint(index))
	structure, err := newStructureFromPointer(cstructure)

	if err != nil {
		return nil
	}

	return structure
}

func (c *caps) GetCapsPointer() *C.GstCaps {
	return c.GstCaps
}
