package tests

import (
	"testing"
	"time"

	"github.com/vinijabes/gostreamer/pkg/gstreamer"
)

func TestVideoPipeline(t *testing.T) {
	pipeline, err := gstreamer.NewPipeline("videotest")
	ok(t, err)

	src, err := gstreamer.NewElement("videotestsrc", "source")
	ok(t, err)

	sink, err := gstreamer.NewElement("autovideosink", "sink")
	ok(t, err)

	equals(t, true, pipeline.Add(src))
	equals(t, true, pipeline.Add(sink))

	equals(t, true, src.Link(sink))

	result := pipeline.SetState(gstreamer.GstStatePlaying)
	equals(t, gstreamer.GstStateChangeAsync, result)

	time.Sleep(2 * time.Second)

	result = pipeline.SetState(gstreamer.GstStateNull)
	equals(t, gstreamer.GstStateChangeSuccess, result)
}

func TestAudioPipeline(t *testing.T) {
	pipeline, err := gstreamer.NewPipeline("audiotest")
	ok(t, err)

	src, err := gstreamer.NewElement("audiotestsrc", "source")
	ok(t, err)

	sink, err := gstreamer.NewElement("autoaudiosink", "sink")
	ok(t, err)

	equals(t, true, pipeline.Add(src))
	equals(t, true, pipeline.Add(sink))

	equals(t, true, src.Link(sink))

	result := pipeline.SetState(gstreamer.GstStatePlaying)
	equals(t, gstreamer.GstStateChangeAsync, result)

	time.Sleep(2 * time.Second)

	result = pipeline.SetState(gstreamer.GstStateNull)
	equals(t, gstreamer.GstStateChangeSuccess, result)
}
