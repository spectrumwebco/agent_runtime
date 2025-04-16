package bubbletea

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type BaseModel struct {
	Width  int
	Height int
}

func (m BaseModel) Init() tea.Cmd {
	return nil
}

func (m BaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}
	return m, nil
}

func (m BaseModel) View() string {
	return "Base Model"
}

type TickMsg time.Time

func Tick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

type ErrorMsg struct {
	Err error
}

func (e ErrorMsg) Error() string {
	return e.Err.Error()
}

func NewErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg{Err: err}
	}
}

type SequentialModel struct {
	Models     []tea.Model
	CurrentIdx int
	Width      int
	Height     int
}

func NewSequentialModel(models ...tea.Model) SequentialModel {
	return SequentialModel{
		Models:     models,
		CurrentIdx: 0,
	}
}

func (m SequentialModel) Init() tea.Cmd {
	if len(m.Models) == 0 {
		return nil
	}
	return m.Models[m.CurrentIdx].Init()
}

func (m SequentialModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(m.Models) == 0 {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		
		for i := range m.Models {
			var cmd tea.Cmd
			m.Models[i], cmd = m.Models[i].Update(msg)
			if cmd != nil {
				return m, cmd
			}
		}
		
		return m, nil
		
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
		
		if msg.String() == "right" || msg.String() == "n" {
			m.CurrentIdx = (m.CurrentIdx + 1) % len(m.Models)
			return m, m.Models[m.CurrentIdx].Init()
		}
		
		if msg.String() == "left" || msg.String() == "p" {
			m.CurrentIdx = (m.CurrentIdx - 1 + len(m.Models)) % len(m.Models)
			return m, m.Models[m.CurrentIdx].Init()
		}
	}

	var cmd tea.Cmd
	m.Models[m.CurrentIdx], cmd = m.Models[m.CurrentIdx].Update(msg)
	return m, cmd
}

func (m SequentialModel) View() string {
	if len(m.Models) == 0 {
		return "No models to display"
	}
	
	view := m.Models[m.CurrentIdx].View()
	
	nav := fmt.Sprintf("\n\n[%d/%d] (← p/n →) ", m.CurrentIdx+1, len(m.Models))
	return view + nav
}
