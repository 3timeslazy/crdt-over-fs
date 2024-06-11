package ui

import (
	"fmt"

	"github.com/3timeslazy/crdt-over-fs/app/todo-cli/tasks"

	"github.com/charmbracelet/bubbles/list"
)

type TaskItem struct {
	Name   string
	Author string
}

func (item TaskItem) FilterValue() string {
	return item.Name
}

func (item TaskItem) Title() string { return item.Name }

func (item TaskItem) Description() string {
	s := fmt.Sprintf("by: %s\n", item.Author)
	return s
}

func ToTaskItem(task tasks.Task) list.Item {
	return TaskItem{
		Name:   task.Title,
		Author: task.Author,
	}
}
