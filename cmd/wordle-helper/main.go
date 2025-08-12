package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MonkieeBoi/wordle-helper/internal/list"
	"github.com/MonkieeBoi/wordle-helper/internal/wordle"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type styles struct {
	wordleStyles map[wordle.Colour]lipgloss.Style
}

var letters map[rune]bool

type model struct {
	wordle wordle.Wordle
	list   list.Model
	styles styles
}

func initialModel() model {
	s := styles{
		wordleStyles: map[wordle.Colour]lipgloss.Style{
			wordle.GREEN: lipgloss.NewStyle().
				Background(lipgloss.Color("2")).
				Foreground(lipgloss.Color("0")),
			wordle.YELLOW: lipgloss.NewStyle().
				Background(lipgloss.Color("3")).
				Foreground(lipgloss.Color("0")),
			wordle.GREY: lipgloss.NewStyle().
				Background(lipgloss.Color("0")).
				Foreground(lipgloss.Color("7")),
			wordle.EMPTY: lipgloss.NewStyle().
				Background(lipgloss.Color("0")).
				Foreground(lipgloss.Color("7")),
		},
	}
	m := model{
		wordle: wordle.NewWordle(),
		list:   list.New(),
		styles: s,
	}
	return m
}

func (m model) Init() tea.Cmd {
	return m.list.Init()
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
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	s := m.list.View()
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
