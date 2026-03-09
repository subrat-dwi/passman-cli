package list

import (
	"fmt"
	"strings"

	"github.com/subrat-dwi/passman-cli/internal/ui/styles"
)

func (m listModel) View() string {
	if m.quitting {
		return ""
	}

	if m.err != "" {
		return fmt.Sprintf("%s\n\nPress ESC to go back", styles.ErrorMsg(m.err))
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
		b.WriteString(fmt.Sprintf("  %s\n\n", styles.SuccessMsg(m.statusMsg)))
	}
	b.WriteString(m.list.View())
	b.WriteString("\n" + styles.Dim.Render("  [enter] view  [u] edit  [d] delete  [q] quit"))
	return b.String()
}

func (m listModel) viewDetail() string {
	if m.selected == nil {
		return ""
	}

	var b strings.Builder
	if m.statusMsg != "" {
		b.WriteString(fmt.Sprintf("  %s\n\n", styles.SuccessMsg(m.statusMsg)))
	}
	b.WriteString(styles.Box("Password Details", 38))
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("  %s  %s\n", styles.Dim.Render("Service:"), m.selected.Name))
	b.WriteString(fmt.Sprintf("  %s %s\n", styles.Dim.Render("Username:"), m.selected.Username))
	b.WriteString(fmt.Sprintf("  %s %s\n", styles.Dim.Render("Password:"), m.selected.Password))
	b.WriteString("\n" + styles.Dim.Render("  [c] copy  [u] edit  [d] delete  [esc] back  [q] quit"))
	return b.String()
}

func (m listModel) viewEdit() string {
	var b strings.Builder
	b.WriteString(styles.Box("Edit Password", 38))
	b.WriteString("\n\n")

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
			cursor = styles.Cursor()
		}
		b.WriteString(fmt.Sprintf(" %s %-10s %s\n", cursor, label+":", input.View()))
	}

	b.WriteString("\n" + styles.Dim.Render("  [tab/↑↓] navigate  [enter] save  [esc] cancel"))
	return b.String()
}

func (m listModel) viewDeleteConfirm() string {
	var b strings.Builder
	b.WriteString(styles.Box("Delete Password?", 38))
	b.WriteString("\n\n")
	b.WriteString("  Are you sure you want to delete this password?\n")
	b.WriteString(fmt.Sprintf("  %s\n\n", styles.Warning.Render("This action cannot be undone.")))
	b.WriteString(styles.Dim.Render("  [y] Yes, delete  [n] No, cancel"))
	return b.String()
}
