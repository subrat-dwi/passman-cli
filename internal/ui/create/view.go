package create

import (
	"fmt"
	"strings"

	"github.com/subrat-dwi/passman-cli/internal/ui/styles"
)

func (m createPasswordModel) View() string {
	var b strings.Builder

	b.WriteString(styles.Box("Save New Password", 38))
	b.WriteString("\n\n")

	labels := []string{"Service", "Username", "Password"}

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
	}

	b.WriteString("\n")

	if m.loading {
		b.WriteString(fmt.Sprintf("  %s Saving password...\n\n", styles.Highlight.Render("⏳")))
	}

	if m.err != "" {
		b.WriteString(fmt.Sprintf("  %s\n\n", styles.ErrorMsg(m.err)))
	}

	b.WriteString(styles.Dim.Render("  [tab/↑↓] switch field  [enter] submit  [esc] quit"))

	return b.String()
}
