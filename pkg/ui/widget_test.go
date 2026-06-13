package ui

import (
	"testing"

	"github.com/FallingSkyQwQ/Narcissus/pkg/layout/flex"
)

func TestNewContainer(t *testing.T) {
	c := NewContainer()
	if c == nil {
		t.Fatal("NewContainer() returned nil")
	}

	// 测试默认属性
	if c.direction != flex.DirectionRow {
		t.Errorf("expected default direction to be DirectionRow, got %v", c.direction)
	}
	if c.wrap != flex.WrapNoWrap {
		t.Errorf("expected default wrap to be WrapNoWrap, got %v", c.wrap)
	}
}

func TestContainerChaining(t *testing.T) {
	c := NewContainer().
		Direction(flex.DirectionColumn).
		Justify(flex.JustifyCenter).
		Align(flex.AlignCenter).
		Gap(10).
		FlexGrow(1).
		Padding(UniformInsets(16))

	if c.direction != flex.DirectionColumn {
		t.Errorf("expected direction to be DirectionColumn, got %v", c.direction)
	}
	if c.flexContainer.Justify != flex.JustifyCenter {
		t.Errorf("expected justify to be JustifyCenter, got %v", c.flexContainer.Justify)
	}
	if c.style.Padding.Top != 16 {
		t.Errorf("expected padding top to be 16, got %f", c.style.Padding.Top)
	}
	if c.style.FlexGrow != 1 {
		t.Errorf("expected flex grow to be 1, got %f", c.style.FlexGrow)
	}
}

func TestContainerAddChild(t *testing.T) {
	c := NewContainer()
	text := NewText("Hello")
	button := NewButton().Text("Click")

	c.Add(text, button)

	children := c.GetChildren()
	if len(children) != 2 {
		t.Errorf("expected 2 children, got %d", len(children))
	}
}

func TestNewText(t *testing.T) {
	text := NewText("Hello World")
	if text == nil {
		t.Fatal("NewText() returned nil")
	}
	if text.GetText() != "Hello World" {
		t.Errorf("expected text to be 'Hello World', got '%s'", text.GetText())
	}
}

func TestTextChaining(t *testing.T) {
	text := NewText("Initial").
		Text("Updated").
		FontSize(20).
		TextColor(ColorRed).
		FlexGrow(1)

	if text.GetText() != "Updated" {
		t.Errorf("expected text to be 'Updated', got '%s'", text.GetText())
	}
	if text.style.Font.Size != 20 {
		t.Errorf("expected font size to be 20, got %f", text.style.Font.Size)
	}
	if text.style.TextColor != ColorRed {
		t.Errorf("expected text color to be Red, got %v", text.style.TextColor)
	}
}

func TestNewButton(t *testing.T) {
	btn := NewButton()
	if btn == nil {
		t.Fatal("NewButton() returned nil")
	}

	// 测试默认样式
	if btn.style.BackgroundColor != ColorPrimary {
		t.Errorf("expected default background color to be ColorPrimary, got %v", btn.style.BackgroundColor)
	}
}

func TestButtonChaining(t *testing.T) {
	clicked := false
	btn := NewButton().
		Text("OK").
		OnClick(func() {
			clicked = true
		}).
		BackgroundColor(ColorGreen).
		BorderRadius(8)

	if btn.GetText() != "OK" {
		t.Errorf("expected text to be 'OK', got '%s'", btn.GetText())
	}
	if btn.style.Border.Radius != 8 {
		t.Errorf("expected border radius to be 8, got %f", btn.style.Border.Radius)
	}

	// 测试点击
	btn.Click()
	if !clicked {
		t.Error("expected OnClick handler to be called")
	}
}

func TestNewImage(t *testing.T) {
	img := NewImage()
	if img == nil {
		t.Fatal("NewImage() returned nil")
	}
	if img.GetFit() != ImageFitContain {
		t.Errorf("expected default fit to be ImageFitContain, got %v", img.GetFit())
	}
}

func TestImageChaining(t *testing.T) {
	img := NewImage().
		Source("test.png").
		Fit(ImageFitCover).
		Size(100, 200)

	if img.GetSource() != "test.png" {
		t.Errorf("expected source to be 'test.png', got '%s'", img.GetSource())
	}
	if img.GetFit() != ImageFitCover {
		t.Errorf("expected fit to be ImageFitCover, got %v", img.GetFit())
	}
	if img.style.Width != 100 || img.style.Height != 200 {
		t.Errorf("expected size to be 100x200, got %fx%f", img.style.Width, img.style.Height)
	}
}

