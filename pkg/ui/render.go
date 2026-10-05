package ui

import (
	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

// kindOf returns the native control kind that backs a widget.
func kindOf(w Widget) ControlKind {
	switch w.(type) {
	case *Container:
		return ControlContainer
	case *Button:
		return ControlButton
	case *Text:
		return ControlText
	case *Checkbox:
		return ControlCheckbox
	case *ComboBox:
		return ControlComboBox
	case *Slider:
		return ControlSlider
	case *Image:
		return ControlImage
	case *TextInput:
		return ControlTextInput
	case *ProgressBar:
		return ControlProgress
	case *Switch:
		return ControlSwitch
	case *RadioButton:
		return ControlRadio
	case *List:
		return ControlList
	case *ScrollView:
		return ControlScroll
	case *Menu:
		return ControlMenu
	default:
		return ControlContainer
	}
}

// propsOf extracts a toolkit-neutral property snapshot from a widget,
// including the resolved accessibility semantics shared by every widget.
func propsOf(w Widget) ControlProps {
	props := widgetProps(w)
	props.AccessibleName = accessibleNameFor(w)
	props.AccessibleDescription = accessibleDescriptionFor(w)
	props.AccessibleRole = accessibleRoleFor(w)
	return props
}

// widgetProps extracts the kind-specific properties of a widget.
func widgetProps(w Widget) ControlProps {
	switch v := w.(type) {
	case *Button:
		return ControlProps{Text: v.GetText()}
	case *Text:
		return ControlProps{Text: v.GetText()}
	case *Checkbox:
		return ControlProps{Text: v.GetLabel(), Checked: v.IsChecked()}
	case *ComboBox:
		return ControlProps{
			Items:       v.GetItems(),
			Selected:    v.GetSelectedIndex(),
			Placeholder: v.GetPlaceholder(),
		}
	case *Slider:
		return ControlProps{
			Value: float64(v.GetValue()),
			Min:   float64(v.GetMin()),
			Max:   float64(v.GetMax()),
			Step:  float64(v.GetStep()),
		}
	case *Image:
		return ControlProps{Source: v.GetSource(), Fit: int(v.GetFit())}
	case *TextInput:
		return ControlProps{
			Text:        v.GetValue(),
			Placeholder: v.GetPlaceholder(),
			Multiline:   v.GetInputType() == TextInputTypeMultiline,
			ReadOnly:    v.IsReadOnly(),
		}
	case *ProgressBar:
		return ControlProps{
			Value:         float64(v.GetValue()),
			Min:           float64(v.GetMin()),
			Max:           float64(v.GetMax()),
			Indeterminate: v.IsIndeterminate(),
			ShowText:      v.IsShowingText(),
		}
	case *Switch:
		return ControlProps{Checked: v.IsChecked(), Text: v.GetLabel()}
	case *RadioButton:
		return ControlProps{Checked: v.IsChecked(), Text: v.GetLabel(), Group: v.GetGroup()}
	case *List:
		return ControlProps{
			Items:         v.GetItems(),
			Selected:      v.SelectedIndex(),
			SelectionMode: int(v.GetSelectionMode()),
		}
	case *Menu:
		return ControlProps{Text: v.GetLabel(), MenuItems: v.GetItems()}
	default:
		return ControlProps{}
	}
}

// widgetBase returns the embedded *BaseWidget of a known widget type, or nil.
func widgetBase(w Widget) *BaseWidget {
	switch v := w.(type) {
	case *Container:
		return v.BaseWidget
	case *Button:
		return v.BaseWidget
	case *Text:
		return v.BaseWidget
	case *Checkbox:
		return v.BaseWidget
	case *ComboBox:
		return v.BaseWidget
	case *Slider:
		return v.BaseWidget
	case *Image:
		return v.BaseWidget
	case *TextInput:
		return v.BaseWidget
	case *ProgressBar:
		return v.BaseWidget
	case *Switch:
		return v.BaseWidget
	case *RadioButton:
		return v.BaseWidget
	case *List:
		return v.BaseWidget
	case *ScrollView:
		return v.BaseWidget
	case *Menu:
		return v.BaseWidget
	default:
		return nil
	}
}

// syncWidget pushes a widget's current state to its already-mounted native
// control. It is a no-op for widgets that have not been mounted yet.
func syncWidget(w Widget) {
	bw := widgetBase(w)
	if bw == nil {
		return
	}
	bw.mu.RLock()
	native := bw.native
	bw.mu.RUnlock()
	if native == nil {
		return
	}
	native.Update(propsOf(w))
	native.SetVisible(w.GetVisible())
	native.SetEnabled(w.GetEnabled())
	native.SetStyle(w.GetStyle())
}

// layoutAndMount runs the framework's flex layout engine over the content tree
// and reflects the result onto the native controls rooted at surface. The flex
// engine remains authoritative on every platform, which keeps layout identical
// across backends: the backend only translates already-computed rectangles.
func layoutAndMount(backend Backend, surface NativeControl, content Widget, width, height float32) error {
	if content == nil {
		return nil
	}
	// TAB traversal walks the tree that was mounted last.
	SetFocusRoot(content)
	if width <= 0 {
		width = 800
	}
	if height <= 0 {
		height = 600
	}

	// Measure against the available space (max only), then lay the root out at
	// the full window size. Passing definite min/max here would clamp every
	// child to the window size instead of its natural size.
	constraints := flex.Constraint{
		MaxWidth:  width,
		MaxHeight: height,
	}
	content.Measure(constraints)
	content.Layout(0, 0, width, height)

	if err := mountRecursive(backend, surface, content, 0, 0); err != nil {
		return err
	}
	// The toolkit is now initialized and the first window laid out.
	fireReady()
	return nil
}

// mountRecursive creates (once) and positions the native control for w inside
// parent, then recurses into children. Coordinates produced by the layout
// engine are absolute to the window, so they are converted to parent-relative
// coordinates by subtracting the parent's origin.
func mountRecursive(backend Backend, parent NativeControl, w Widget, parentX, parentY float32) error {
	if w == nil {
		return nil
	}

	bw := widgetBase(w)
	if bw == nil {
		return nil
	}

	bw.mu.RLock()
	ctrl := bw.native
	bw.mu.RUnlock()

	if ctrl == nil {
		created, err := backend.CreateControl(w, kindOf(w), propsOf(w))
		if err != nil {
			return err
		}
		if created == nil {
			return nil
		}
		bw.SetNativeControl(created)
		ctrl = created
	}

	ctrl.AttachTo(parent)

	rect := w.GetRect()
	ctrl.SetBounds(
		rect.Position.X-parentX,
		rect.Position.Y-parentY,
		rect.Size.Width,
		rect.Size.Height,
	)
	ctrl.SetVisible(w.GetVisible())
	ctrl.SetEnabled(w.GetEnabled())
	ctrl.SetStyle(w.GetStyle())
	ctrl.Update(propsOf(w))

	for _, child := range w.GetChildren() {
		if err := mountRecursive(backend, ctrl, child, rect.Position.X, rect.Position.Y); err != nil {
			return err
		}
	}

	return nil
}

// renderWidget mounts (or refreshes) a single widget. It exists so that the
// Widget.Render method can stay meaningful without the caller having to know
// about the backend.
func renderWidget(w Widget) error {
	bw := widgetBase(w)
	if bw == nil {
		return nil
	}
	bw.mu.RLock()
	mounted := bw.native != nil
	bw.mu.RUnlock()
	if mounted {
		syncWidget(w)
		return nil
	}
	// Not mounted yet (no window/backend surface). Rendering is driven by the
	// window when the widget tree is attached, so this is a successful no-op.
	return nil
}
