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

//Pipeline is a wrapper to GstPipeline
type Pipeline interface {
	Bin

	GetPipelinePointer() *C.GstPipeline
}

type pipeline struct {
	bin
}

//NewPipeline ...
func NewPipeline(name string) (Pipeline, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cpipeline := C.gst_pipeline_new(cname)

	if cpipeline == nil {
		return nil, fmt.Errorf("failed to create bin: %s", name)
	}

	pipeline := &pipeline{
		bin: bin{
			element: element{
				object: object{
					GstObject: convertPointerToObject(unsafe.Pointer(cpipeline)),
				},
			},
		},
	}

	runtime.SetFinalizer(pipeline, func(p Pipeline) {
		p.Unref()
	})

	return pipeline, nil
}

func (p *pipeline) GetPipelinePointer() *C.GstPipeline {
	return (*C.GstPipeline)(unsafe.Pointer(p.GstObject))
}
