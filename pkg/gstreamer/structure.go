package gstreamer

/*
#cgo CFLAGS: -I ../../include
#cgo pkg-config: gstreamer-1.0
#include "../../include/message.h"
*/
import "C"
import "errors"

type Structure interface {
}

type structure struct {
	GstStructure *C.GstStructure
}

var ErrFailedToCreateStructure = errors.New("failed to create structure")

func newStructureFromPointer(pointer *C.GstStructure) (Structure, error) {
	if pointer == nil {
		return nil, ErrFailedToCreateStructure
	}

	structure := &structure{}
	structure.GstStructure = pointer

	return structure, nil
}

func (s *structure) String() string {
	return C.GoString(C.gst_structure_to_string(s.GstStructure))
}
