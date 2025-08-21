package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

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
		case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
			"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z":
			if m.active < wordle.WORD_LEN {
				m.word[m.active].Val = []rune(msg.String())[0]
				if m.word[m.active].Colour == wordle.EMPTY {
					m.word[m.active].Colour = wordle.GREY
				}
				m.active = min(wordle.WORD_LEN, m.active+1)
				if m.active < wordle.WORD_LEN {
					m.word[m.active].Val = '_'
				}
			}
		case "enter":
			if err := m.wordle.AddWord(m.word); err == nil {
				m.active = 0
				m.word = wordle.Word{
					wordle.Char{Val: '_'},
					wordle.Char{Val: ' '},
					wordle.Char{Val: ' '},
					wordle.Char{Val: ' '},
					wordle.Char{Val: ' '},
				}
				m.list, cmd = m.list.Update(list.SetContentMsg{Content: strings.Join(
					filter.GetWords(
						m.wordle.Green,
						m.wordle.Yello,
						m.wordle.Greys,
					),
					" ",
				)})
				cmds = append(cmds, cmd)
			}
		case "ctrl+c":
			cmds = append(cmds, tea.Quit)
		}
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func cellView(c wordle.Char, s styles) string {
	style := s.wordleStyles[c.Colour]
	padTB := style.Render("       ")
	padLR := style.Render("   ")
	ch := lipgloss.JoinHorizontal(0, padLR, style.Render(string(c.Val)), padLR)
	return lipgloss.JoinVertical(lipgloss.Center, padTB, ch, padTB)
}

func wordView(w wordle.Word, s styles) string {
	cells := []string{}
	for _, c := range w {
		cells = append(cells, cellView(c, s))
	}
	return lipgloss.JoinHorizontal(0, cells...)
}

func boardView(w wordle.Wordle, s styles) string {
	words := []string{}
	for _, word := range w.Board {
		renderedWord := ""
		renderedWord = wordView(word, s)
		words = append(words, renderedWord)
	}
	return lipgloss.JoinVertical(0, words...)
}

func (m model) View() string {
	list := m.list.View()
	board := boardView(m.wordle, m.styles)
	input := wordView(m.word, m.styles)
	right := input
	if len(m.wordle.Board) != 0 {
		right = lipgloss.JoinVertical(0, board, input)
	}
	return lipgloss.JoinHorizontal(0, right, list)
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
