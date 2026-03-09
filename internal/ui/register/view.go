package register

import (
	"fmt"
	"strings"

	"github.com/subrat-dwi/passman-cli/internal/ui/styles"
)

func (m RegisterModel) View() string {
	var b strings.Builder

	b.WriteString(styles.Box("Register to PassMan", 38))
	b.WriteString("\n\n")

	labels := []string{"Email", "Password", "Confirm"}

	for i, input := range m.inputs {
		cursor := " "
		if i == m.focus {
			cursor = styles.Cursor()
		}

		b.WriteString(fmt.Sprintf(" %s %-10s %s\n", cursor, labels[i]+":", input.View()))

		// Show field error
		if m.fieldErr[i] != "" {
			b.WriteString(fmt.Sprintf("              %s\n", styles.WarningMsg(m.fieldErr[i])))
		}

		// Show password strength indicator
		if i == fieldPassword && m.pwStrength != "" && m.fieldErr[i] == "" {
			var strengthMsg string
			switch m.pwStrength {
			case "Strong":
				strengthMsg = styles.SuccessMsg("Strength: " + m.pwStrength)
			case "Fair":
				strengthMsg = styles.InfoMsg("Strength: " + m.pwStrength)
			default:
				strengthMsg = styles.WarningMsg("Strength: " + m.pwStrength)
			}
			b.WriteString(fmt.Sprintf("              %s\n", strengthMsg))
		}
	}

	b.WriteString("\n")

	if m.loading {
		b.WriteString(fmt.Sprintf("  %s Registering...\n\n", styles.Highlight.Render("⏳")))
	}

	if m.err != "" {
		b.WriteString(fmt.Sprintf("  %s\n\n", styles.ErrorMsg(m.err)))
	}

	b.WriteString(styles.Dim.Render("  [tab/↑↓] switch field  [enter] submit  [esc] quit"))

	return b.String()
}