func TestNewTextInput(t *testing.T) {
	input := NewTextInput()
	if input == nil {
		t.Fatal("NewTextInput() returned nil")
	}
	if input.GetInputType() != TextInputTypeText {
		t.Errorf("expected default input type to be TextInputTypeText, got %v", input.GetInputType())
	}
}

func TestTextInputChaining(t *testing.T) {
	changed := false
	input := NewTextInput().
		Placeholder("Enter text...").
		Value("initial").
		OnChange(func(value string) {
			changed = true
		}).
		MaxLength(100)

	if input.GetPlaceholder() != "Enter text..." {
		t.Errorf("expected placeholder to be 'Enter text...', got '%s'", input.GetPlaceholder())
	}
	if input.GetValue() != "initial" {
		t.Errorf("expected value to be 'initial', got '%s'", input.GetValue())
	}
	if input.GetMaxLength() != 100 {
		t.Errorf("expected max length to be 100, got %d", input.GetMaxLength())
	}

	// 测试设置文本
	input.SetText("new value")
	if input.GetValue() != "new value" {
		t.Errorf("expected value to be 'new value', got '%s'", input.GetValue())
	}
	if !changed {
		t.Error("expected OnChange handler to be called")
	}
}

func TestNewCheckbox(t *testing.T) {
	cb := NewCheckbox()
	if cb == nil {
		t.Fatal("NewCheckbox() returned nil")
	}
	if cb.IsChecked() {
		t.Error("expected checkbox to be unchecked by default")
	}
}

func TestCheckboxChaining(t *testing.T) {
	changed := false
	cb := NewCheckbox().
		Label("Enable feature").
		Checked(true).
		OnChange(func(checked bool) {
			changed = true
		})

	if cb.GetLabel() != "Enable feature" {
		t.Errorf("expected label to be 'Enable feature', got '%s'", cb.GetLabel())
	}
	if !cb.IsChecked() {
		t.Error("expected checkbox to be checked")
	}

	// 测试切换
	cb.Toggle()
	if cb.IsChecked() {
		t.Error("expected checkbox to be unchecked after toggle")
	}
	if !changed {
		t.Error("expected OnChange handler to be called")
	}
}

func TestNewSlider(t *testing.T) {
	slider := NewSlider()
	if slider == nil {
		t.Fatal("NewSlider() returned nil")
	}
	if slider.GetMin() != 0 || slider.GetMax() != 100 {
		t.Errorf("expected default range to be 0-100, got %f-%f", slider.GetMin(), slider.GetMax())
	}
}

func TestSliderChaining(t *testing.T) {
	changed := false
	slider := NewSlider().
		Range(0, 10).
		Value(5).
		Step(1).
		OnChange(func(value float32) {
			changed = true
		})

	if slider.GetMin() != 0 || slider.GetMax() != 10 {
		t.Errorf("expected range to be 0-10, got %f-%f", slider.GetMin(), slider.GetMax())
	}
	if slider.GetValue() != 5 {
		t.Errorf("expected value to be 5, got %f", slider.GetValue())
	}

	// 测试设置值
	slider.Value(7)
	if slider.GetValue() != 7 {
		t.Errorf("expected value to be 7, got %f", slider.GetValue())
	}
	if !changed {
		t.Error("expected OnChange handler to be called")
	}
}

func TestNewComboBox(t *testing.T) {
	cb := NewComboBox()
	if cb == nil {
		t.Fatal("NewComboBox() returned nil")
	}
	if cb.GetSelectedIndex() != -1 {
		t.Errorf("expected default selected index to be -1, got %d", cb.GetSelectedIndex())
	}
}

func TestComboBoxChaining(t *testing.T) {
	changed := false
	cb := NewComboBox().
		Items("Option 1", "Option 2", "Option 3").
		SelectedIndex(0).
		Placeholder("Select an option").
		OnChange(func(index int, value string) {
			changed = true
		})

	items := cb.GetItems()
	if len(items) != 3 {
		t.Errorf("expected 3 items, got %d", len(items))
	}
	if cb.GetSelectedIndex() != 0 {
		t.Errorf("expected selected index to be 0, got %d", cb.GetSelectedIndex())
	}
	if cb.GetSelectedValue() != "Option 1" {
		t.Errorf("expected selected value to be 'Option 1', got '%s'", cb.GetSelectedValue())
	}

	// 测试选择
	cb.SelectItem(1)
	if cb.GetSelectedIndex() != 1 {
		t.Errorf("expected selected index to be 1, got %d", cb.GetSelectedIndex())
	}
	if !changed {
		t.Error("expected OnChange handler to be called")
	}
}

