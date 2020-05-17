package tests

import (
	"testing"

	"github.com/vinijabes/gostreamer/pkg/gstreamer"
)

func TestBin(t *testing.T) {
	bin, err := gstreamer.NewBin("testbin")
	ok(t, err)

	src, err := gstreamer.NewElement("videotestsrc", "testsrc")
	ok(t, err)

	sink, err := gstreamer.NewElement("fakesink", "tessink")
	ok(t, err)

	equals(t, true, bin.Add(src))
	equals(t, true, bin.Add(sink))

	equals(t, true, bin.Remove(src))
	equals(t, true, bin.Remove(sink))

	equals(t, false, bin.Remove(src))
}
