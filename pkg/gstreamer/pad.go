package gstreamer

/*
#cgo pkg-config: gstreamer-1.0
#include "../../include/gostreamer.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type Pad interface {
	Object
}

type pad struct {
	object
}

type PadTemplate interface {
	Object

	GetPadTemplatePointer() *C.GstPadTemplate
}

type padTemplate struct {
	object
}

type GstPadDirection int
type GstPadPresence int

const (
	GstPadUnknown GstPadDirection = iota
	GstPadSrc
	GstPadSink
)

const (
	GstPadAlways GstPadPresence = iota
	GstPadSometimes
	GstPadRequest
)

func NewPad(name *string, direction GstPadDirection) (Pad, error) {
	pad := &pad{}

	if name == nil {
		cpad := C.gst_pad_new(nil, C.GstPadDirection(direction))
		if cpad == nil {
			return nil, fmt.Errorf("failed to create Pad")
		}

		pad.GstObject = convertPointerToObject(unsafe.Pointer(cpad))
	} else {
		cname := C.CString(*name)
		defer C.free(unsafe.Pointer(cname))

		cpad := C.gst_pad_new(nil, C.GstPadDirection(direction))
		if cpad == nil {
			return nil, fmt.Errorf("failed to create Pad")
		}
		pad.GstObject = convertPointerToObject(unsafe.Pointer(cpad))
	}

	return pad, nil
}

func NewPadTemplate(name *string, direction GstPadDirection, presence GstPadPresence, caps Caps) (PadTemplate, error) {
	padTemplate := &padTemplate{}

	if name == nil {
		cpadTemplate := C.gst_pad_template_new(nil, C.GstPadDirection(direction), C.GstPadPresence(presence), caps.GetCapsPointer())
		if cpadTemplate == nil {
			return nil, fmt.Errorf("failed to create PadTemplate")
		}

		padTemplate.GstObject = convertPointerToObject(unsafe.Pointer(cpadTemplate))
	} else {
		cpadName := C.CString(*name)
		defer C.free(unsafe.Pointer(cpadName))

		cpadTemplate := C.gst_pad_template_new(cpadName, C.GstPadDirection(direction), C.GstPadPresence(presence), caps.GetCapsPointer())
		if cpadTemplate == nil {
			return nil, fmt.Errorf("failed to create PadTemplate")
		}

		padTemplate.GstObject = convertPointerToObject(unsafe.Pointer(cpadTemplate))
	}

	return padTemplate, nil
}

func newPadFromPointer(pointer *C.GstPad) Pad {
	pad := &pad{}
	pad.GstObject = convertPointerToObject(unsafe.Pointer(pointer))

	return pad
}

func newPadTemplateFromPointer(pointer *C.GstPadTemplate) PadTemplate {
	padTemplate := &padTemplate{}
	padTemplate.GstObject = convertPointerToObject(unsafe.Pointer(pointer))

	return padTemplate
}

func (pt *padTemplate) GetPadTemplatePointer() *C.GstPadTemplate {
	return (*C.GstPadTemplate)(unsafe.Pointer(pt.GstObject))
}