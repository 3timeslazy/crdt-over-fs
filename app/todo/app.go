package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: add sync button

type App struct {
	id    string
	tasks list.Model

	fs *FSWrapper
}

func NewApp(id string) *App {
	const width, height = 0, 0
	items := []list.Item{
		// Task{Name: "Buy groceries"},
		// Task{Name: "Buy a new jacket"},
		// Task{Name: "Play Shadow of the Erdtree"},
	}
	tasks := list.New(items, list.NewDefaultDelegate(), width, height)
	tasks.SetFilteringEnabled(true)
	tasks.SetShowStatusBar(true)
	tasks.Title = fmt.Sprintf("TODO Over FS (%s)", id)

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
	}

	tasks.AdditionalShortHelpKeys = func() []key.Binding {
		return additional
	}
	tasks.AdditionalFullHelpKeys = func() []key.Binding {
		return additional
	}

	return &App{
		tasks: tasks,
		id:    id,
		fs:    NewFSWrapper(id),
	}
}

func (app *App) Init() tea.Cmd {
	return func() tea.Msg {
		err := app.fs.SetupDir()
		if err != nil {
			return EventErrorFS(err)
		}

		tasks, err := app.fs.LoadTasks()
		if err != nil {
			return EventErrorFS(err)
		}

		return EventTasksLoaded(tasks)
	}
}

func (app *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		app.tasks.SetSize(msg.Width, msg.Height)
		return app, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return app, tea.Quit

		case "+":
			cb := func(task Task, added bool) (tea.Model, tea.Cmd) {
				if !added {
					return app, nil
				}

				return app, app.tasks.InsertItem(0, task)
			}
			form := AddTaskForm(cb)
			return form, form.Init()

		case "-":
			app.tasks.RemoveItem(app.tasks.Cursor())
			return app, nil

		case "s":
			tasks := []Task{}
			for _, task := range app.tasks.Items() {
				tasks = append(tasks, task.(Task))
			}

			return app, func() tea.Msg {
				err := app.fs.SaveTasks(tasks)
				return EventErrorFS(err)
			}
		}

	case EventTasksLoaded:
		teatasks := []list.Item{}
		for _, task := range msg {
			teatasks = append(teatasks, task)
		}

		return app, app.tasks.SetItems(teatasks)

	case EventErrorFS:
		text := "⚠️ ERROR ⚠️\n\n"
		text += msg.Error()
		text += "\n\nPress any key to hide the message."
		banner := NewBanner(text, app)
		return banner, banner.Init()
	}

	var cmd tea.Cmd
	app.tasks, cmd = app.tasks.Update(msg)
	return app, cmd
}

func (app *App) View() string {
	return app.tasks.View()
}
