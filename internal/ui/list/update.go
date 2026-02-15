package list

import (
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/subrat-dwi/passman-cli/internal/passwordmanager"
	"github.com/subrat-dwi/passman-cli/internal/service"
)

type fetchPasswordMsg struct {
	password *service.DecryptedPassword
	err      error
}

type updateResultMsg struct {
	err error
}

type deleteResultMsg struct {
	deletedID string
	err       error
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Pass window size to list for proper rendering
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)
		return m, nil

	case tea.KeyMsg:
		// Handle quit everywhere
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

		switch m.state {
		case stateList:
			return m.updateList(msg)
		case stateDetail:
			return m.updateDetail(msg)
		case stateEdit:
			return m.updateEdit(msg)
		case stateDeleteConfirm:
			return m.updateDeleteConfirm(msg)
		}

	case fetchPasswordMsg:
		if msg.err != nil {
			m.err = msg.err.Error()
			return m, nil
		}
		m.selected = msg.password
		m.state = stateDetail
		return m, nil

	case fetchForEditMsg:
		if msg.err != nil {
			m.err = msg.err.Error()
			return m, nil
		}
		m.selected = msg.password
		m.enterEditMode()
		return m, nil

	case updateResultMsg:
		if msg.err != nil {
			m.err = msg.err.Error()
			m.state = stateDetail
			return m, nil
		}
		m.statusMsg = "Password updated successfully"
		// Refresh list and return to list view
		return m, m.refreshList()

	case deleteResultMsg:
		if msg.err != nil {
			m.err = msg.err.Error()
			m.state = stateList
			return m, nil
		}
		m.statusMsg = "Password deleted successfully"
		// Refresh list
		return m, m.refreshList()

	case refreshListMsg:
		if msg.err != nil {
			m.err = msg.err.Error()
			return m, nil
		}
		// Update list items
		items := make([]list.Item, len(msg.entries))
		for i, e := range msg.entries {
			items[i] = passwordItem{entry: e}
		}
		m.list.SetItems(items)
		m.state = stateList
		m.selected = nil
		return m, nil

	default:
		// Pass other messages (like cursor blink) to the list when in list state
		if m.state == stateList {
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m listModel) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// If filtering, pass all keys to the list component
	if m.list.FilterState() == list.Filtering {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	switch msg.String() {
	case "q":
		m.quitting = true
		return m, tea.Quit

	case "enter":
		selected, ok := m.list.SelectedItem().(passwordItem)
		if !ok {
			return m, nil
		}
		return m, m.fetchPassword(selected.entry.ID)

	case "d":
		selected, ok := m.list.SelectedItem().(passwordItem)
		if !ok {
			return m, nil
		}
		m.deleteTargetID = selected.entry.ID
		m.state = stateDeleteConfirm
		return m, nil

	case "u":
		selected, ok := m.list.SelectedItem().(passwordItem)
		if !ok {
			return m, nil
		}
		// First fetch password details, then switch to edit mode
		return m, m.fetchPasswordForEdit(selected.entry.ID)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m listModel) updateDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		m.quitting = true
		return m, tea.Quit

	case "esc":
		m.state = stateList
		m.selected = nil
		m.err = ""
		return m, nil

	case "c":
		if m.selected != nil {
			copyWithAutoClear(m.selected.Password, 60*time.Second)
			m.statusMsg = "Password copied to clipboard (clears in 60s)"
		}
		return m, nil

	case "d":
		if m.selected != nil {
			m.deleteTargetID = m.selected.ID
			m.state = stateDeleteConfirm
		}
		return m, nil

	case "u":
		if m.selected != nil {
			m.enterEditMode()
		}
		return m, nil
	}

	return m, nil
}

func (m *listModel) enterEditMode() {
	m.editInputs[fieldName].SetValue(m.selected.Name)
	m.editInputs[fieldUsername].SetValue(m.selected.Username)
	m.editInputs[fieldPassword].SetValue(m.selected.Password)
	m.editFocus = fieldName
	m.editInputs[fieldName].Focus()
	m.editInputs[fieldUsername].Blur()
	m.editInputs[fieldPassword].Blur()
	m.state = stateEdit
}

func (m listModel) updateEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = stateDetail
		return m, nil

	case "tab", "down":
		m.editInputs[m.editFocus].Blur()
		m.editFocus = (m.editFocus + 1) % 3
		m.editInputs[m.editFocus].Focus()
		return m, nil

	case "shift+tab", "up":
		m.editInputs[m.editFocus].Blur()
		m.editFocus = (m.editFocus + 2) % 3 // -1 mod 3
		m.editInputs[m.editFocus].Focus()
		return m, nil

	case "enter":
		// Save changes
		name := m.editInputs[fieldName].Value()
		username := m.editInputs[fieldUsername].Value()
		password := m.editInputs[fieldPassword].Value()

		if name == "" || username == "" || password == "" {
			m.err = "All fields are required"
			return m, nil
		}

		return m, m.saveUpdate(m.selected.ID, name, username, password)
	}

	// Update the focused input
	var cmd tea.Cmd
	m.editInputs[m.editFocus], cmd = m.editInputs[m.editFocus].Update(msg)
	return m, cmd
}

func (m listModel) updateDeleteConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		return m, m.deletePassword(m.deleteTargetID)

	case "n", "N", "esc":
		m.deleteTargetID = ""
		if m.selected != nil {
			m.state = stateDetail
		} else {
			m.state = stateList
		}
		return m, nil
	}

	return m, nil
}

// Commands

func (m listModel) fetchPassword(id string) tea.Cmd {
	return func() tea.Msg {
		password, err := m.app.PasswordService.Get(id)
		return fetchPasswordMsg{password: password, err: err}
	}
}

type fetchForEditMsg struct {
	password *service.DecryptedPassword
	err      error
}

func (m listModel) fetchPasswordForEdit(id string) tea.Cmd {
	return func() tea.Msg {
		password, err := m.app.PasswordService.Get(id)
		if err != nil {
			return fetchPasswordMsg{password: nil, err: err}
		}
		return fetchForEditMsg{password: password, err: nil}
	}
}

func (m listModel) saveUpdate(id, name, username, password string) tea.Cmd {
	return func() tea.Msg {
		err := m.app.PasswordService.Update(id, name, username, password)
		return updateResultMsg{err: err}
	}
}

func (m listModel) deletePassword(id string) tea.Cmd {
	return func() tea.Msg {
		err := m.app.PasswordService.Delete(id)
		return deleteResultMsg{deletedID: id, err: err}
	}
}

type refreshListMsg struct {
	entries []passwordmanager.PasswordEntry
	err     error
}

func (m listModel) refreshList() tea.Cmd {
	return func() tea.Msg {
		entries, err := m.app.PasswordService.List()
		return refreshListMsg{entries: entries, err: err}
	}
}

func copyWithAutoClear(text string, ttl time.Duration) {
	_ = clipboard.WriteAll(text)

	// Clear clipboard after TTL
	time.AfterFunc(ttl, func() {
		current, _ := clipboard.ReadAll()
		if current == text {
			_ = clipboard.WriteAll("")
		}
	})
}
