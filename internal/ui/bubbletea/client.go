package bubbletea

import (
	"context"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Client struct {
	program *tea.Program
	options []tea.ProgramOption
}

func NewClient(opts ...tea.ProgramOption) *Client {
	return &Client{
		options: opts,
	}
}

func DefaultOptions() []tea.ProgramOption {
	return []tea.ProgramOption{
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	}
}

func (c *Client) Run(model tea.Model) error {
	c.program = tea.NewProgram(model, c.options...)
	_, err := c.program.Run()
	return err
}

func (c *Client) RunWithContext(ctx context.Context, model tea.Model) error {
	c.program = tea.NewProgram(model, c.options...)
	
	errCh := make(chan error, 1)
	go func() {
		_, err := c.program.Run()
		errCh <- err
	}()

	select {
	case <-ctx.Done():
		c.program.Quit()
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (c *Client) Quit() {
	if c.program != nil {
		c.program.Quit()
	}
}

func (c *Client) Send(msg tea.Msg) {
	if c.program != nil {
		c.program.Send(msg)
	}
}

func EnableLogging(filename string, prefix string) error {
	if filename == "" {
		filename = "bubbletea.log"
	}
	
	if prefix == "" {
		prefix = "kled"
	}
	
	_, err := tea.LogToFile(filename, prefix)
	return err
}

func WithOutput(w io.Writer) tea.ProgramOption {
	return tea.WithOutput(w)
}

func WithInput(r io.Reader) tea.ProgramOption {
	return tea.WithInput(r)
}

func WithAltScreen() tea.ProgramOption {
	return tea.WithAltScreen()
}

func WithMouseCellMotion() tea.ProgramOption {
	return tea.WithMouseCellMotion()
}

func WithoutCatchPanics() tea.ProgramOption {
	return tea.WithoutCatchPanics()
}

func PrintInline(msg string) {
	fmt.Fprint(os.Stdout, msg)
}

func PrintlnInline(msg string) {
	fmt.Fprintln(os.Stdout, msg)
}
