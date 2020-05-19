package gstreamer

/*
#cgo pkg-config: gstreamer-1.0
#include "../../include/object.h"
*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

//Object is a GSTObject wrapper interface
type Object interface {
	GetParent() Object
	Unref()

	SetName(string)
	GetName() string

	GetObjectPointer() *C.GstObject

	Set(string, interface{})
	Get(string, interface{})
}

//object ...
type object struct {
	GstObject *C.GstObject
}

func (o *object) GetParent() Object {
	if o.GstObject == nil {
		return nil
	}

	parent := C.gst_object_get_parent(o.GstObject)
	if parent != nil {
		newObj := &object{
			GstObject: parent,
		}

		return newObj
	}

	return nil
}

func (o *object) Unref() {
	C.gst_object_unref(C.gpointer(unsafe.Pointer(o.GstObject)))
}

func (o *object) GetName() (name string) {
	n := C.gst_object_get_name(o.GstObject)

	if n != nil {
		name = C.GoString(n)
	}

	return
}

func (o *object) SetName(name string) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	C.gst_object_set_name(o.GstObject, cname)
}

func (o *object) GetObjectPointer() *C.GstObject {
	return o.GstObject
}

func (o *object) Set(name string, value interface{}) {
	cname := (*C.gchar)(unsafe.Pointer(C.CString(name)))
	defer C.g_free(C.gpointer(unsafe.Pointer(cname)))

	switch value.(type) {
	case string:
		str := (*C.gchar)(unsafe.Pointer(C.CString(value.(string))))
		defer C.g_free(C.gpointer(unsafe.Pointer(str)))
		C.gostreamer_object_set_string(o.GstObject, cname, str)
	case int:
		C.gostreamer_object_set_int(o.GstObject, cname, C.gint(value.(int)))
	case uint32:
		C.gostreamer_object_set_uint(o.GstObject, cname, C.guint(value.(uint32)))
	case float32:
		C.gostreamer_object_set_double(o.GstObject, cname, C.gdouble(value.(float32)))
	case bool:
		var cvalue int
		if value.(bool) == true {
			cvalue = 1
		} else {
			cvalue = 0
		}
		C.gostreamer_object_set_bool(o.GstObject, cname, C.gboolean(cvalue))
	case Caps:
		caps := value.(Caps)
		C.gostreamer_object_set_caps(o.GstObject, cname, caps.GetCapsPointer())
	}
}

func (o *object) Get(name string, value interface{}) {
	cname := (*C.gchar)(unsafe.Pointer(C.CString(name)))
	defer C.g_free(C.gpointer(unsafe.Pointer(cname)))

	switch value.(type) {
	case *string:
		*(value.(*string)) = C.GoString(C.gostreamer_object_get_string(o.GstObject, cname))
	case *int:
		*(value.(*int)) = int(C.gostreamer_object_get_int(o.GstObject, cname))
	case *uint32:
		*(value.(*uint32)) = uint32(C.gostreamer_object_get_uint(o.GstObject, cname))
	case *float32:
		*(value.(*float32)) = float32(C.gostreamer_object_get_double(o.GstObject, cname))
	case *bool:
		*(value.(*bool)) = !(int(C.gostreamer_object_get_bool(o.GstObject, cname)) == 0)
	default:
		fmt.Printf("Type not found %s", reflect.TypeOf(value))
	}
}

func convertPointerToObject(pointer unsafe.Pointer) *C.GstObject {
	return (*C.GstObject)(pointer)
}
