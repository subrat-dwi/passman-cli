package login

import (
	tea "github.com/charmbracelet/bubbletea"
)

type loginSubmitMsg struct {
	email    string
	password string
}

type loginResultMsg struct {
	err error
}

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				return loginSubmitMsg{
					email:    m.inputs[fieldEmail].Value(),
					password: m.inputs[fieldPassword].Value(),
				}
			}
		}

	case loginSubmitMsg:
		// Perform login in a command to keep it non-blocking
		return m, func() tea.Msg {
			err := m.app.AuthService.Login(msg.email, msg.password)
			return loginResultMsg{err: err}
		}

	case loginResultMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err.Error()
			// Allow retry - focus on password field
			m.inputs[fieldPassword].SetValue("")
			m.inputs[m.focus].Blur()
			m.focus = fieldPassword
			m.inputs[m.focus].Focus()
			return m, nil
		}
		m.success = true
		return m, tea.Quit
	}

	// Update focused input and validate on change
	var cmd tea.Cmd
	m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
	cmds = append(cmds, cmd)

	// Live validation for password strength
	if m.focus == fieldPassword {
		m.validateField(fieldPassword)
	}

	return m, tea.Batch(cmds...)
}
