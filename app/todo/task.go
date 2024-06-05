package main

type Task struct {
	Name string
}

func (task Task) FilterValue() string {
	return task.Name
}

func (task Task) Title() string { return task.Name }

func (task Task) Description() string { return "" }
