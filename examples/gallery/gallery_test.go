package main

import (
	"strings"
	"testing"

	"github.com/FallingSkyQwQ/Narcissus/pkg/ui"
)

// TestBuildConstructsEveryPage is a headless smoke test: it assembles the full
// widget tree without a backend and checks the navigation wiring.
func TestBuildConstructsEveryPage(t *testing.T) {
	g := &gallery{}
	g.pages = g.definePages()

	root := g.build()
	if root == nil {
		t.Fatal("build returned nil root")
	}
	if len(g.pageUI) != len(g.pages) {
		t.Fatalf("page widgets = %d, want %d", len(g.pageUI), len(g.pages))
	}

	visible := 0
	for _, w := range g.pageUI {
		if w.GetVisible() {
			visible++
		}
	}
	if visible != 1 {
		t.Errorf("visible pages = %d, want exactly 1", visible)
	}
}

func TestShowTogglesPageVisibility(t *testing.T) {
	g := &gallery{}
	g.pages = g.definePages()
	g.build()

	g.show(3)
	if g.current != 3 {
		t.Fatalf("current = %d, want 3", g.current)
	}
	if !g.pageUI[3].GetVisible() {
		t.Error("selected page is not visible")
	}
	if g.pageUI[0].GetVisible() {
		t.Error("previous page is still visible")
	}

	// Out-of-range indices are ignored.
	g.show(999)
	if g.current != 3 {
		t.Errorf("current changed for out-of-range index: %d", g.current)
	}
}

func TestAccessibilityDump(t *testing.T) {
	submit := ui.NewButton().Text("Sign in")
	submit.SetAccessibleRole(ui.RoleButton)
	submit.SetAccessibleDescription("Submits the form")
	root := ui.NewContainer().Add(submit)

	dump := dumpAccessibility(ui.AccessibilityTree(root), 0)
	if !strings.Contains(dump, "button: Sign in") {
		t.Errorf("dump missing button node:\n%s", dump)
	}
	if !strings.Contains(dump, "Submits the form") {
		t.Errorf("dump missing description:\n%s", dump)
	}
}
