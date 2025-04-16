package bubbletea

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spectrumwebco/agent_runtime/internal/ui/bubbletea"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Module struct {
	client *bubbletea.Client
	config *config.Config
}

func NewModule(cfg *config.Config) *Module {
	return &Module{
		client: bubbletea.NewClient(bubbletea.DefaultOptions()...),
		config: cfg,
	}
}

func (m *Module) Name() string {
	return "bubbletea"
}

func (m *Module) Description() string {
	return "Terminal UI framework for building rich terminal user interfaces"
}

func (m *Module) Initialize(ctx context.Context) error {
	return nil
}

func (m *Module) Shutdown(ctx context.Context) error {
	if m.client != nil {
		m.client.Quit()
	}
	return nil
}

func (m *Module) GetClient() *bubbletea.Client {
	return m.client
}

func (m *Module) RunModel(model tea.Model) error {
	return m.client.Run(model)
}

func (m *Module) RunModelWithContext(ctx context.Context, model tea.Model) error {
	return m.client.RunWithContext(ctx, model)
}

func (m *Module) CreateSpinner() bubbletea.SpinnerModel {
	return bubbletea.NewSpinnerModel()
}

func (m *Module) CreateProgressBar(total int) bubbletea.ProgressBarModel {
	return bubbletea.NewProgressBarModel(total)
}

func (m *Module) CreateList(items []string) bubbletea.ListModel {
	return bubbletea.NewListModel(items)
}

func (m *Module) CreateSequentialModel(models ...tea.Model) bubbletea.SequentialModel {
	return bubbletea.NewSequentialModel(models...)
}

func (m *Module) RunExample() {
	fmt.Println("Running BubbleTea example...")
	
	type countdownModel struct {
		bubbletea.BaseModel
		Count int
	}
	
	initFunc := func(m countdownModel) tea.Cmd {
		return bubbletea.Tick(1 * time.Second)
	}
	
	updateFunc := func(m countdownModel, msg tea.Msg) (tea.Model, tea.Cmd) {
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
	
	viewFunc := func(m countdownModel) string {
		return fmt.Sprintf("Countdown: %d\n\nPress q to quit...", m.Count)
	}
	
	model := countdownModel{Count: 5}
	model.Init = func() tea.Cmd { return initFunc(model) }
	model.Update = func(msg tea.Msg) (tea.Model, tea.Cmd) { return updateFunc(model, msg) }
	model.View = func() string { return viewFunc(model) }
	
	if err := m.RunModel(model); err != nil {
		fmt.Printf("Error running example: %v\n", err)
	}
}
