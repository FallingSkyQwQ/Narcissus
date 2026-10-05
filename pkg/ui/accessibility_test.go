package ui

import (
	"testing"

	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

func TestAccessibleDefaultsFromWidgetKind(t *testing.T) {
	button := NewButton().Text("Save")
	check := NewCheckbox().Label("Remember me")
	label := NewText("Title")
	input := NewTextInput().Placeholder("Email")

	cases := []struct {
		w    Widget
		role AccessibleRole
		name string
	}{
		{button, RoleButton, "Save"},
		{check, RoleCheckbox, "Remember me"},
		{label, RoleLabel, "Title"},
		{input, RoleTextInput, "Email"},
	}

	for _, tc := range cases {
		props := propsOf(tc.w)
		if props.AccessibleRole != tc.role {
			t.Errorf("%T role = %v, want %v", tc.w, props.AccessibleRole, tc.role)
		}
		if props.AccessibleName != tc.name {
			t.Errorf("%T name = %q, want %q", tc.w, props.AccessibleName, tc.name)
		}
	}
}

func TestAccessibleExplicitOverridesDefaults(t *testing.T) {
	button := NewButton().Text("OK")
	button.SetAccessibility(RoleGroup, "Confirm", "Applies the changes")
	button.SetAccessibleRole(RoleButton)

	props := propsOf(button)
	if props.AccessibleRole != RoleButton {
		t.Errorf("role = %v, want RoleButton", props.AccessibleRole)
	}
	if props.AccessibleName != "Confirm" {
		t.Errorf("name = %q, want Confirm", props.AccessibleName)
	}
	if props.AccessibleDescription != "Applies the changes" {
		t.Errorf("description = %q, want %q", props.AccessibleDescription, "Applies the changes")
	}
}

func TestAccessibilityTreeMirrorsWidgetTree(t *testing.T) {
	title := NewText("Sign in")
	name := NewTextInput().Placeholder("Name")
	ok := NewButton().Text("Continue")
	root := NewContainer().Direction(flex.DirectionColumn).Add(title, name, ok)

	node := AccessibilityTree(root)
	if node == nil {
		t.Fatal("AccessibilityTree returned nil")
	}
	if node.Role != RoleGroup {
		t.Errorf("root role = %v, want RoleGroup", node.Role)
	}
	if len(node.Children) != 3 {
		t.Fatalf("root children = %d, want 3", len(node.Children))
	}
	if node.Children[0].Role != RoleLabel || node.Children[0].Name != "Sign in" {
		t.Errorf("first child = %v/%q, want label/Sign in", node.Children[0].Role, node.Children[0].Name)
	}
	if node.Children[2].Role != RoleButton || node.Children[2].Name != "Continue" {
		t.Errorf("button node = %v/%q, want button/Continue", node.Children[2].Role, node.Children[2].Name)
	}
}

func TestAccessibilityTreeSkipsHiddenWidgets(t *testing.T) {
	visible := NewButton().Text("visible")
	hidden := NewButton().Text("hidden")
	hidden.SetVisible(false)
	root := NewContainer().Add(visible, hidden)

	node := AccessibilityTree(root)
	if len(node.Children) != 1 {
		t.Fatalf("children = %d, want 1 (hidden widget excluded)", len(node.Children))
	}
	if node.Children[0].Name != "visible" {
		t.Errorf("child name = %q, want visible", node.Children[0].Name)
	}
}

func TestAccessibilityTreeNilForNilRoot(t *testing.T) {
	if AccessibilityTree(nil) != nil {
		t.Error("expected nil tree for nil root")
	}
}

func TestMountPublishesAccessibilityProps(t *testing.T) {
	button := NewButton().Text("Submit")
	button.SetAccessibleDescription("Submits the form")
	root := NewContainer().Add(button)

	backend := &fakeBackend{}
	if err := layoutAndMount(backend, &fakeControl{kind: ControlContainer}, root, 200, 100); err != nil {
		t.Fatalf("layoutAndMount: %v", err)
	}

	var found bool
	for _, c := range backend.created {
		if c.kind != ControlButton {
			continue
		}
		found = true
		if c.props.AccessibleName != "Submit" {
			t.Errorf("props name = %q, want Submit", c.props.AccessibleName)
		}
		if c.props.AccessibleDescription != "Submits the form" {
			t.Errorf("props description = %q, want %q", c.props.AccessibleDescription, "Submits the form")
		}
		if c.props.AccessibleRole != RoleButton {
			t.Errorf("props role = %v, want RoleButton", c.props.AccessibleRole)
		}
	}
	if !found {
		t.Fatal("no button control created")
	}
}

func TestAccessibleRoleString(t *testing.T) {
	if RoleButton.String() != "button" {
		t.Errorf("RoleButton.String() = %q", RoleButton.String())
	}
	if RoleNone.String() != "none" {
		t.Errorf("RoleNone.String() = %q", RoleNone.String())
	}
}
