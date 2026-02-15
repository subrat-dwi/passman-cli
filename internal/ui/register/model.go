package register

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/subrat-dwi/passman-cli/internal/app"
	"github.com/subrat-dwi/passman-cli/internal/validation"
)

const (
	fieldEmail = iota
	fieldPassword
	fieldConfirm
)

type RegisterModel struct {
	inputs     []textinput.Model
	focus      int
	app        *app.App
	err        string
	fieldErr   []string
	loading    bool
	pwStrength string
	success    bool
}

// Success returns whether registration completed successfully
func (m RegisterModel) Success() bool {
	return m.success
}

func NewRegisterModel(app *app.App) RegisterModel {
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

	confirm := textinput.New()
	confirm.Placeholder = "Confirm Password"
	confirm.CharLimit = 128
	confirm.Width = 40
	confirm.EchoMode = textinput.EchoPassword
	confirm.EchoCharacter = '•'
	confirm.Prompt = "  "

	return RegisterModel{
		inputs:   []textinput.Model{email, password, confirm},
		focus:    fieldEmail,
		app:      app,
		fieldErr: make([]string, 3),
	}
}

func (m RegisterModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *RegisterModel) validateField(field int) {
	switch field {
	case fieldEmail:
		if err := validation.ValidateEmail(m.inputs[fieldEmail].Value()); err != nil {
			m.fieldErr[fieldEmail] = err.Error()
		} else {
			m.fieldErr[fieldEmail] = ""
		}
	case fieldPassword:
		pw := m.inputs[fieldPassword].Value()
		if err := validation.ValidateMasterPassword(pw); err != nil {
			m.fieldErr[fieldPassword] = err.Error()
			m.pwStrength = ""
		} else {
			m.fieldErr[fieldPassword] = ""
			_, m.pwStrength = validation.GetPasswordStrength(pw)
		}
		// Also re-validate confirm if it has value
		if m.inputs[fieldConfirm].Value() != "" {
			m.validateField(fieldConfirm)
		}
	case fieldConfirm:
		confirm := m.inputs[fieldConfirm].Value()
		if confirm == "" {
			m.fieldErr[fieldConfirm] = "please confirm your password"
		} else if confirm != m.inputs[fieldPassword].Value() {
			m.fieldErr[fieldConfirm] = "passwords do not match"
		} else {
			m.fieldErr[fieldConfirm] = ""
		}
	}
}

func (m *RegisterModel) validateAll() bool {
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

func (m *RegisterModel) nextField() {
	m.validateField(m.focus)
	m.inputs[m.focus].Blur()
	m.focus = (m.focus + 1) % len(m.inputs)
	m.inputs[m.focus].Focus()
}

func (m *RegisterModel) prevField() {
	m.validateField(m.focus)
	m.inputs[m.focus].Blur()
	m.focus = (m.focus + len(m.inputs) - 1) % len(m.inputs)
	m.inputs[m.focus].Focus()
}
