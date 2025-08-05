package main

import (
    "fmt"
    "os"

    tea "github.com/charmbracelet/bubbletea"
)

type model struct {
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
		}
	}

    return m, nil
}

func (m model) View() string {
    return "Wordle"
}

func main() {
    p := tea.NewProgram(model{}, tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
