package ui

import (
	"fmt"

	"github.com/3timeslazy/crdt-over-fs/app/todo-cli/tasks"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	user string

	tasksView list.Model
	tasks     *tasks.Manager
}

func NewApp(device, user string, tasks *tasks.Manager) *App {
	return &App{
		user:      user,
		tasksView: NewTaskList(user, device),
		tasks:     tasks,
	}
}

func (app *App) Init() tea.Cmd {
	return func() tea.Msg {
		err := app.tasks.Init()
		if err != nil {
			return fmt.Errorf("init tasks manager: %w", err)
		}

		return ManagerReady
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
			cb := func(inputs *TaskAddInputs) (tea.Model, tea.Cmd) {
				if inputs == nil {
					return app, nil
				}

				task := tasks.Task{
					Title:  inputs.Title,
					Author: app.user,
				}
				app.tasks.PushFront(task)
				return app, app.tasksView.InsertItem(0, ToTaskItem(task))
			}
			form := NewTaskAddForm(cb)
			return form, form.Init()

		case "-":
			currentIdx := app.tasksView.Cursor()
			app.tasks.Remove(currentIdx)
			app.tasksView.RemoveItem(currentIdx)
			return app, nil

		case "s":
			return app, func() tea.Msg {
				return app.tasks.Persist()
			}

		case "*":
			changes, err := app.tasks.MergeRemote()
			if err != nil {
				return app, func() tea.Msg {
					return err
				}
			}

			text := "No new changes."

			if len(changes) > 0 {
				text = "Successfully synced.\n\n"
				for _, change := range changes {
					text += "  " + change.Message + "\n"
				}
			}

			teatasks := tasks.Map(app.tasks, ToTaskItem)
			banner := NewBanner(text, app)
			return banner, app.tasksView.SetItems(teatasks)
		}

	case ManagerReadyMsg:
		teatasks := tasks.Map(app.tasks, ToTaskItem)
		return app, app.tasksView.SetItems(teatasks)

	case error:
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
