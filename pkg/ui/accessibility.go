package ui

// AccessibleRole describes what a widget is, so assistive technology can
// announce it correctly. The zero value (RoleNone) means "use the role implied
// by the widget kind".
type AccessibleRole int

const (
	// RoleNone leaves the role to be inferred from the widget kind.
	RoleNone AccessibleRole = iota
	// RoleGeneric is a container or unclassified element.
	RoleGeneric
	// RoleGroup groups related widgets.
	RoleGroup
	// RoleButton is a clickable button.
	RoleButton
	// RoleCheckbox is a checkable input.
	RoleCheckbox
	// RoleRadio is a radio button.
	RoleRadio
	// RoleSwitch is an on/off toggle.
	RoleSwitch
	// RoleComboBox is a drop-down selector.
	RoleComboBox
	// RoleSlider is a range control.
	RoleSlider
	// RoleTextInput is an editable text field.
	RoleTextInput
	// RoleLabel is static text.
	RoleLabel
	// RoleImage is a graphic.
	RoleImage
	// RoleProgress is a progress indicator.
	RoleProgress
	// RoleList is a list of items.
	RoleList
	// RoleMenu is a menu button.
	RoleMenu
	// RoleDialog is a modal dialog.
	RoleDialog
	// RoleAlert is an important, time-sensitive message such as a toast.
	RoleAlert
)

// String returns a stable, human readable name for an AccessibleRole.
func (r AccessibleRole) String() string {
	switch r {
	case RoleGeneric:
		return "generic"
	case RoleGroup:
		return "group"
	case RoleButton:
		return "button"
	case RoleCheckbox:
		return "checkbox"
	case RoleRadio:
		return "radio"
	case RoleSwitch:
		return "switch"
	case RoleComboBox:
		return "combobox"
	case RoleSlider:
		return "slider"
	case RoleTextInput:
		return "textinput"
	case RoleLabel:
		return "label"
	case RoleImage:
		return "image"
	case RoleProgress:
		return "progress"
	case RoleList:
		return "list"
	case RoleMenu:
		return "menu"
	case RoleDialog:
		return "dialog"
	case RoleAlert:
		return "alert"
	default:
		return "none"
	}
}

// AccessibleNode is a node in the accessibility tree. It mirrors the widget
// tree but carries only the semantics an assistive technology consumes: a
// role, an accessible name and an optional description.
type AccessibleNode struct {
	// Widget is the widget this node describes.
	Widget Widget
	// Role is the resolved role (never RoleNone).
	Role AccessibleRole
	// Name is the accessible name (explicit, else derived from content).
	Name string
	// Description is the accessible description.
	Description string
	// Children are the child nodes in document order.
	Children []*AccessibleNode
}

// AccessibilityTree builds the accessibility tree rooted at w. It is the
// toolkit-neutral view a backend or screen reader bridge can walk; it needs no
// display, so it can be inspected in tests.
func AccessibilityTree(w Widget) *AccessibleNode {
	if w == nil {
		return nil
	}
	return buildAccessibleNode(w)
}

func buildAccessibleNode(w Widget) *AccessibleNode {
	node := &AccessibleNode{
		Widget:      w,
		Role:        accessibleRoleFor(w),
		Name:        accessibleNameFor(w),
		Description: accessibleDescriptionFor(w),
	}
	for _, child := range w.GetChildren() {
		if child == nil || !child.GetVisible() {
			continue
		}
		node.Children = append(node.Children, buildAccessibleNode(child))
	}
	return node
}

// accessibleRoleFor resolves the role of a widget: an explicit role wins,
// otherwise the role implied by the widget kind.
func accessibleRoleFor(w Widget) AccessibleRole {
	if bw := widgetBase(w); bw != nil {
		if role := bw.AccessibleRole(); role != RoleNone {
			return role
		}
	}
	return defaultAccessibleRole(kindOf(w))
}

// accessibleNameFor resolves the accessible name: an explicit name wins,
// otherwise a sensible default derived from the widget's content.
func accessibleNameFor(w Widget) string {
	if bw := widgetBase(w); bw != nil {
		if name := bw.AccessibleName(); name != "" {
			return name
		}
	}
	return defaultAccessibleName(w)
}

func accessibleDescriptionFor(w Widget) string {
	if bw := widgetBase(w); bw != nil {
		return bw.AccessibleDescription()
	}
	return ""
}

// defaultAccessibleRole maps a control kind to the role assistive technology
// should announce when the widget does not specify one.
func defaultAccessibleRole(kind ControlKind) AccessibleRole {
	switch kind {
	case ControlButton:
		return RoleButton
	case ControlText:
		return RoleLabel
	case ControlCheckbox:
		return RoleCheckbox
	case ControlComboBox:
		return RoleComboBox
	case ControlSlider:
		return RoleSlider
	case ControlImage:
		return RoleImage
	case ControlTextInput:
		return RoleTextInput
	case ControlProgress:
		return RoleProgress
	case ControlSwitch:
		return RoleSwitch
	case ControlRadio:
		return RoleRadio
	case ControlList:
		return RoleList
	case ControlScroll:
		return RoleGroup
	case ControlMenu:
		return RoleMenu
	case ControlDialog:
		return RoleDialog
	case ControlToast:
		return RoleAlert
	case ControlContainer:
		return RoleGroup
	default:
		return RoleGeneric
	}
}

// defaultAccessibleName derives an accessible name from the widget's own
// content so common controls are labelled without extra configuration.
func defaultAccessibleName(w Widget) string {
	switch v := w.(type) {
	case *Button:
		return v.GetText()
	case *Checkbox:
		return v.GetLabel()
	case *RadioButton:
		return v.GetLabel()
	case *Switch:
		return v.GetLabel()
	case *Text:
		return v.GetText()
	case *TextInput:
		return v.GetPlaceholder()
	case *ComboBox:
		return v.GetPlaceholder()
	case *Menu:
		return v.GetLabel()
	default:
		return ""
	}
}
