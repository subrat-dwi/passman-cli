package login

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/subrat-dwi/passman-cli/internal/app"
	"github.com/subrat-dwi/passman-cli/internal/validation"
)

const (
	fieldEmail = iota
	fieldPassword
)

type LoginModel struct {
	inputs     []textinput.Model
	focus      int
	app        *app.App
	err        string
	fieldErr   []string // per-field validation errors
	loading    bool
	pwStrength string
	success    bool
}

// Success returns whether login completed successfully
func (m LoginModel) Success() bool {
	return m.success
}

func NewLoginModel(app *app.App) LoginModel {
	email := textinput.New()
	email.Placeholder = "Email"
	email.Focus()
	email.CharLimit = 254
	email.Width = 40
	email.Prompt = "  "

	password := textinput.New()
	password.Placeholder = "Master Password"
	password.CharLimit = 128
	password.Width = 40
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '•'
	password.Prompt = "  "

	return LoginModel{
		inputs:   []textinput.Model{email, password},
		focus:    fieldEmail,
		app:      app,
		fieldErr: make([]string, 2),
	}
}

func (m LoginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *LoginModel) validateField(field int) {
	switch field {
	case fieldEmail:
		if err := validation.ValidateEmail(m.inputs[fieldEmail].Value()); err != nil {
			m.fieldErr[fieldEmail] = err.Error()
		} else {
			m.fieldErr[fieldEmail] = ""
		}
	case fieldPassword:
		pw := m.inputs[fieldPassword].Value()
		if pw == "" {
			m.fieldErr[fieldPassword] = "password cannot be empty"
			m.pwStrength = ""
		} else {
			m.fieldErr[fieldPassword] = ""
			_, m.pwStrength = validation.GetPasswordStrength(pw)
		}
	}
}

func (m *LoginModel) validateAll() bool {
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

func (m *LoginModel) nextField() {
	m.validateField(m.focus)
	m.inputs[m.focus].Blur()
	m.focus = (m.focus + 1) % len(m.inputs)
	m.inputs[m.focus].Focus()
}

func (m *LoginModel) prevField() {
	m.validateField(m.focus)
	m.inputs[m.focus].Blur()
	m.focus = (m.focus + len(m.inputs) - 1) % len(m.inputs)
	m.inputs[m.focus].Focus()
}
