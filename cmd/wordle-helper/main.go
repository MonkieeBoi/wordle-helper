package main

import (
	"flag"
	"fmt"
	"os"
	"path"
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
	w := wordle.NewWord()
	w[0].Val = '_'
	m := model{
		wordle: wordle.NewWordle(),
		list:   list.New(),
		styles: s,
		word:   w,
	}
	return m
}

func (m model) Init() tea.Cmd {
	return m.list.Init()
}

func (m model) updateKeys(msg tea.KeyMsg) (model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg.Type {
	case tea.KeyBackspace:
		if m.active < wordle.WORD_LEN {
			m.word[m.active].Val = ' '
			m.word[m.active].Colour = wordle.EMPTY
		}
		m.active = max(0, m.active-1)
		m.word[m.active].Colour = wordle.EMPTY
		m.word[m.active].Val = '_'
	case tea.KeyRunes:
		if m.active >= wordle.WORD_LEN {
			return m, nil
		}
		s := msg.String()
		r := []rune(s)[len(s)-1]
		if !('A' <= r && r <= 'Z') && !('a' <= r && r <= 'z') {
			return m, nil
		}
		if msg.Alt {
			m.word[m.active].Colour = wordle.GREEN
		}
		if r <= 'Z' {
			m.word[m.active].Colour = wordle.YELLOW
			r += 'a' - 'A'
		}
		m.word[m.active].Val = r
		if m.word[m.active].Colour == wordle.EMPTY {
			m.word[m.active].Colour = wordle.GREY
		}
		m.active = min(wordle.WORD_LEN, m.active+1)
		if m.active < wordle.WORD_LEN {
			m.word[m.active].Val = '_'
		}
	case tea.KeyEnter:
		if err := m.wordle.AddWord(m.word); err == nil {
			m.active = 0
			m.word = wordle.NewWord()
			m.word[0].Val = '_'
			m.list, cmd = m.list.Update(list.ContentMsg{Content: strings.Join(
				filter.GetWords(
					m.wordle.Greens(),
					m.wordle.Yellows(),
					m.wordle.Greys(),
				),
				" ",
			)})
			cmds = append(cmds, cmd)
		}
	case tea.KeyCtrlC:
		cmds = append(cmds, tea.Quit)
	}
	return m, tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list, cmd = m.list.Update(list.SizeMsg{
			Width:  msg.Width - (wordle.WORD_LEN * 7),
			Height: msg.Height,
		})
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		m, cmd = m.updateKeys(msg)
		cmds = append(cmds, cmd)
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
	for _, word := range w.Board() {
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
	if len(m.wordle.Board()) != 0 {
		right = lipgloss.JoinVertical(0, board, input)
	}
	return lipgloss.JoinHorizontal(0, right, list)
}

func main() {
	fallback := path.Join(os.Getenv("XDG_DATA_HOME"), "wordle-helper", "words")
	wordFile := flag.String("f", fallback, "Text file containing words on seperate lines")
	flag.Parse()
	err := filter.InitWords(*wordFile)
	if err != nil {
		fmt.Printf("Could not open words file at '%s'\n", *wordFile)
		os.Exit(1)
	}
	m := initialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
