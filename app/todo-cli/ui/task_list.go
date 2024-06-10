package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

func NewTaskList(user, device string) list.Model {
	const width, height = 0, 0
	tasks := list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	tasks.SetFilteringEnabled(true)
	tasks.SetShowStatusBar(true)
	tasks.Title = fmt.Sprintf("TODO Over FS\nUser: %s\nDevice: %s", user, device)

	additional := []key.Binding{
		key.NewBinding(
			key.WithKeys("add task", "+"),
			key.WithHelp("+", "add task"),
		),
		key.NewBinding(
			key.WithKeys("delete task", "-"),
			key.WithHelp("-", "delete task"),
		),
		key.NewBinding(
			key.WithKeys("save state", "s"),
			key.WithHelp("s", "save state"),
		),
		key.NewBinding(
			key.WithKeys("sync state", "*"),
			key.WithHelp("*", "sync state"),
		),
	}

	tasks.AdditionalShortHelpKeys = func() []key.Binding {
		return additional
	}
	tasks.AdditionalFullHelpKeys = func() []key.Binding {
		return additional
	}

	return tasks
}
