package gstreamer

/*
#cgo pkg-config: gstreamer-1.0 gstreamer-base-1.0 gstreamer-app-1.0 gstreamer-plugins-base-1.0 gstreamer-video-1.0 gstreamer-audio-1.0 gstreamer-plugins-bad-1.0
#include "gstreamer.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"log"
	"sync"
	"unsafe"
)

//export goPrint
func goPrint(str *C.char) {
	log.Println(C.GoString(str))
}

//Pipeline ...
type Pipeline struct {
	pipeline *C.GstPipeline
	messages chan *Message
	name     string
	id       int
}

//GoHandleSignalCallback ...
type GoHandleSignalCallback func(element *Element)

//GoHandlePadAddedSignalCallback ...
type GoHandlePadAddedSignalCallback func(element *Element, pad *Pad)

type callback struct {
	padAdded GoHandlePadAddedSignalCallback
}

//Element ...
type Element struct {
	element   *C.GstElement
	factory   string
	name      string
	id        int
	callbacks callback
}

//Pad ...
type Pad struct {
	pad *C.GstPad
}

//PadTemplate ...
type PadTemplate struct {
	padTemplate *C.GstPadTemplate
}

//Caps ...
type Caps struct {
	caps *C.GstCaps
}

//Structure ...
type Structure struct {
	structure *C.GstStructure
}

//Bus ...
type Bus struct {
	bus *C.GstBus
}

//Message ...
type Message struct {
	message *C.GstMessage
}
type MessageType int

//MessageType constants
const (
	MessageUnknown      MessageType = C.GST_MESSAGE_UNKNOWN
	MessageEOS          MessageType = C.GST_MESSAGE_EOS
	MessageError        MessageType = C.GST_MESSAGE_ERROR
	MessageWarning      MessageType = C.GST_MESSAGE_WARNING
	MessageInfo         MessageType = C.GST_MESSAGE_INFO
	MessageTag          MessageType = C.GST_MESSAGE_TAG
	MessageBuffering    MessageType = C.GST_MESSAGE_BUFFERING
	MessageStateChanged MessageType = C.GST_MESSAGE_STATE_CHANGED
	MessageAny          MessageType = C.GST_MESSAGE_ANY
)

var pipelines = make(map[int]*Pipeline)
var elements = make(map[int]*Element)
var lock sync.Mutex
var gstIDGenerate = 10000

func init() {
	log.Println("Gstreamer Initializing")
	C.gstreamer_init()
}

//NewPipeline ...
func NewPipeline(name string) (*Pipeline, error) {
	pipelineStrUnsafe := C.CString(name)
	defer C.free(unsafe.Pointer(pipelineStrUnsafe))
	cpipeline := C.gstreamer_create_pipeline(pipelineStrUnsafe)
	if cpipeline == nil {
		return nil, errors.New("create pipeline error")
	}

	pipeline := &Pipeline{
		pipeline: cpipeline,
		name:     name,
	}

	lock.Lock()
	defer lock.Unlock()
	gstIDGenerate++
	pipeline.id = gstIDGenerate
	pipelines[pipeline.id] = pipeline

	return pipeline, nil
}

//NewElement ...
func NewElement(factory string, name string) (*Element, error) {
	elementFactoryStrUnsafe := C.CString(factory)
	elementNameStrUnsafe := C.CString(name)
	defer C.free(unsafe.Pointer(elementFactoryStrUnsafe))
	defer C.free(unsafe.Pointer(elementNameStrUnsafe))

	celement := C.gstreamer_element_factory_make(elementFactoryStrUnsafe, elementNameStrUnsafe)
	if celement == nil {
		return nil, fmt.Errorf("create element error(%s, %s)", factory, name)
	}

	element := &Element{
		element: celement,
		name:    name,
		factory: factory,
	}

	lock.Lock()
	defer lock.Unlock()
	gstIDGenerate++
	element.id = gstIDGenerate
	elements[element.id] = element

	return element, nil
}

//GetType ...
func (m *Message) GetType() MessageType {
	c := C.toGstMessageType(unsafe.Pointer(m.message))
	return MessageType(c)
}

//GetTimestamp ...
func (m *Message) GetTimestamp() uint64 {
	c := C.messageTimestamp(unsafe.Pointer(m.message))
	return uint64(c)
}

//GetTypeName ...
func (m *Message) GetTypeName() string {
	c := C.messageTypeName(unsafe.Pointer(m.message))
	return C.GoString(c)
}

//Start ...
func (p *Pipeline) Start() {
	C.gstreamer_pipeline_start(p.pipeline)
}

//Pause ...
func (p *Pipeline) Pause() {
	C.gstreamer_pipeline_pause(p.pipeline)
}

//Stop ...
func (p *Pipeline) Stop() {
	C.gstreamer_pipeline_stop(p.pipeline)
}

//SendEOS ...
func (p *Pipeline) SendEOS() {
	C.gstreamer_pipeline_sendeos(p.pipeline)
}

//Add ...
func (p *Pipeline) Add(e *Element) {
	C.gstreamer_bin_add_element(p.pipeline, e.element)
}

//PullMessage ...
func (p *Pipeline) PullMessage() <-chan *Message {
	p.messages = make(chan *Message, 5)
	C.gstreamer_pipeline_bus_watch(p.pipeline, C.int(p.id))
	return p.messages
}

//export goHandleBusMessage
func goHandleBusMessage(message *C.GstMessage, pipelineID C.int) {
	lock.Lock()
	defer lock.Unlock()

	msg := &Message{message: message}

	if pipeline, ok := pipelines[int(pipelineID)]; ok {
		log.Println(msg.GetTypeName())
		if pipeline.messages != nil {
			pipeline.messages <- msg
		}
	} else {
		fmt.Printf("discarding message, no pipeline with id %d", int(pipelineID))
	}
}

//Link ...
func (e *Element) Link(dest *Element) error {
	result := C.gstreamer_element_link(e.element, dest.element)

	if int(result) == 0 {
		return fmt.Errorf("link element failed(%s, %s)", e.name, dest.name)
	}

	return nil
}

//Set ...
func (e *Element) Set(property string, value string) {
	elementPropertyStrUnsafe := C.CString(property)
	elementValueStrUnsafe := C.CString(value)
	defer C.free(unsafe.Pointer(elementPropertyStrUnsafe))
	defer C.free(unsafe.Pointer(elementValueStrUnsafe))

	C.gstreamer_object_set(e.element, elementPropertyStrUnsafe, elementValueStrUnsafe)
}

//GetStaticPad ...
func (e *Element) GetStaticPad(padName string) (*Pad, error) {
	elementPadNameStrUnsafe := C.CString(padName)
	defer C.free(unsafe.Pointer(elementPadNameStrUnsafe))

	cpad := C.gst_element_get_static_pad(e.element, elementPadNameStrUnsafe)

	if cpad == nil {
		return nil, fmt.Errorf("GetStaticPad with name %s failed", padName)
	}

	pad := &Pad{pad: cpad}
	return pad, nil
}

//ConnectPadAddedSignal ...
func (e *Element) ConnectPadAddedSignal(cb GoHandlePadAddedSignalCallback) {
	e.callbacks.padAdded = cb
	C.gstreamer_element_pad_added_signal_connect(e.element, C.int(e.id))
}

//SetCapsFromString ...
func (e *Element) SetCapsFromString(caps string) {
	capsStr := C.CString(caps)
	defer C.free(unsafe.Pointer(capsStr))
	C.gstreamer_set_caps(e.element, capsStr)
}

//GetPadTemplate ...
func (e *Element) GetPadTemplate(name string) (*PadTemplate, error) {
	cpadTemplate := C.gstreamer_get_pad_template(e.element, C.CString(name))
	if cpadTemplate == nil {
		return nil, errors.New("create pad template error")
	}

	padTemplate := &PadTemplate{padTemplate: cpadTemplate}

	return padTemplate, nil
}

//RequestPad ...
func (e *Element) RequestPad(template *PadTemplate) (*Pad, error) {
	cpad := C.gstreamer_element_request_pad(e.element, template.padTemplate)
	if cpad == nil {
		return nil, errors.New("create pad error")
	}
	pad := &Pad{pad: cpad}
	return pad, nil
}

//export goHandlePadAddedSignal
func goHandlePadAddedSignal(elementID C.int, cpad *C.GstPad) {
	lock.Lock()
	defer lock.Unlock()

	pad := &Pad{pad: cpad}

	if element, ok := elements[int(elementID)]; ok {
		if element.callbacks.padAdded != nil {
			element.callbacks.padAdded(element, pad)
		}
	}
}

//Link ...
func (p *Pad) Link(dest *Pad) int {
	ret := int(C.gstreamer_pad_link(p.pad, dest.pad))
	return int(ret)
}

//GetCurrentCaps ...
func (p *Pad) GetCurrentCaps() *Caps {
	ccaps := C.gst_pad_get_current_caps(p.pad)
	if ccaps == nil {
		return nil
	}

	caps := &Caps{caps: ccaps}

	return caps
}

//GetStructure ...
func (c *Caps) GetStructure() *Structure {
	cstructure := C.gst_caps_get_structure(c.caps, C.uint(0))
	if cstructure == nil {
		return nil
	}
	structure := &Structure{structure: cstructure}

	return structure
}

//Unref ...
func (c *Caps) Unref() {
	C.gst_caps_unref(c.caps)
}

//GetName ...
func (s *Structure) GetName() string {
	res := C.GoString(C.gst_structure_get_name(s.structure))
	return res
}

//ScanPathForPlugins ...
func ScanPathForPlugins(directory string) {
	C.gst_registry_scan_path(C.gst_registry_get(), C.CString(directory))
}

//CheckPlugins ...
func CheckPlugins(plugins []string) error {

	var plugin *C.GstPlugin
	var registry *C.GstRegistry

	registry = C.gst_registry_get()

	for _, pluginstr := range plugins {
		plugincstr := C.CString(pluginstr)
		plugin = C.gst_registry_find_plugin(registry, plugincstr)
		C.free(unsafe.Pointer(plugincstr))
		if plugin == nil {
			return fmt.Errorf("Required gstreamer plugin %s not found", pluginstr)
		}
	}

	return nil
}
