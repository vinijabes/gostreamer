package tests

import (
	"testing"

	"github.com/vinijabes/gostreamer/pkg/gstreamer"
)

func TestElement(t *testing.T) {
	element, err := gstreamer.NewElement("compositor", "element")
	ok(t, err)

	template, err := element.GetPadTemplate("sink_%u")
	ok(t, err)

	_, err = element.RequestPad(template, nil, nil)
	ok(t, err)
}

func TestFactory(t *testing.T) {
	factory, err := gstreamer.NewElementFactory("videotestsrc")
	ok(t, err)

	equals(t, "videotestsrc", factory.GetName())

	t.Run("CreateElement", func(t *testing.T) {
		element := factory.Create("element")

		equals(t, "element", element.GetName())
	})
}
