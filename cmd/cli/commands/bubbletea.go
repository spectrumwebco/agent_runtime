package commands

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spectrumwebco/agent_runtime/internal/ui/bubbletea"
	bubbleteaModule "github.com/spectrumwebco/agent_runtime/pkg/modules/bubbletea"
)

func NewBubbleTeaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bubbletea",
		Short: "BubbleTea terminal UI framework",
		Long:  `BubbleTea is a terminal UI framework for building rich terminal user interfaces.`,
	}

	cmd.AddCommand(newBubbleTeaExampleCommand())
	cmd.AddCommand(newBubbleTeaSpinnerCommand())
	cmd.AddCommand(newBubbleTeaProgressCommand())
	cmd.AddCommand(newBubbleTeaListCommand())

	return cmd
}

func newBubbleTeaExampleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "example",
		Short: "Run a BubbleTea example",
		Long:  `Run a simple BubbleTea example to demonstrate the framework.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := bubbleteaModule.NewModule(cfg)
			module.RunExample()

			return nil
		},
	}

	return cmd
}

func newBubbleTeaSpinnerCommand() *cobra.Command {
	var message string
	var duration int

	cmd := &cobra.Command{
		Use:   "spinner",
		Short: "Run a BubbleTea spinner",
		Long:  `Run a BubbleTea spinner with a custom message and duration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := bubbleteaModule.NewModule(cfg)
			
			spinner := bubbletea.NewSpinnerModel()
			spinner.Message = message
			
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
			defer cancel()
			
			return module.RunModelWithContext(ctx, spinner)
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "Loading...", "Message to display next to the spinner")
	cmd.Flags().IntVarP(&duration, "duration", "d", 5, "Duration in seconds to run the spinner")

	return cmd
}

func newBubbleTeaProgressCommand() *cobra.Command {
	var total int
	var label string
	var autoIncrement bool

	cmd := &cobra.Command{
		Use:   "progress",
		Short: "Run a BubbleTea progress bar",
		Long:  `Run a BubbleTea progress bar with a custom total and label.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := bubbleteaModule.NewModule(cfg)
			
			progressBar := bubbletea.NewProgressBarModel(total)
			progressBar.Label = label
			
			if autoIncrement {
				type progressModel struct {
					bubbletea.BaseModel
					ProgressBar bubbletea.ProgressBarModel
				}
				
				model := progressModel{
					ProgressBar: progressBar,
				}
				
				model.Init = func() tea.Cmd {
					return bubbletea.Tick(100 * time.Millisecond)
				}
				
				model.Update = func(msg tea.Msg) (tea.Model, tea.Cmd) {
					switch msg := msg.(type) {
					case bubbletea.TickMsg:
						model.ProgressBar.Increment()
						if model.ProgressBar.Current >= model.ProgressBar.Total {
							return model, tea.Quit
						}
						return model, bubbletea.Tick(100 * time.Millisecond)
					case tea.KeyMsg:
						switch msg.String() {
						case "ctrl+c", "q":
							return model, tea.Quit
						}
					}
					return model, nil
				}
				
				model.View = func() string {
					return model.ProgressBar.View()
				}
				
				return module.RunModel(model)
			}
			
			return module.RunModel(progressBar)
		},
	}

	cmd.Flags().IntVarP(&total, "total", "t", 100, "Total value for the progress bar")
	cmd.Flags().StringVarP(&label, "label", "l", "Progress", "Label to display next to the progress bar")
	cmd.Flags().BoolVarP(&autoIncrement, "auto", "a", true, "Automatically increment the progress bar")

	return cmd
}

func newBubbleTeaListCommand() *cobra.Command {
	var title string
	var itemsFile string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Run a BubbleTea list",
		Long:  `Run a BubbleTea list with custom items.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := bubbleteaModule.NewModule(cfg)
			
			var items []string
			
			if itemsFile != "" {
				data, err := os.ReadFile(itemsFile)
				if err != nil {
					return fmt.Errorf("failed to read items file: %w", err)
				}
				
				lines := strings.Split(string(data), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line != "" {
						items = append(items, line)
					}
				}
			} else {
				items = args
				if len(items) == 0 {
					items = []string{"Item 1", "Item 2", "Item 3", "Item 4", "Item 5"}
				}
			}
			
			list := bubbletea.NewListModel(items)
			list.Title = title
			
			type listModel struct {
				bubbletea.BaseModel
				List bubbletea.ListModel
				Done bool
			}
			
			model := listModel{
				List: list,
			}
			
			model.Init = func() tea.Cmd {
				return nil
			}
			
			model.Update = func(msg tea.Msg) (tea.Model, tea.Cmd) {
				switch msg := msg.(type) {
				case tea.KeyMsg:
					switch msg.String() {
					case "ctrl+c", "q":
						return model, tea.Quit
					case "enter":
						model.Done = true
						return model, tea.Quit
					}
				}
				
				var cmd tea.Cmd
				model.List, cmd = model.List.Update(msg)
				return model, cmd
			}
			
			model.View = func() string {
				if model.Done {
					return fmt.Sprintf("Selected: %s\n", model.List.GetSelectedItem())
				}
				return model.List.View() + "\n\nPress Enter to select, q to quit"
			}
			
			if err := module.RunModel(model); err != nil {
				return err
			}
			
			if model.Done {
				fmt.Printf("You selected: %s\n", model.List.GetSelectedItem())
			}
			
			return nil
		},
	}

	cmd.Flags().StringVarP(&title, "title", "t", "Select an item", "Title to display above the list")
	cmd.Flags().StringVarP(&itemsFile, "file", "f", "", "File containing list items (one per line)")

	return cmd
}
