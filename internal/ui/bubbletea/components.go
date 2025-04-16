package bubbletea

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type SpinnerModel struct {
	Frames   []string
	Current  int
	Speed    time.Duration
	Message  string
	Spinning bool
}

func NewSpinnerModel() SpinnerModel {
	return SpinnerModel{
		Frames:   []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		Speed:    time.Millisecond * 100,
		Spinning: true,
	}
}

func (m SpinnerModel) Init() tea.Cmd {
	return Tick(m.Speed)
}

func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		if !m.Spinning {
			return m, nil
		}
		m.Current = (m.Current + 1) % len(m.Frames)
		return m, Tick(m.Speed)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case " ":
			m.Spinning = !m.Spinning
			if m.Spinning {
				return m, Tick(m.Speed)
			}
			return m, nil
		}
	}
	return m, nil
}

func (m SpinnerModel) View() string {
	if m.Spinning {
		return fmt.Sprintf("%s %s", m.Frames[m.Current], m.Message)
	}
	return fmt.Sprintf("● %s", m.Message)
}

func (m *SpinnerModel) Start() {
	m.Spinning = true
}

func (m *SpinnerModel) Stop() {
	m.Spinning = false
}

func (m *SpinnerModel) SetMessage(msg string) {
	m.Message = msg
}

type ProgressBarModel struct {
	Total      int
	Current    int
	Width      int
	ShowPercentage bool
	Label      string
	BarChar    string
	EmptyChar  string
}

func NewProgressBarModel(total int) ProgressBarModel {
	return ProgressBarModel{
		Total:      total,
		Current:    0,
		Width:      40,
		ShowPercentage: true,
		BarChar:    "█",
		EmptyChar:  "░",
	}
}

func (m ProgressBarModel) Init() tea.Cmd {
	return nil
}

func (m ProgressBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ProgressBarModel) View() string {
	percentage := float64(m.Current) / float64(m.Total)
	if percentage > 1.0 {
		percentage = 1.0
	}
	
	filled := int(float64(m.Width) * percentage)
	
	bar := strings.Repeat(m.BarChar, filled) + strings.Repeat(m.EmptyChar, m.Width-filled)
	
	if m.ShowPercentage {
		return fmt.Sprintf("%s [%s] %d%%", m.Label, bar, int(percentage*100))
	}
	
	return fmt.Sprintf("%s [%s]", m.Label, bar)
}

func (m *ProgressBarModel) SetProgress(current int) {
	m.Current = current
	if m.Current > m.Total {
		m.Current = m.Total
	}
}

func (m *ProgressBarModel) Increment() {
	m.Current++
	if m.Current > m.Total {
		m.Current = m.Total
	}
}

func (m *ProgressBarModel) SetLabel(label string) {
	m.Label = label
}

type ListModel struct {
	Items       []string
	Selected    int
	Width       int
	Height      int
	Title       string
	ShowCursor  bool
	CursorChar  string
}

func NewListModel(items []string) ListModel {
	return ListModel{
		Items:      items,
		Selected:   0,
		ShowCursor: true,
		CursorChar: "→",
	}
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.Selected > 0 {
				m.Selected--
			}
		case "down", "j":
			if m.Selected < len(m.Items)-1 {
				m.Selected++
			}
		case "home":
			m.Selected = 0
		case "end":
			m.Selected = len(m.Items) - 1
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}
	return m, nil
}

func (m ListModel) View() string {
	var sb strings.Builder
	
	if m.Title != "" {
		sb.WriteString(m.Title + "\n\n")
	}
	
	for i, item := range m.Items {
		if i == m.Selected && m.ShowCursor {
			sb.WriteString(fmt.Sprintf("%s %s\n", m.CursorChar, item))
		} else {
			sb.WriteString(fmt.Sprintf("  %s\n", item))
		}
	}
	
	return sb.String()
}

func (m ListModel) GetSelectedItem() string {
	if len(m.Items) == 0 {
		return ""
	}
	return m.Items[m.Selected]
}

func (m ListModel) GetSelectedIndex() int {
	return m.Selected
}

func (m *ListModel) SetItems(items []string) {
	m.Items = items
	if m.Selected >= len(m.Items) {
		m.Selected = len(m.Items) - 1
		if m.Selected < 0 {
			m.Selected = 0
		}
	}
}

func (m *ListModel) SetTitle(title string) {
	m.Title = title
}
