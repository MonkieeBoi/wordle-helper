package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MonkieeBoi/wordle-helper/internal/filter"
	"github.com/MonkieeBoi/wordle-helper/internal/list"
	"github.com/MonkieeBoi/wordle-helper/internal/wordle"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type styles struct {
	wordleStyles map[wordle.Colour]lipgloss.Style
}

type model struct {
	wordle wordle.Wordle
	list   list.Model
	styles styles
	word   wordle.Word
	active int
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
				Foreground(lipgloss.Color("7")).
				Blink(true),
		},
	}
	m := model{
		wordle: wordle.NewWordle(),
		list:   list.New(),
		styles: s,
		word: wordle.Word{
			wordle.Char{Val: '_'},
			wordle.Char{Val: ' '},
			wordle.Char{Val: ' '},
			wordle.Char{Val: ' '},
			wordle.Char{Val: ' '},
		},
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
		case "ctrl+g":
			m.word[m.active].Colour = wordle.GREEN
		case "ctrl+y":
			m.word[m.active].Colour = wordle.YELLOW
		case "backspace":
			if m.active < wordle.WORD_LEN {
				m.word[m.active].Val = ' '
				m.word[m.active].Colour = wordle.EMPTY
			}
			m.active = max(0, m.active-1)
			m.word[m.active].Colour = wordle.EMPTY
			m.word[m.active].Val = '_'
		case "ctrl+c":
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
	wordFile := flag.String("f", "", "Text file containing words on seperate lines")
	flag.Parse()
	err := filter.InitWords(*wordFile)
	if err != nil {
		fmt.Println("Could not open word file!")
		os.Exit(1)
	}
	m := initialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
