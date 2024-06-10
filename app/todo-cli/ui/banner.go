package ui

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type Banner struct {
	text     textarea.Model
	redirect tea.Model
}

func NewBanner(text string, redirect tea.Model) *Banner {
	area := textarea.New()
	area.SetValue(text)
	area.Focus()

	return &Banner{
		text:     area,
		redirect: redirect,
	}
}

func (b *Banner) Init() tea.Cmd {
	return nil
}

func (b *Banner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return b.redirect, nil
	}

	return b, nil
}

func (b *Banner) View() string {
	return b.text.Value()
}
