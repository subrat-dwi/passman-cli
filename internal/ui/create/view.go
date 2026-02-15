package create

import (
	"fmt"
	"strings"
)

func (m createPasswordModel) View() string {
	var b strings.Builder

	b.WriteString("╭──────────────────────────────────────╮\n")
	b.WriteString("│           Save New Password          │\n")
	b.WriteString("╰──────────────────────────────────────╯\n\n")

	labels := []string{"Service", "Username", "Password"}

	for i, input := range m.inputs {
		cursor := " "
		if i == m.focus {
			cursor = ">"
		}

		b.WriteString(fmt.Sprintf(" %s %-10s %s\n", cursor, labels[i]+":", input.View()))

		// Show field error
		if m.fieldErr[i] != "" {
			b.WriteString(fmt.Sprintf("              ⚠ %s\n", m.fieldErr[i]))
		}
	}

	b.WriteString("\n")

	if m.loading {
		b.WriteString("  ⏳ Saving password...\n\n")
	}

	if m.err != "" {
		b.WriteString(fmt.Sprintf("  ✗ Error: %s\n\n", m.err))
	}

	b.WriteString("  [tab/↑↓] switch field  [enter] submit  [esc] quit")

	return b.String()
}
