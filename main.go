package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

var letters map[rune]bool

type model struct {
	green    []rune
	yello    map[rune]bool
	greys    map[rune]bool
	spinner  spinner.Model
	viewport viewport.Model
	loading  bool
}

func initialModel() model {
	m := model{
		green:    []rune{'0', '0', '0', '0', '0'},
		yello:    make(map[rune]bool, len(letters)),
		greys:    make(map[rune]bool, len(letters)),
		spinner:  spinner.New(),
		viewport: viewport.New(0, 0),
		loading:  true,
	}
	m.spinner.Spinner = spinner.Monkey
	return m
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmds = append(cmds, tea.Quit)
		}
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
	}

	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	s := ""
	if m.loading {
		spin := m.spinner.View()
		s = fmt.Sprintf("%s Loading Words %s", spin, spin)
	} else {
		s = m.viewport.View()
	}
	return s
}

func main() {
	m := initialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	letters = make(map[rune]bool, 'z'-'a'+1)
	for i := 'a'; i <= 'z'; i++ {
		letters[rune(i)] = true
	}
}
