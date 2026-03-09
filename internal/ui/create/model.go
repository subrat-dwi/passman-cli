package create

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/subrat-dwi/passman-cli/internal/app"
	"github.com/subrat-dwi/passman-cli/internal/validation"
)

const (
	fieldName = iota
	fieldUsername
	fieldPassword
)

type createPasswordModel struct {
	inputs   []textinput.Model
	focus    int
	err      string
	fieldErr []string
	loading  bool
	success  bool
	app      *app.App
}

func NewCreatePasswordModel(app *app.App) createPasswordModel {
	name := textinput.New()
	name.Placeholder = "Service (1-64 chars)"
	name.Focus()
	name.CharLimit = 64
	name.Width = 40
	name.Prompt = "  "

	username := textinput.New()
	username.Placeholder = "Username"
	username.CharLimit = 128
	username.Width = 40
	username.Prompt = "  "

	password := textinput.New()
	password.Placeholder = "Password"
	password.CharLimit = 256
	password.Width = 40
	password.Prompt = "  "
	password.EchoMode = textinput.EchoPassword

	return createPasswordModel{
		inputs:   []textinput.Model{name, username, password},
		focus:    fieldName,
		app:      app,
		fieldErr: make([]string, 3),
	}
}

func (m createPasswordModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *createPasswordModel) validateField(field int) {
	switch field {
	case fieldName:
		if err := validation.ValidateServiceName(m.inputs[fieldName].Value()); err != nil {
			m.fieldErr[fieldName] = err.Error()
		} else {
			m.fieldErr[fieldName] = ""
		}
	case fieldUsername:
		if err := validation.ValidateUsername(m.inputs[fieldUsername].Value()); err != nil {
			m.fieldErr[fieldUsername] = err.Error()
		} else {
			m.fieldErr[fieldUsername] = ""
		}
	case fieldPassword:
		if err := validation.ValidatePassword(m.inputs[fieldPassword].Value()); err != nil {
			m.fieldErr[fieldPassword] = err.Error()
		} else {
			m.fieldErr[fieldPassword] = ""
		}
	}
}

func (m *createPasswordModel) validateAll() bool {
	for i := range m.inputs {
		m.validateField(i)
	}
	for _, e := range m.fieldErr {
		if e != "" {
			return false
		}
	}
	return true
}

func (m *createPasswordModel) nextField() {
	m.validateField(m.focus)
	m.inputs[m.focus].Blur()
	m.focus = (m.focus + 1) % len(m.inputs)
	m.inputs[m.focus].Focus()
}

func (m *createPasswordModel) prevField() {
	m.validateField(m.focus)
	m.inputs[m.focus].Blur()
	m.focus = (m.focus + len(m.inputs) - 1) % len(m.inputs)
	m.inputs[m.focus].Focus()
}
