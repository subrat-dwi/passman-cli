package create

import tea "github.com/charmbracelet/bubbletea"

type createPasswordMsg struct {
	name     string
	username string
	password string
}

type createResultMsg struct {
	err error
}

func (m createPasswordModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Clear general error on any key press
		if m.err != "" {
			m.err = ""
		}

		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "tab", "down":
			m.nextField()
			return m, nil

		case "shift+tab", "up":
			m.prevField()
			return m, nil

		case "enter":
			if m.loading {
				return m, nil
			}

			// If not on last field, move to next
			if m.focus < len(m.inputs)-1 {
				m.nextField()
				return m, nil
			}

			// Validate all fields before submit
			if !m.validateAll() {
				return m, nil
			}

			// Submit
			m.loading = true
			return m, func() tea.Msg {
				return createPasswordMsg{
					name:     m.inputs[fieldName].Value(),
					username: m.inputs[fieldUsername].Value(),
					password: m.inputs[fieldPassword].Value(),
				}
			}
		}

	case createPasswordMsg:
		// Perform create in a command to keep it non-blocking
		return m, func() tea.Msg {
			err := m.app.PasswordService.Create(msg.name, msg.username, msg.password)
			return createResultMsg{err: err}
		}

	case createResultMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err.Error()
			// Allow retry
			return m, nil
		}
		m.success = true
		return m, tea.Quit
	}

	// Update focused input
	var cmd tea.Cmd
	m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
