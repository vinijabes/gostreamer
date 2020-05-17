package gstreamer

/*
#cgo CFLAGS: -I ../../include
#cgo pkg-config: gstreamer-1.0
#include "../../include/element.h"
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

//Element is a element gstreamer wrapper
type Element interface {
	Object

	Link(Element) bool
	Unlink(Element) bool

	GetElementPointer() *C.GstElement
}

//Factory is a wrapper to element factory
type Factory interface {
	Object

	Create(name string) Element
	GetElementFactoryPointer() *C.GstElementFactory
}

type element struct {
	object
}

type elementFactory struct {
	object
}

//NewElement ...
func NewElement(factory string, name string) (Element, error) {
	elementFactoryStrUnsafe := C.CString(factory)
	elementNameStrUnsafe := C.CString(name)
	defer C.free(unsafe.Pointer(elementFactoryStrUnsafe))
	defer C.free(unsafe.Pointer(elementNameStrUnsafe))

	celement := C.gst_element_factory_make(elementFactoryStrUnsafe, elementNameStrUnsafe)
	if celement == nil {
		return nil, fmt.Errorf("create element error(%s, %s)", factory, name)
	}

	element := &element{
		object: object{
			GstObject: convertPointerToObject(unsafe.Pointer(celement)),
		},
	}

	runtime.SetFinalizer(element, func(e Element) {
		e.Unref()
	})

	return element, nil
}

//NewElementFactory ...
func NewElementFactory(name string) (Factory, error) {
	factoryNameStrUnsafe := C.CString(name)
	defer C.free(unsafe.Pointer(factoryNameStrUnsafe))

	celement := C.gst_element_factory_find(factoryNameStrUnsafe)
	if celement == nil {
		return nil, fmt.Errorf("Failed to find factory %s", name)
	}

	factory := &elementFactory{
		object: object{
			GstObject: convertPointerToObject(unsafe.Pointer(celement)),
		},
	}

	return factory, nil
}

func (e *element) Link(other Element) bool {
	return !(int(C.gst_element_link(e.GetElementPointer(), other.GetElementPointer())) == 0)
}

func (e *element) Unlink(other Element) bool {
	return !(int(C.gst_element_link(e.GetElementPointer(), other.GetElementPointer())) == 0)
}

func (e *element) GetElementPointer() *C.GstElement {
	return (*C.GstElement)(unsafe.Pointer(e.GstObject))
}

func (ef *elementFactory) Create(name string) Element {
	elementNameStrUnsafe := C.CString(name)
	defer C.free(unsafe.Pointer(elementNameStrUnsafe))

	celement := C.gst_element_factory_create(ef.GetElementFactoryPointer(), elementNameStrUnsafe)

	element := &element{
		object: object{
			GstObject: convertPointerToObject(unsafe.Pointer(celement)),
		},
	}

	runtime.SetFinalizer(element, func(e Element) {
		ef.Unref()
	})

	return element
}

func (ef *elementFactory) GetElementFactoryPointer() *C.GstElementFactory {
	return (*C.GstElementFactory)(unsafe.Pointer(ef.GstObject))
}
