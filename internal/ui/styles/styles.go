package styles

import "github.com/charmbracelet/lipgloss"

// Colors based on bubbletea's default theme
// Uses AdaptiveColor for light/dark terminal support
var (
	// Primary accent - pink/magenta like bubbletea defaults
	Accent = lipgloss.AdaptiveColor{Light: "205", Dark: "205"}
	// Subtle text
	Subtle = lipgloss.AdaptiveColor{Light: "241", Dark: "241"}
	// Success green
	Green = lipgloss.AdaptiveColor{Light: "34", Dark: "76"}
	// Error red
	Red = lipgloss.AdaptiveColor{Light: "196", Dark: "203"}
	// Warning yellow
	Yellow = lipgloss.AdaptiveColor{Light: "214", Dark: "214"}
	// Info cyan
	Cyan = lipgloss.AdaptiveColor{Light: "39", Dark: "86"}
)

// Text styles
var (
	// Bold text
	Bold = lipgloss.NewStyle().Bold(true)

	// Dimmed/subtle text
	Dim = lipgloss.NewStyle().Foreground(Subtle)

	// Accent colored text (pink)
	Highlight = lipgloss.NewStyle().Foreground(Accent)

	// Success style (green with checkmark)
	Success = lipgloss.NewStyle().Foreground(Green)

	// Error style (red)
	Error = lipgloss.NewStyle().Foreground(Red)

	// Warning style (yellow)
	Warning = lipgloss.NewStyle().Foreground(Yellow)

	// Info style (cyan)
	Info = lipgloss.NewStyle().Foreground(Cyan)
)

// Box styles for headers
var (
	// Box border color
	BoxBorder = lipgloss.NewStyle().Foreground(Accent)

	// Box title
	BoxTitle = lipgloss.NewStyle().Bold(true).Foreground(Accent)
)

// Helper functions for common patterns

// SuccessMsg formats a success message with green checkmark
func SuccessMsg(msg string) string {
	return Success.Render("✓") + " " + msg
}

// ErrorMsg formats an error message with red X
func ErrorMsg(msg string) string {
	return Error.Render("✗") + " " + msg
}

// WarningMsg formats a warning message with yellow triangle
func WarningMsg(msg string) string {
	return Warning.Render("⚠") + " " + msg
}

// InfoMsg formats an info message with cyan arrow
func InfoMsg(msg string) string {
	return Info.Render("→") + " " + msg
}

// Cursor returns the selection cursor (highlighted >)
func Cursor() string {
	return Highlight.Render(">")
}

// Version formats version string
func Version(v string) string {
	return Bold.Render("pman") + " " + Highlight.Render(v)
}

// UserError formats a user error with message and hint
func UserError(message, hint string) string {
	if hint != "" {
		return Error.Render("✗") + " " + message + "\n  " + Dim.Render("→ "+hint)
	}
	return Error.Render("✗") + " " + message
}

// Box creates a simple bordered box header
func Box(title string, width int) string {
	if width < len(title)+2 {
		width = len(title) + 2
	}

	padding := (width - len(title)) / 2
	paddingRight := width - len(title) - padding

	top := BoxBorder.Render("╭") + BoxBorder.Render(repeat("─", width)) + BoxBorder.Render("╮")
	middle := BoxBorder.Render("│") + repeat(" ", padding) + BoxTitle.Render(title) + repeat(" ", paddingRight) + BoxBorder.Render("│")
	bottom := BoxBorder.Render("╰") + BoxBorder.Render(repeat("─", width)) + BoxBorder.Render("╯")

	return top + "\n" + middle + "\n" + bottom
}

func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
