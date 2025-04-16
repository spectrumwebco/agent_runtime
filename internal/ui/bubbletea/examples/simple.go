package examples

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spectrumwebco/agent_runtime/internal/ui/bubbletea"
)

type SimpleModel struct {
	bubbletea.BaseModel
	Count int
}

func NewSimpleModel(count int) SimpleModel {
	return SimpleModel{
		Count: count,
	}
}

func (m SimpleModel) Init() tea.Cmd {
	return bubbletea.Tick(1 * time.Second)
}

func (m SimpleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case bubbletea.TickMsg:
		m.Count--
		if m.Count <= 0 {
			return m, tea.Quit
		}
		return m, bubbletea.Tick(1 * time.Second)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m SimpleModel) View() string {
	return fmt.Sprintf("Countdown: %d\n\nPress q to quit...", m.Count)
}

func RunSimpleExample() error {
	client := bubbletea.NewClient(bubbletea.DefaultOptions()...)
	model := NewSimpleModel(5)
	return client.Run(model)
}
