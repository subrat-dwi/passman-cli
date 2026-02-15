package list

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/subrat-dwi/passman-cli/internal/app"
	"github.com/subrat-dwi/passman-cli/internal/passwordmanager"
	"github.com/subrat-dwi/passman-cli/internal/service"
)

type passwordItem struct {
	entry passwordmanager.PasswordEntry
}

func (p passwordItem) Title() string       { return p.entry.Name }
func (p passwordItem) Description() string { return p.entry.Username }
func (p passwordItem) FilterValue() string { return p.entry.Name + " " + p.entry.Username }

type viewState int

const (
	stateList viewState = iota
	stateDetail
	stateEdit
	stateDeleteConfirm
)

type editField int

const (
	fieldName editField = iota
	fieldUsername
	fieldPassword
)

type listModel struct {
	list           list.Model
	app            *app.App
	state          viewState
	selected       *service.DecryptedPassword
	err            string
	quitting       bool
	editInputs     []textinput.Model
	editFocus      editField
	deleteTargetID string
	statusMsg      string
}

func NewListModel(app *app.App, entries []passwordmanager.PasswordEntry) listModel {
	items := make([]list.Item, len(entries))
	for i, e := range entries {
		items[i] = passwordItem{entry: e}
	}

	l := list.New(items, list.NewDefaultDelegate(), 50, 20)
	l.Title = "Your Passwords"
	l.SetShowHelp(true)

	// Initialize edit inputs
	inputs := make([]textinput.Model, 3)

	inputs[fieldName] = textinput.New()
	inputs[fieldName].Placeholder = "Service name"
	inputs[fieldName].CharLimit = 64
	inputs[fieldName].Width = 30

	inputs[fieldUsername] = textinput.New()
	inputs[fieldUsername].Placeholder = "Username"
	inputs[fieldUsername].CharLimit = 64
	inputs[fieldUsername].Width = 30

	inputs[fieldPassword] = textinput.New()
	inputs[fieldPassword].Placeholder = "Password"
	inputs[fieldPassword].CharLimit = 128
	inputs[fieldPassword].Width = 30

	return listModel{
		list:       l,
		app:        app,
		state:      stateList,
		editInputs: inputs,
	}
}

func (m listModel) Init() tea.Cmd {
	return nil
}
