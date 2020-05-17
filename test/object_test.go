package tests

import (
	"testing"

	"github.com/vinijabes/gostreamer/pkg/gstreamer"
)

func TestObject(t *testing.T) {
	element, err := gstreamer.NewElement("fakesrc", "fakesrc_test")
	ok(t, err)

	object := element.(gstreamer.Object)

	t.Run("GetName", func(t *testing.T) {
		equals(t, "fakesrc_test", object.GetName())
	})

	t.Run("SetName", func(t *testing.T) {
		object.SetName("new_name")
		equals(t, "new_name", object.GetName())
	})

	t.Run("GetParent", func(t *testing.T) {
		equals(t, nil, object.GetParent())
	})

	t.Run("SetAndGetString", func(t *testing.T) {
		var name string
		object.Set("name", "string_name")

		object.Get("name", &name)
		equals(t, "string_name", name)
	})

	t.Run("SetAndGetInt", func(t *testing.T) {
		var datarate int
		object.Set("datarate", 1024)

		object.Get("datarate", &datarate)
		equals(t, 1024, datarate)
	})

	t.Run("SetAndGetUint", func(t *testing.T) {
		var blocksize uint32
		object.Set("blocksize", 1024)

		object.Get("blocksize", &blocksize)
		equals(t, uint32(1024), blocksize)
	})

	t.Run("SetAndGetBool", func(t *testing.T) {
		var doTimestamp bool
		object.Set("do-timestamp", true)

		object.Get("do-timestamp", &doTimestamp)
		equals(t, true, doTimestamp)
	})
}
