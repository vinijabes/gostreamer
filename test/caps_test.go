package tests

import (
	"testing"

	"github.com/vinijabes/gostreamer/pkg/gstreamer"
)

func TestNewCaps(t *testing.T) {
	t.Run("NewEmptyCaps", func(t *testing.T) {
		_, err := gstreamer.NewCapsEmpty()
		ok(t, err)
	})

	t.Run("NewAnyCaps", func(t *testing.T) {
		_, err := gstreamer.NewCapsAny()
		ok(t, err)
	})

	t.Run("NewSinkDirectionPad", func(t *testing.T) {
		_, err := gstreamer.NewCapsEmptySimple("video/x-raw")
		ok(t, err)
	})
}
