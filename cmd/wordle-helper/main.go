package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var letters map[rune]bool

type model struct {
	green           []rune
	yello           map[rune]bool
	greys           map[rune]bool
	spinner         spinner.Model
	viewport        viewport.Model
	viewportStyle   lipgloss.Style
	viewportContent string
	loading         bool
}

func initialModel() model {
	m := model{
		green:         []rune{'0', '0', '0', '0', '0'},
		yello:         make(map[rune]bool, len(letters)),
		greys:         make(map[rune]bool, len(letters)),
		spinner:       spinner.New(),
		viewport:      viewport.New(0, 0),
		viewportStyle: lipgloss.NewStyle(),
		loading:       true,
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
		m.viewportStyle = m.viewportStyle.Width(msg.Width)
		m.viewport.SetContent(m.viewportStyle.Render(m.viewportContent))
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
	_ = flag.String("f", "", "Text file containing words on seperate lines")
	flag.Parse()
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
