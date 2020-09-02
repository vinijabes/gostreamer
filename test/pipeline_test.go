package tests

import (
	"errors"
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

		fmt.Println(message.GetStructure())
	}
	equals(t, false, bus.HavePending())

	result = pipeline.SetState(gstreamer.GstStateNull)
	equals(t, gstreamer.GstStateChangeSuccess, result)

	pipeline = nil
	convert = nil
	sink = nil
	runtime.GC()
}

func TestAppSinkPipeline(t *testing.T) {
	PrintMemUsage()

	pipeline, err := gstreamer.NewPipeline("appsinktest")
	ok(t, err)

	src, err := gstreamer.NewElement("videotestsrc", "source")
	ok(t, err)

	sink, err := gstreamer.NewElement("appsink", "sink")
	ok(t, err)

	src.Set("is-live", true)
	src.Set("do-timestamp", true)

	caps, err := gstreamer.NewCapsFromString("video/x-raw,format=RGB,width=640,height=480,bpp=24,depth=24")
	ok(t, err)

	sink.Set("caps", caps)
	sink.Set("emit-signals", true)

	samples := make(chan gstreamer.Sample)

	sink.SetOnSampleAddedCallback(func(e gstreamer.Element, sample gstreamer.Sample) {
		samples <- sample
	})

	equals(t, true, pipeline.Add(src))
	equals(t, true, pipeline.Add(sink))

	equals(t, true, src.Link(sink))

	result := pipeline.SetState(gstreamer.GstStatePlaying)
	equals(t, gstreamer.GstStateChangeAsync, result)

	PrintMemUsage()

	select {
	case <-samples:
		break
	case <-time.After(3 * time.Second):
		ok(t, errors.New("Timed Out"))
		break
	}

	result = pipeline.SetState(gstreamer.GstStateNull)
	equals(t, gstreamer.GstStateChangeSuccess, result)

	pipeline = nil
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

func TestUriDecodePipeline(t *testing.T) {
	PrintMemUsage()

	pipeline, err := gstreamer.NewPipeline("urisrctest")
	ok(t, err)

	src, err := gstreamer.NewElement("uridecodebin", "source")
	ok(t, err)

	audioconvert, err := gstreamer.NewElement("audioconvert", "audioconvert")
	ok(t, err)

	audioresample, err := gstreamer.NewElement("audioresample", "resample")
	ok(t, err)

	audiosink, err := gstreamer.NewElement("autoaudiosink", "audiosink")
	ok(t, err)

	videoconvert, err := gstreamer.NewElement("videoconvert", "videoconvert")
	ok(t, err)

	videosink, err := gstreamer.NewElement("autovideosink", "videosink")
	ok(t, err)

	equals(t, true, pipeline.Add(src))
	equals(t, true, pipeline.Add(audioconvert))
	equals(t, true, pipeline.Add(audioresample))
	equals(t, true, pipeline.Add(audiosink))
	equals(t, true, pipeline.Add(videoconvert))
	equals(t, true, pipeline.Add(videosink))

	equals(t, true, audioconvert.Link(audioresample))
	equals(t, true, audioresample.Link(audiosink))
	equals(t, true, videoconvert.Link(videosink))

	src.Set("uri", "https://www.freedesktop.org/software/gstreamer-sdk/data/media/sintel_trailer-480p.webm")

	src.SetOnPadAddedCallback(func(element gstreamer.Element, pad gstreamer.Pad) {
		fmt.Println("PadAdded")
		var sinkPad gstreamer.Pad

		if pad.GetCurrentCaps().GetStructure(0).GetName() == "audio/x-raw" {
			sinkPad, err = audioconvert.GetStaticPad("sink")
		} else {
			sinkPad, err = videoconvert.GetStaticPad("sink")
		}
		ok(t, err)

		result := pad.Link(sinkPad)
		equals(t, gstreamer.GstPadLinkOk, result)
	})

	result := pipeline.SetState(gstreamer.GstStatePlaying)
	equals(t, gstreamer.GstStateChangeAsync, result)

	PrintMemUsage()

	time.Sleep(10 * time.Second)

	result = pipeline.SetState(gstreamer.GstStateNull)
	equals(t, gstreamer.GstStateChangeSuccess, result)

	pipeline = nil
	src = nil
	audioconvert = nil
	audioresample = nil
	audiosink = nil
	runtime.GC()
}

func TestRtspPipeline(t *testing.T) {
	PrintMemUsage()

	pipeline, err := gstreamer.NewPipeline("rtspsrcpipeline")
	ok(t, err)

	src, err := gstreamer.NewElement("rtspsrc", "source")
	ok(t, err)

	decode, err := gstreamer.NewElement("decodebin", "videodecoder")
	ok(t, err)

	timeoverlay, err := gstreamer.NewElement("timeoverlay", "timeoverlay")
	ok(t, err)

	convert, err := gstreamer.NewElement("videoconvert", "videoconvert")
	ok(t, err)

	sink, err := gstreamer.NewElement("autovideosink", "videosink")
	ok(t, err)

	equals(t, true, pipeline.Add(src))
	equals(t, true, pipeline.Add(decode))
	equals(t, true, pipeline.Add(timeoverlay))
	equals(t, true, pipeline.Add(convert))
	equals(t, true, pipeline.Add(sink))

	equals(t, true, convert.Link(timeoverlay))
	equals(t, true, timeoverlay.Link(sink))

	src.Set("location", "rtsp://177.188.103.117:8554/0")
	src.Set("latency", 1000)

	src.SetOnPadAddedCallback(func(element gstreamer.Element, pad gstreamer.Pad) {
		fmt.Println("PadAddedSrc")
		var sinkPad gstreamer.Pad

		sinkPad, err = decode.GetStaticPad("sink")
		ok(t, err)

		result := pad.Link(sinkPad)
		equals(t, gstreamer.GstPadLinkOk, result)
	})

	decode.SetOnPadAddedCallback(func(element gstreamer.Element, pad gstreamer.Pad) {
		fmt.Println("PadAddedDecodeBin")
		var sinkPad gstreamer.Pad

		sinkPad, err = convert.GetStaticPad("sink")
		ok(t, err)

		result := pad.Link(sinkPad)
		equals(t, gstreamer.GstPadLinkOk, result)
	})

	result := pipeline.SetState(gstreamer.GstStatePlaying)
	equals(t, gstreamer.GstStateChangeAsync, result)

	PrintMemUsage()

	time.Sleep(60 * time.Second)

	result = pipeline.SetState(gstreamer.GstStateNull)
	equals(t, gstreamer.GstStateChangeSuccess, result)

	pipeline = nil
	src = nil
	decode = nil
	convert = nil
	sink = nil
	runtime.GC()
}
