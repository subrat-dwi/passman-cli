package register

import (
	"fmt"
	"strings"
)

func (m RegisterModel) View() string {
	var b strings.Builder

	b.WriteString("╭──────────────────────────────────────╮\n")
	b.WriteString("│         Register to PassMan          │\n")
	b.WriteString("╰──────────────────────────────────────╯\n\n")

	labels := []string{"Email", "Password", "Confirm"}

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

		// Show password strength indicator
		if i == fieldPassword && m.pwStrength != "" && m.fieldErr[i] == "" {
			strengthColor := ""
			switch m.pwStrength {
			case "Strong":
				strengthColor = "✓"
			case "Fair":
				strengthColor = "○"
			default:
				strengthColor = "⚠"
			}
			b.WriteString(fmt.Sprintf("              %s Strength: %s\n", strengthColor, m.pwStrength))
		}
	}

	b.WriteString("\n")

	if m.loading {
		b.WriteString("  ⏳ Registering...\n\n")
	}

	if m.err != "" {
		b.WriteString(fmt.Sprintf("  ✗ Error: %s\n\n", m.err))
	}

	b.WriteString("  [tab/↑↓] switch field  [enter] submit  [esc] quit")

	return b.String()
}
