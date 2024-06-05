package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type AddTask struct {
	input textinput.Model
	cb    AddTaskCallback
}

type AddTaskCallback func(Task, bool) (tea.Model, tea.Cmd)

func AddTaskForm(onTask AddTaskCallback) *AddTask {
	input := textinput.New()
	input.Placeholder = "Task name"

	return &AddTask{
		input: input,
		cb:    onTask,
	}
}

func (form *AddTask) Init() tea.Cmd {
	form.input.Focus()
	return textinput.Blink
}

func (form *AddTask) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			task := Task{
				Name: form.input.Value(),
			}
			return form.cb(task, true)

		case tea.KeyEsc:
			return form.cb(Task{}, false)
		}
	}

	var cmd tea.Cmd
	form.input, cmd = form.input.Update(msg)
	return form, cmd
}

func (form *AddTask) View() string {
	s := "Let's add a new task\n"
	s += form.input.View()
	return s
}
