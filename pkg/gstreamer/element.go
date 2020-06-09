package gstreamer

/*
#cgo CFLAGS: -I ../../include
#cgo pkg-config: gstreamer-1.0
#include "../../include/element.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"unsafe"
)

//Element is a element gstreamer wrapper
type Element interface {
	Object

	Link(Element) bool
	Unlink(Element)

	SetState(GstState) GstStateChangeReturn

	GetPadTemplate(string) (PadTemplate, error)
	GetStaticPad(padName string) (Pad, error)
	RequestPad(PadTemplate, *string, Caps) (Pad, error)

	GetBus() (Bus, error)

	Push(buffer []byte) error

	SetOnPadAddedCallback(PadAddedCallback)
	SetOnPadRemovedCallback(PadRemovedCallback)

	GetElementPointer() *C.GstElement
}

//Factory is a wrapper to element factory
type Factory interface {
	Object

	Create(name string) Element
	GetElementFactoryPointer() *C.GstElementFactory
}

type ElementSignalCallback struct {
	element      Element
	callbackID   uint64
	callbackFunc interface{}

	handlerID uint64
}

var (
	callbackID  uint64 = 0
	callbackMap        = map[uint64]*ElementSignalCallback{}
	mutex       sync.Mutex
)

type element struct {
	object

	onPadAdded   *ElementSignalCallback
	onPadRemoved *ElementSignalCallback
}

type elementFactory struct {
	object
}

type GstState int
type GstStateChangeReturn int

type PadAddedCallback func(Element, Pad)
type PadRemovedCallback func(Element, Pad)

const (
	GstStateVoidPending GstState = iota
	GstStateNull
	GstStateReady
	GstStatePaused
	GstStatePlaying
)

const (
	GstStateChangeFailure GstStateChangeReturn = iota
	GstStateChangeSuccess
	GstStateChangeAsync
	GstStateChangeNoPreroll
)

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
			needUnref: true,
		},
	}

	runtime.SetFinalizer(element, func(e Element) {
		if e.IsAutoUnrefEnabled() {
			e.Unref()
		}
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

func (e *element) Unlink(other Element) {
	C.gst_element_unlink(e.GetElementPointer(), other.GetElementPointer())
}

func (e *element) SetState(state GstState) GstStateChangeReturn {
	result := GstStateChangeReturn(C.gst_element_set_state(e.GetElementPointer(), C.GstState(state)))
	return result
}

func (e *element) GetPadTemplate(name string) (PadTemplate, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cpadTemplate := C.gst_element_get_pad_template(e.GetElementPointer(), cname)
	return newPadTemplateFromPointer(cpadTemplate)
}

func (e *element) GetStaticPad(padName string) (Pad, error) {
	elementPadNameStrUnsafe := C.CString(padName)
	defer C.free(unsafe.Pointer(elementPadNameStrUnsafe))

	cpad := C.gst_element_get_static_pad(e.GetElementPointer(), elementPadNameStrUnsafe)
	return newPadFromPointer(cpad)
}

func (e *element) RequestPad(template PadTemplate, name *string, caps Caps) (Pad, error) {
	var cname *C.char
	var ccaps *C.GstCaps

	if name == nil {
		cname = nil
	} else {
		cname = C.CString(*name)
	}

	if caps == nil {
		ccaps = nil
	}

	cpad := C.gst_element_request_pad(
		e.GetElementPointer(),
		template.GetPadTemplatePointer(),
		cname,
		ccaps,
	)

	if cname != nil {
		C.free(unsafe.Pointer(cname))
	}

	return newPadFromPointer(cpad)
}

func (e *element) GetBus() (Bus, error) {
	cbus := C.gst_element_get_bus(e.GetElementPointer())
	return newBusFromPointer(cbus)
}

func (e *element) SetOnPadAddedCallback(cb PadAddedCallback) {
	elementCallback := &ElementSignalCallback{
		element:      e,
		callbackID:   callbackID,
		callbackFunc: cb,
	}

	handlerID := C.gostreamer_add_pad_added_signal(e.GetElementPointer(), C.guint64(callbackID))
	elementCallback.handlerID = uint64(handlerID)

	e.onPadAdded = elementCallback

	mutex.Lock()
	callbackMap[callbackID] = elementCallback
	mutex.Unlock()

	callbackID++
}

func (e *element) SetOnPadRemovedCallback(cb PadRemovedCallback) {
	elementCallback := &ElementSignalCallback{
		element:      e,
		callbackID:   callbackID,
		callbackFunc: cb,
	}

	handlerID := C.gostreamer_add_pad_removed_signal(e.GetElementPointer(), C.guint64(callbackID))
	elementCallback.handlerID = uint64(handlerID)

	e.onPadRemoved = elementCallback

	mutex.Lock()
	callbackMap[callbackID] = elementCallback
	mutex.Unlock()

	callbackID++
}

func (e *element) GetElementPointer() *C.GstElement {
	return (*C.GstElement)(unsafe.Pointer(e.GstObject))
}

func (e *element) Push(buffer []byte) (err error) {
	b := C.CBytes(buffer)
	defer C.free(unsafe.Pointer(b))
	gstReturn := C.gostreamer_element_push_buffer(e.GetElementPointer(), b, C.int(len(buffer)))

	if gstReturn != C.GST_FLOW_OK {
		err = errors.New("could not push buffer on appsrc element")
		return
	}

	return
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

//export go_pad_added_callback
func go_pad_added_callback(celement *C.GstElement, cpad *C.GstPad, callbackID C.guint64) {
	mutex.Lock()
	callback, ok := callbackMap[uint64(callbackID)]
	mutex.Unlock()

	if ok {
		pad, err := newPadFromPointer(cpad)
		if err != nil {
			return
		}
		pad.DisableAutoUnref()

		padAdded := callback.callbackFunc.(PadAddedCallback)
		padAdded(callback.element, pad)
	}
}

//export go_pad_removed_callback
func go_pad_removed_callback(celement *C.GstElement, cpad *C.GstPad, callbackID C.guint64) {
	mutex.Lock()
	callback, ok := callbackMap[uint64(callbackID)]
	mutex.Unlock()

	if ok {
		pad, err := newPadFromPointer(cpad)
		if err != nil {
			return
		}
		pad.DisableAutoUnref()

		padRemoved := callback.callbackFunc.(PadRemovedCallback)
		padRemoved(callback.element, pad)
	}
}
