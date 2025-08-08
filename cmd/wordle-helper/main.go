package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MonkieeBoi/wordle-helper/internal/list"
	tea "github.com/charmbracelet/bubbletea"
)

var letters map[rune]bool

type model struct {
	green []rune
	yello map[rune][]int
	greys map[rune]bool
	list  list.Model
}

func initialModel() model {
	m := model{
		green: []rune{'0', '0', '0', '0', '0'},
		yello: make(map[rune][]int, len(letters)),
		greys: make(map[rune]bool, len(letters)),
		list:  list.New(),
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
