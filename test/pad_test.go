package tests

import (
	"testing"

	"github.com/vinijabes/gostreamer/pkg/gstreamer"
)

func TestNewPad(t *testing.T) {
	t.Run("NewUnknownDirectionPad", func(t *testing.T) {
		name := "name"
		_, err := gstreamer.NewPad(&name, gstreamer.GstPadUnknown)
		ok(t, err)
	})

	t.Run("NewSrcDirectionPad", func(t *testing.T) {
		name := "name"
		_, err := gstreamer.NewPad(&name, gstreamer.GstPadSrc)
		ok(t, err)
	})

	t.Run("NewSinkDirectionPad", func(t *testing.T) {
		name := "name"
		_, err := gstreamer.NewPad(&name, gstreamer.GstPadSink)
		ok(t, err)
	})

	t.Run("CreateUnnamedPad", func(t *testing.T) {
		_, err := gstreamer.NewPad(nil, gstreamer.GstPadUnknown)
		ok(t, err)
	})
}

func TestPad(t *testing.T) {

}
