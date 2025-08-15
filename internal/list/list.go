package list

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	spinner  spinner.Model
	viewport viewport.Model
	style    lipgloss.Style
	content  string
	loading  bool
}

type SetContentMsg struct {
	Content string
}

func New() Model {
	return Model{
		spinner:  spinner.New(spinner.WithSpinner(spinner.Monkey)),
		viewport: viewport.New(0, 0),
		style:    lipgloss.NewStyle(),
		loading:  true,
	}
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
		m.style = m.style.Width(m.viewport.Width)
		m.viewport.SetContent(m.style.Render(m.content))
	case SetContentMsg:
		m.content = msg.Content
		m.viewport.SetContent(m.style.Render(m.content))
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

func (m Model) View() string {
	s := ""
	if m.loading {
		spin := m.spinner.View()
		s = fmt.Sprintf("%s Loading Words %s", spin, spin)
	} else {
		s = m.viewport.View()
	}
	return s
}
