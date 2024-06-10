package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AddTaskForm struct {
	title textinput.Model
	cb    TaskAddCallback
}

type TaskAddCallback func(*TaskAddInputs) (tea.Model, tea.Cmd)

type TaskAddInputs struct {
	Title string
}

func NewTaskAddForm(cb TaskAddCallback) *AddTaskForm {
	title := textinput.New()
	title.Placeholder = "Title"
	title.Focus()

	return &AddTaskForm{
		title: title,
		cb:    cb,
	}
}

func (form *AddTaskForm) Init() tea.Cmd {
	return textinput.Blink
}

func (form *AddTaskForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return form.cb(&TaskAddInputs{
				Title: form.title.Value(),
			})

		case tea.KeyEsc:
			return form.cb(nil)
		}
	}

	var cmd tea.Cmd
	form.title, cmd = form.title.Update(msg)
	return form, cmd
}

var titleStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("62")).
	Foreground(lipgloss.Color("230")).
	Padding(0, 1)

func (form *AddTaskForm) View() string {
	s := titleStyle.Render("Let's add a new task") + "\n\n"
	s += form.title.View() + "\n"
	return s
}
