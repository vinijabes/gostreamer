package tests

import (
	"testing"

	"github.com/vinijabes/gostreamer/pkg/gstreamer"
)

func TestBus(t *testing.T) {
	bus, err := gstreamer.NewBus()
	ok(t, err)

	equals(t, false, bus.HavePending())
}
