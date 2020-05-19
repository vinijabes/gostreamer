package gstreamer

/*
#cgo pkg-config: gstreamer-1.0
#include "../../include/gostreamer.h"
*/
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

type Pad interface {
	Object

	Link(Pad) bool
	Unlink(Pad) bool

	GetPadPointer() *C.GstPad
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

var (
	ErrFailedToCreatePadTemplate = errors.New("gstreamer: failed to create pad template")
	ErrFailedToCreatePad         = errors.New("gstreamer: failed to create pad")
)

func NewPad(name *string, direction GstPadDirection) (Pad, error) {
	var cname *C.char
	if name == nil {
		cname = nil
	} else {
		cname := C.CString(*name)
		defer C.free(unsafe.Pointer(cname))
	}

	cpad := C.gst_pad_new(cname, C.GstPadDirection(direction))

	return newPadFromPointer(cpad)
}

func NewPadTemplate(name *string, direction GstPadDirection, presence GstPadPresence, caps Caps) (PadTemplate, error) {
	var cname *C.char
	if name == nil {
		cname = nil
	} else {
		cname := C.CString(*name)
		defer C.free(unsafe.Pointer(cname))
	}

	cpadTemplate := C.gst_pad_template_new(cname, C.GstPadDirection(direction), C.GstPadPresence(presence), caps.GetCapsPointer())
	return newPadTemplateFromPointer(cpadTemplate)
}

func newPadFromPointer(pointer *C.GstPad) (Pad, error) {
	if pointer == nil {
		return nil, ErrFailedToCreatePadTemplate
	}

	pad := &pad{}
	pad.GstObject = convertPointerToObject(unsafe.Pointer(pointer))

	runtime.SetFinalizer(pad, func(p Pad) {
		p.Unref()
	})

	return pad, nil
}

func newPadTemplateFromPointer(pointer *C.GstPadTemplate) (PadTemplate, error) {
	if pointer == nil {
		return nil, ErrFailedToCreatePad
	}

	padTemplate := &padTemplate{}
	padTemplate.GstObject = convertPointerToObject(unsafe.Pointer(pointer))

	runtime.SetFinalizer(padTemplate, func(pt PadTemplate) {
		pt.Unref()
	})

	return padTemplate, nil
}

func (p *pad) Link(other Pad) bool {
	return !(int(C.gst_pad_link(p.GetPadPointer(), other.GetPadPointer())) == 0)
}

func (p *pad) Unlink(other Pad) bool {
	return !(int(C.gst_pad_unlink(p.GetPadPointer(), other.GetPadPointer())) == 0)
}

func (p *pad) GetPadPointer() *C.GstPad {
	return (*C.GstPad)(unsafe.Pointer(p.GstObject))
}

func (pt *padTemplate) GetPadTemplatePointer() *C.GstPadTemplate {
	return (*C.GstPadTemplate)(unsafe.Pointer(pt.GstObject))
}
