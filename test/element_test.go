package tests

import (
	"testing"

	"github.com/vinijabes/gostreamer/pkg/gstreamer"
)

func TestFactory(t *testing.T) {
	factory, err := gstreamer.NewElementFactory("videotestsrc")
	ok(t, err)

	equals(t, "videotestsrc", factory.GetName())

	t.Run("CreateElement", func(t *testing.T) {
		element := factory.Create("element")

		equals(t, "element", element.GetName())
	})
}
