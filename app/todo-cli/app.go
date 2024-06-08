package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: robust sync of view and crdt representations
// TODO: commands instead of characters

type App struct {
	id   string
	user string

	tasksView list.Model
	tasks     *Tasks
	repo      *Repository
}

func NewApp(device, user string, repo *Repository) *App {
	const width, height = 0, 0
	items := []list.Item{}
	tasks := list.New(items, list.NewDefaultDelegate(), width, height)
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

	return &App{
		id:        device,
		user:      user,
		tasksView: tasks,
		repo:      repo,
	}
}

func (app *App) Init() tea.Cmd {
	return func() tea.Msg {
		tasks, err := app.repo.LoadTasks()
		if err != nil {
			return EventErrorFS(err)
		}

		return EventTasksLoaded(tasks)
	}
}

func (app *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		app.tasksView.SetSize(msg.Width, msg.Height)
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

				task.CreatedBy = app.user

				app.tasks.PushFront(task)
				return app, app.tasksView.InsertItem(0, task)
			}
			form := AddTaskForm(cb)
			return form, form.Init()

		case "-":
			currentIdx := app.tasksView.Cursor()
			app.tasks.Remove(currentIdx)
			app.tasksView.RemoveItem(currentIdx)
			return app, nil

		case "s":
			return app, func() tea.Msg {
				err := app.repo.SaveTasks(app.tasks)
				return EventErrorFS(err)
			}

		case "*":
			newTasks, changes, err := app.repo.Sync(app.tasks)
			if err != nil {
				panic(err)
			}

			if len(changes) == 0 {
				text := "No new changes"
				banner := NewBanner(text, app)
				return banner, banner.Init()
			}

			app.tasks = newTasks

			teatasks := []list.Item{}
			for _, task := range app.tasks.All() {
				teatasks = append(teatasks, task)
			}

			text := "Successfully synced.\n\n"
			for neighbour, changes := range changes {
				text += neighbour + "\n"
				for _, change := range changes {
					text += fmt.Sprintf("  Hash: %s", change.Hash)
				}
			}
			banner := NewBanner(text, app)
			return banner, app.tasksView.SetItems(teatasks)
		}

	case EventTasksLoaded:
		teatasks := []list.Item{}
		for _, task := range msg.All() {
			teatasks = append(teatasks, task)
		}

		app.tasks = msg

		return app, app.tasksView.SetItems(teatasks)

	case EventErrorFS:
		text := "⚠️ ERROR ⚠️\n\n"
		text += msg.Error()
		text += "\n\nPress any key to hide the message."
		banner := NewBanner(text, app)
		return banner, banner.Init()
	}

	var cmd tea.Cmd
	app.tasksView, cmd = app.tasksView.Update(msg)
	return app, cmd
}

func (app *App) View() string {
	return app.tasksView.View()
}