func TestStyleChaining(t *testing.T) {
	style := NewStyle().
		WithMargin(UniformInsets(10)).
		WithPadding(UniformInsets(8)).
		WithFlexGrow(1).
		WithBackgroundColor(ColorBlue).
		WithTextColor(ColorWhite).
		WithBorderRadius(4)

	if style.Margin.Top != 10 {
		t.Errorf("expected margin top to be 10, got %f", style.Margin.Top)
	}
	if style.Padding.Top != 8 {
		t.Errorf("expected padding top to be 8, got %f", style.Padding.Top)
	}
	if style.FlexGrow != 1 {
		t.Errorf("expected flex grow to be 1, got %f", style.FlexGrow)
	}
	if style.BackgroundColor != ColorBlue {
		t.Errorf("expected background color to be Blue, got %v", style.BackgroundColor)
	}
	if style.Border.Radius != 4 {
		t.Errorf("expected border radius to be 4, got %f", style.Border.Radius)
	}
}

func TestTheme(t *testing.T) {
	// 测试浅色主题
	lightTheme := NewLightTheme()
	if lightTheme.GetType() != ThemeLight {
		t.Errorf("expected theme type to be ThemeLight, got %v", lightTheme.GetType())
	}

	// 测试深色主题
	darkTheme := NewDarkTheme()
	if darkTheme.GetType() != ThemeDark {
		t.Errorf("expected theme type to be ThemeDark, got %v", darkTheme.GetType())
	}

	// 测试主题颜色
	bgColor := lightTheme.GetColor(TokenBackground)
	if bgColor != ColorLightBackground {
		t.Errorf("expected light theme background to be ColorLightBackground, got %v", bgColor)
	}
}

func TestEventRegistry(t *testing.T) {
	registry := NewEventRegistry()

	handled := false
	unsubscribe := registry.AddHandler(EventClick, func(e Event) {
		handled = true
	})

	// 创建测试事件
	event := &BaseEvent{
		Type:      EventClick,
		Target:    nil,
		Timestamp: 0,
	}

	// 分发事件
	result := registry.Dispatch(event)
	if !result {
		t.Error("expected event to be handled")
	}
	if !handled {
		t.Error("expected handler to be called")
	}

	// 取消订阅
	unsubscribe()

	// 再次分发事件
	handled = false
	result = registry.Dispatch(event)
	if result {
		t.Error("expected event not to be handled after unsubscribe")
	}
}

func TestComplexUI(t *testing.T) {
	// 创建一个复杂的 UI 布局
	root := NewContainer().
		Direction(flex.DirectionColumn).
		Padding(UniformInsets(16)).
		Gap(10)

	header := NewContainer().
		Direction(flex.DirectionRow).
		Justify(flex.JustifySpaceBetween)

	title := NewText("My App").FontSize(24).FontWeight(FontWeightBold)
	settingsBtn := NewButton().Text("Settings")

	header.Add(title, settingsBtn)

	content := NewContainer().
		Direction(flex.DirectionColumn).
		FlexGrow(1).
		Gap(8)

	input := NewTextInput().Placeholder("Enter your name").FlexGrow(1)
	checkbox := NewCheckbox().Label("Enable notifications")
	slider := NewSlider().Range(0, 100).Value(50)

	content.Add(input, checkbox, slider)

	footer := NewContainer().
		Direction(flex.DirectionRow).
		Justify(flex.JustifyFlexEnd).
		Gap(8)

	cancelBtn := NewButton().Text("Cancel").BackgroundColor(ColorGray)
	okBtn := NewButton().Text("OK")

	footer.Add(cancelBtn, okBtn)

	root.Add(header, content, footer)

	// 验证结构
	children := root.GetChildren()
	if len(children) != 3 {
		t.Errorf("expected 3 children in root, got %d", len(children))
	}

	headerChildren := header.GetChildren()
	if len(headerChildren) != 2 {
		t.Errorf("expected 2 children in header, got %d", len(headerChildren))
	}

	contentChildren := content.GetChildren()
	if len(contentChildren) != 3 {
		t.Errorf("expected 3 children in content, got %d", len(contentChildren))
	}
}
