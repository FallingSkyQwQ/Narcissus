package ui

import (
	"testing"

	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

func TestListSingleSelectionReplacesPrevious(t *testing.T) {
	l := NewList().Items("a", "b", "c")

	l.Select(0)
	l.Select(2)

	if l.SelectedIndex() != 2 {
		t.Errorf("selected index = %d, want 2", l.SelectedIndex())
	}
	if l.IsSelected(0) {
		t.Error("selecting a second item should clear the first in single mode")
	}
	if l.SelectedValue() != "c" {
		t.Errorf("selected value = %q, want c", l.SelectedValue())
	}
}

func TestListMultipleSelection(t *testing.T) {
	l := NewList().Items("a", "b", "c").SelectionMode(ListMultipleSelection)

	l.Select(0)
	l.Select(2)

	indices := l.SelectedIndices()
	if len(indices) != 2 || indices[0] != 0 || indices[1] != 2 {
		t.Errorf("selected indices = %v, want [0 2]", indices)
	}
}

func TestListSelectionModeSwitchKeepsFirst(t *testing.T) {
	l := NewList().Items("a", "b", "c").SelectionMode(ListMultipleSelection)
	l.Select(1)
	l.Select(2)

	l.SelectionMode(ListSingleSelection)

	if got := l.SelectedIndices(); len(got) != 1 || got[0] != 1 {
		t.Errorf("indices after switching to single = %v, want [1]", got)
	}
}

func TestListItemsTruncationClearsStaleSelection(t *testing.T) {
	l := NewList().Items("a", "b", "c")
	l.Select(2)
	l.Items("a")

	if l.SelectedIndex() != -1 {
		t.Errorf("selected index = %d, want -1 after shrinking items", l.SelectedIndex())
	}
}

func TestListOnSelectCallback(t *testing.T) {
	l := NewList().Items("a", "b")
	var gotIndex int
	var gotValue string
	l.OnSelect(func(index int, value string) { gotIndex, gotValue = index, value })

	l.Select(1)
	if gotIndex != 1 || gotValue != "b" {
		t.Errorf("callback = (%d, %q), want (1, \"b\")", gotIndex, gotValue)
	}
}

func TestListRenderProps(t *testing.T) {
	l := NewList().Items("a", "b").Select(1)
	if kindOf(l) != ControlList {
		t.Errorf("kindOf = %v, want list", kindOf(l))
	}
	props := propsOf(l)
	if len(props.Items) != 2 || props.Selected != 1 {
		t.Errorf("unexpected list props: %+v", props)
	}
}

func TestTabsSelectSwapsContent(t *testing.T) {
	first := NewText("first")
	second := NewText("second")
	tabs := NewTabs().AddTab("One", first).AddTab("Two", second)

	// The first tab is selected on creation.
	if tabs.CurrentIndex() != 0 || tabs.CurrentTitle() != "One" {
		t.Fatalf("initial tab = %d/%q, want 0/One", tabs.CurrentIndex(), tabs.CurrentTitle())
	}
	if got := tabs.body.GetChildren(); len(got) != 1 || got[0] != Widget(first) {
		t.Fatalf("body should show the first tab content")
	}

	tabs.Select(1)
	got := tabs.body.GetChildren()
	if len(got) != 1 || got[0] != Widget(second) {
		t.Errorf("body should show the second tab content after Select(1)")
	}
	if tabs.CurrentTitle() != "Two" {
		t.Errorf("current title = %q, want Two", tabs.CurrentTitle())
	}
}

func TestTabsOnChangeFires(t *testing.T) {
	tabs := NewTabs().AddTab("One", NewText("1")).AddTab("Two", NewText("2"))
	var changes []int
	tabs.OnChange(func(index int) { changes = append(changes, index) })

	tabs.Select(1)
	if len(changes) != 1 || changes[0] != 1 {
		t.Errorf("onChange = %v, want [1]", changes)
	}
}

func TestTabsMeasureAndLayout(t *testing.T) {
	tabs := NewTabs().AddTab("One", NewText("content"))

	size := tabs.Measure(flex.Constraint{MaxWidth: 400, MaxHeight: 300})
	if size.Width <= 0 || size.Height <= 0 {
		t.Fatalf("tab size = %gx%g, want positive", size.Width, size.Height)
	}

	tabs.Layout(0, 0, 400, 300)
	headerRect := tabs.header.GetRect()
	bodyRect := tabs.body.GetRect()
	if bodyRect.Position.Y < headerRect.Size.Height {
		t.Errorf("body starts at y=%g, want at least header height %g",
			bodyRect.Position.Y, headerRect.Size.Height)
	}
}

func TestScrollViewContent(t *testing.T) {
	content := NewText("scroll me")
	sv := NewScrollView().SetContent(content).Size(200, 100)

	if sv.Content() != Widget(content) {
		t.Fatal("content not stored")
	}
	if len(sv.GetChildren()) != 1 {
		t.Fatalf("scroll view children = %d, want 1", len(sv.GetChildren()))
	}

	size := sv.Measure(flex.Constraint{})
	if size.Width != 200 || size.Height != 100 {
		t.Errorf("size = %gx%g, want 200x100", size.Width, size.Height)
	}
	if kindOf(sv) != ControlScroll {
		t.Errorf("kindOf = %v, want scroll", kindOf(sv))
	}
}

func TestScrollViewSetContentReplaces(t *testing.T) {
	sv := NewScrollView().SetContent(NewText("first"))
	sv.SetContent(NewText("second"))

	if len(sv.GetChildren()) != 1 {
		t.Errorf("children = %d, want 1 after replacing content", len(sv.GetChildren()))
	}
}

func TestListDeselectAndClearSelection(t *testing.T) {
	l := NewList().Items("a", "b", "c").SelectionMode(ListMultipleSelection)
	l.Select(0)
	l.Select(1)
	l.Select(2)

	l.Deselect(1)
	if l.IsSelected(1) {
		t.Error("Deselect did not clear the selection")
	}
	if got := l.SelectedIndices(); len(got) != 2 || got[0] != 0 || got[1] != 2 {
		t.Errorf("indices after deselect = %v, want [0 2]", got)
	}

	l.ClearSelection()
	if got := l.SelectedIndices(); len(got) != 0 {
		t.Errorf("indices after clear = %v, want empty", got)
	}
	if l.SelectedIndex() != -1 || l.SelectedValue() != "" {
		t.Errorf("selection after clear = %d/%q, want -1/\"\"", l.SelectedIndex(), l.SelectedValue())
	}
}

func TestListSelectIgnoresOutOfRange(t *testing.T) {
	l := NewList().Items("a")
	l.Select(-1)
	l.Select(5)
	if l.SelectedIndex() != -1 {
		t.Errorf("selected index = %d, want -1 for out-of-range selects", l.SelectedIndex())
	}
}

func TestListSelectedIndicesAreSorted(t *testing.T) {
	l := NewList().Items("a", "b", "c").SelectionMode(ListMultipleSelection)
	l.Select(2)
	l.Select(0)

	if got := l.SelectedIndices(); len(got) != 2 || got[0] != 0 || got[1] != 2 {
		t.Errorf("selected indices = %v, want ascending [0 2]", got)
	}
}

func TestListRenderPropsSelectionMode(t *testing.T) {
	single := NewList().Items("a")
	if props := propsOf(single); props.SelectionMode != int(ListSingleSelection) {
		t.Errorf("single-mode props = %d, want %d", props.SelectionMode, int(ListSingleSelection))
	}

	multi := NewList().Items("a").SelectionMode(ListMultipleSelection)
	if props := propsOf(multi); props.SelectionMode != int(ListMultipleSelection) {
		t.Errorf("multi-mode props = %d, want %d", props.SelectionMode, int(ListMultipleSelection))
	}
}

func TestListClickHandledWhenEnabled(t *testing.T) {
	l := NewList().Items("a")
	if !l.HandleEvent(newTestClickEvent(l)) {
		t.Error("click should be handled by an enabled list")
	}

	l.SetEnabled(false)
	if l.HandleEvent(newTestClickEvent(l)) {
		t.Error("disabled list should not handle clicks")
	}
}

func TestListItemsReturnsCopy(t *testing.T) {
	l := NewList().Items("a", "b")
	items := l.GetItems()
	items[0] = "mutated"
	if l.GetItems()[0] != "a" {
		t.Error("GetItems did not return a copy")
	}
}

func TestScrollViewMeasureDefaults(t *testing.T) {
	size := NewScrollView().Measure(flex.Constraint{})
	if size.Width != 200 || size.Height != 150 {
		t.Errorf("default size = %gx%g, want 200x150", size.Width, size.Height)
	}

	sized := NewScrollView().Size(120, 60)
	size = sized.Measure(flex.Constraint{})
	if size.Width != 120 || size.Height != 60 {
		t.Errorf("explicit size = %gx%g, want 120x60", size.Width, size.Height)
	}
}

func TestScrollViewLayoutGivesContentNaturalHeight(t *testing.T) {
	SetTextMeasurer(fixedMeasurer{perRune: 10, height: 7})
	defer SetTextMeasurer(nil)

	content := NewText("scroll me")                         // 9 runes * 10 = 90 wide, 7 tall
	sv := NewScrollView().SetContent(content).Size(200, 30) // viewport shorter than content
	sv.Layout(5, 10, 200, 30)

	rect := content.GetRect()
	if rect.Position.X != 5 || rect.Position.Y != 10 {
		t.Errorf("content position = %v, want (5, 10)", rect.Position)
	}
	if rect.Size.Width != 200 {
		t.Errorf("content width = %g, want 200", rect.Size.Width)
	}
	// The content keeps its natural height; the overflow is
	// handled by the native scrollbars.
	if rect.Size.Height != 7 {
		t.Errorf("content height = %g, want natural 7", rect.Size.Height)
	}
}

func TestScrollViewSetContentNilClears(t *testing.T) {
	sv := NewScrollView().SetContent(NewText("x"))
	sv.SetContent(nil)

	if sv.Content() != nil {
		t.Error("content should be nil after SetContent(nil)")
	}
	if len(sv.GetChildren()) != 0 {
		t.Errorf("children = %d, want 0 after clearing content", len(sv.GetChildren()))
	}
}
