package ui

import (
	"fmt"
	"io"

	"github.com/3timeslazy/crdt-over-fs/app/todo-cli/tasks"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type TaskItem struct {
	Name   string
	Author string
	Done   bool
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
		Done:   task.Done,
	}
}

type ItemDelegate struct {
	defaultDelegate list.DefaultDelegate
	doneDelegate    list.DefaultDelegate
}

func NewTaskItemDelegate() *ItemDelegate {
	doneDelegate := list.NewDefaultDelegate()

	doneDelegate.Styles.NormalTitle = doneDelegate.Styles.NormalTitle.Strikethrough(true)
	doneDelegate.Styles.DimmedTitle = doneDelegate.Styles.DimmedTitle.Strikethrough(true)
	doneDelegate.Styles.SelectedTitle = doneDelegate.Styles.SelectedTitle.Strikethrough(true)

	return &ItemDelegate{
		defaultDelegate: list.NewDefaultDelegate(),
		doneDelegate:    doneDelegate,
	}
}

func (d *ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(TaskItem)
	if !ok {
		d.defaultDelegate.Render(w, m, index, listItem)
		return
	}
	if item.Done {
		d.doneDelegate.Render(w, m, index, listItem)
		return
	}

	d.defaultDelegate.Render(w, m, index, listItem)
}

func (d *ItemDelegate) Height() int { return d.defaultDelegate.Height() }

func (d *ItemDelegate) Spacing() int { return d.defaultDelegate.Spacing() }

func (d *ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return d.defaultDelegate.Update(msg, m)
}
