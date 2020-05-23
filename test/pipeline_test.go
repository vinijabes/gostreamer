package tests

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/vinijabes/gostreamer/pkg/gstreamer"
)

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func TestVideoPipeline(t *testing.T) {
	PrintMemUsage()

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

	PrintMemUsage()

	time.Sleep(2 * time.Second)

	result = pipeline.SetState(gstreamer.GstStateNull)
	equals(t, gstreamer.GstStateChangeSuccess, result)

	equals(t, true, pipeline.Remove(src))
	equals(t, true, pipeline.Remove(sink))

	pipeline = nil
	src = nil
	sink = nil
	runtime.GC()
}

func TestAudioPipeline(t *testing.T) {
	PrintMemUsage()

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

	PrintMemUsage()

	time.Sleep(2 * time.Second)

	result = pipeline.SetState(gstreamer.GstStateNull)
	equals(t, gstreamer.GstStateChangeSuccess, result)

	equals(t, true, pipeline.Remove(src))
	equals(t, true, pipeline.Remove(sink))

	pipeline = nil
	src = nil
	sink = nil
	runtime.GC()
}

func TestAppSrcPipeline(t *testing.T) {
	PrintMemUsage()

	pipeline, err := gstreamer.NewPipeline("appsrctest")
	ok(t, err)

	bus, err := pipeline.GetBus()
	ok(t, err)

	src, err := gstreamer.NewElement("appsrc", "source")
	ok(t, err)

	convert, err := gstreamer.NewElement("videoconvert", "convert")
	ok(t, err)

	sink, err := gstreamer.NewElement("autovideosink", "sink")
	ok(t, err)

	src.Set("format", 3)
	src.Set("is-live", true)
	src.Set("do-timestamp", true)

	caps, err := gstreamer.NewCapsFromString("video/x-raw,format=RGB,width=640,height=480,bpp=24,depth=24")
	ok(t, err)

	src.Set("caps", caps)

	equals(t, true, pipeline.Add(src))
	equals(t, true, pipeline.Add(convert))
	equals(t, true, pipeline.Add(sink))

	equals(t, true, src.Link(convert))
	equals(t, true, convert.Link(sink))

	result := pipeline.SetState(gstreamer.GstStatePlaying)
	equals(t, gstreamer.GstStateChangeAsync, result)

	PrintMemUsage()

	time.Sleep(1 * time.Second)

	i := 0
	for {
		if i > 10 {
			break
		}

		data := make([]byte, 640*480*3)
		for j := 0; j < 640*480*3; j++ {
			data[j] = uint8(rand.Intn(255))
		}

		ok(t, src.Push(data))

		i++
		time.Sleep(16 * time.Millisecond)
	}

	for bus.HavePending() {
		message, err := bus.Pop()
		ok(t, err)

		fmt.Println(message.GetName())
	}
	equals(t, false, bus.HavePending())

	result = pipeline.SetState(gstreamer.GstStateNull)
	equals(t, gstreamer.GstStateChangeSuccess, result)

	pipeline = nil
	convert = nil
	sink = nil
	runtime.GC()
}

func TestEncodeDecodePipeline(t *testing.T) {
	PrintMemUsage()

	pipeline, err := gstreamer.NewPipeline("appsrctest")
	ok(t, err)

	src, err := gstreamer.NewElement("videotestsrc", "source")
	ok(t, err)

	enc, err := gstreamer.NewElement("vp8enc", "encoder")
	ok(t, err)

	pay, err := gstreamer.NewElement("rtpvp8pay", "payloader")
	ok(t, err)

	depay, err := gstreamer.NewElement("rtpvp8depay", "depayloader")
	ok(t, err)

	dec, err := gstreamer.NewElement("vp8dec", "decoder")
	ok(t, err)

	sink, err := gstreamer.NewElement("autovideosink", "sink")
	ok(t, err)

	equals(t, true, pipeline.Add(src))
	equals(t, true, pipeline.Add(enc))
	equals(t, true, pipeline.Add(pay))
	equals(t, true, pipeline.Add(depay))
	equals(t, true, pipeline.Add(dec))
	equals(t, true, pipeline.Add(sink))

	equals(t, true, src.Link(enc))
	equals(t, true, enc.Link(pay))
	equals(t, true, pay.Link(depay))
	equals(t, true, depay.Link(dec))
	equals(t, true, dec.Link(sink))

	result := pipeline.SetState(gstreamer.GstStatePlaying)
	equals(t, gstreamer.GstStateChangeAsync, result)

	PrintMemUsage()

	time.Sleep(2 * time.Second)

	result = pipeline.SetState(gstreamer.GstStateNull)
	equals(t, gstreamer.GstStateChangeSuccess, result)

	pipeline = nil
	src = nil
	enc = nil
	depay = nil
	dec = nil
	sink = nil
	runtime.GC()
}
