package list

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	viewport viewport.Model
	style    lipgloss.Style
	content  string
}

type SizeMsg struct {
	Width  int
	Height int
}

type ContentMsg struct {
	Content string
}

func New() Model {
	v := viewport.New(0, 0)
	v.KeyMap.Down.SetKeys("ctrl+j")
	v.KeyMap.Up.SetKeys("ctrl+k")
	return Model{
		viewport: v,
		style:    lipgloss.NewStyle().AlignHorizontal(lipgloss.Center),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case SizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
		m.style = m.style.Width(m.viewport.Width)
		m.viewport.SetContent(m.style.Render(m.content))
	case ContentMsg:
		m.content = msg.Content
		m.viewport.SetContent(m.style.Render(m.content))
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.viewport.View()
}
