package main

import (
	"os"
	"testing"
	"time"

	"github.com/FallingSkyQwQ/Narcissus/pkg/ui"
)

// TestInteractive walks every page and toggles the theme on a real display,
// then quits. It is opt-in because it needs a windowing system:
//
//	NARC_GALLERY_INTERACTIVE=1 xvfb-run -a go test ./examples/gallery -run TestInteractive -v
func TestInteractive(t *testing.T) {
	if os.Getenv("NARC_GALLERY_INTERACTIVE") != "1" {
		t.Skip("set NARC_GALLERY_INTERACTIVE=1 and run under a display")
	}

	// GtkApplication parses os.Args; hide the test flags from it.
	os.Args = []string{"narcissus-gallery-selftest"}

	g := &gallery{}
	g.pages = g.definePages()
	g.app = ui.NewApp("Gallery self-test", ui.WithSize(1000, 760))
	g.root = g.build()
	g.app.SetContent(g.root)

	post := func(fn func()) {
		if backend := ui.CurrentBackend(); backend != nil {
			_ = backend.Post(fn)
		}
	}

	go func() {
		time.Sleep(1500 * time.Millisecond)
		for i := range g.pages {
			index := i
			post(func() { g.show(index) })
			time.Sleep(120 * time.Millisecond)
		}
		post(func() { g.setDark(true) })
		time.Sleep(500 * time.Millisecond)
		post(func() { g.setDark(false) })
		time.Sleep(500 * time.Millisecond)
		post(func() { g.app.Quit() })
	}()

	if err := g.app.Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
}
