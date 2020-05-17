package gstreamer

/*
#cgo CFLAGS: -I ../../include
#cgo pkg-config: gstreamer-1.0
#include "../../include/pch.h"
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

type Bin interface {
	Element

	Add(Element) bool
	Remove(Element) bool

	GetBinPointer() *C.GstBin
}

type bin struct {
	element
}

//NewBin ...
func NewBin(name string) (Bin, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cbin := C.gst_bin_new(cname)

	if cbin == nil {
		return nil, fmt.Errorf("failed to create bin: %s", name)
	}

	bin := &bin{
		element: element{
			object: object{
				GstObject: convertPointerToObject(unsafe.Pointer(cbin)),
			},
		},
	}

	runtime.SetFinalizer(bin, func(b Bin) {
		b.Unref()
	})

	return bin, nil
}

func (b *bin) Add(element Element) bool {
	return !(int(C.gst_bin_add(b.GetBinPointer(), element.GetElementPointer())) == 0)
}

func (b *bin) Remove(element Element) bool {
	return !(int(C.gst_bin_remove(b.GetBinPointer(), element.GetElementPointer())) == 0)
}

func (b *bin) GetBinPointer() *C.GstBin {
	return (*C.GstBin)(unsafe.Pointer(b.GstObject))
}
