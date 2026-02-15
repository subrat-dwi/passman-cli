package list

import (
	"fmt"
	"strings"
)

func (m listModel) View() string {
	if m.quitting {
		return ""
	}

	if m.err != "" {
		return fmt.Sprintf("Error: %s\n\nPress ESC to go back", m.err)
	}

	switch m.state {
	case stateDetail:
		return m.viewDetail()
	case stateEdit:
		return m.viewEdit()
	case stateDeleteConfirm:
		return m.viewDeleteConfirm()
	default:
		return m.viewList()
	}
}

func (m listModel) viewList() string {
	var b strings.Builder
	if m.statusMsg != "" {
		b.WriteString(fmt.Sprintf("  ✓ %s\n\n", m.statusMsg))
	}
	b.WriteString(m.list.View())
	b.WriteString("\n  [enter] view  [u] edit  [d] delete  [q] quit")
	return b.String()
}

func (m listModel) viewDetail() string {
	if m.selected == nil {
		return ""
	}

	var b strings.Builder
	if m.statusMsg != "" {
		b.WriteString(fmt.Sprintf("  ✓ %s\n\n", m.statusMsg))
	}
	b.WriteString("╭──────────────────────────────────────╮\n")
	b.WriteString("│          Password Details            │\n")
	b.WriteString("╰──────────────────────────────────────╯\n\n")
	b.WriteString(fmt.Sprintf("  Service:  %s\n", m.selected.Name))
	b.WriteString(fmt.Sprintf("  Username: %s\n", m.selected.Username))
	b.WriteString(fmt.Sprintf("  Password: %s\n", m.selected.Password))
	b.WriteString("\n  [c] copy  [u] edit  [d] delete  [esc] back  [q] quit")
	return b.String()
}

func (m listModel) viewEdit() string {
	var b strings.Builder
	b.WriteString("╭──────────────────────────────────────╮\n")
	b.WriteString("│           Edit Password              │\n")
	b.WriteString("╰──────────────────────────────────────╯\n\n")

	for i, input := range m.editInputs {
		label := ""
		switch editField(i) {
		case fieldName:
			label = "Service"
		case fieldUsername:
			label = "Username"
		case fieldPassword:
			label = "Password"
		}

		cursor := " "
		if editField(i) == m.editFocus {
			cursor = ">"
		}
		b.WriteString(fmt.Sprintf(" %s %-10s %s\n", cursor, label+":", input.View()))
	}

	b.WriteString("\n  [tab/↑↓] navigate  [enter] save  [esc] cancel")
	return b.String()
}

func (m listModel) viewDeleteConfirm() string {
	var b strings.Builder
	b.WriteString("╭──────────────────────────────────────╮\n")
	b.WriteString("│         Delete Password?             │\n")
	b.WriteString("╰──────────────────────────────────────╯\n\n")
	b.WriteString("  Are you sure you want to delete this password?\n")
	b.WriteString("  This action cannot be undone.\n\n")
	b.WriteString("  [y] Yes, delete  [n] No, cancel")
	return b.String()
}
