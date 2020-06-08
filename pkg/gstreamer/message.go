package gstreamer

/*
#cgo CFLAGS: -I ../../include
#cgo pkg-config: gstreamer-1.0
#include "../../include/message.h"
*/
import "C"
import (
	"runtime"
)

type Message interface {
	GetType() MessageType
	GetName() string

	GetStructure() Structure

	Unref()
	GetMessagePointer() *C.GstMessage
}

type message struct {
	GstMessage *C.GstMessage
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

func newMessageFromPointer(pointer *C.GstMessage) (Message, error) {
	if pointer == nil {
		return nil, ErrFailedToCreateMessage
	}

	message := &message{}
	message.GstMessage = pointer

	runtime.SetFinalizer(message, func(m Message) {
		m.Unref()
	})

	return message, nil
}

func (m *message) GetType() MessageType {
	return MessageType(C.gostreamer_get_message_type(m.GetMessagePointer()))
}

func (m *message) GetName() string {
	messageType := m.GetType()
	cname := C.gst_message_type_get_name(C.GstMessageType(messageType))
	name := C.GoString(cname)

	return name
}

func (m *message) GetStructure() Structure {
	cstructure := C.gst_message_get_structure(m.GstMessage)

	str := &structure{}
	str.GstStructure = cstructure

	return str
}

func (m *message) Unref() {
	C.gst_message_unref(m.GetMessagePointer())
}

func (m *message) GetMessagePointer() *C.GstMessage {
	return m.GstMessage
}
