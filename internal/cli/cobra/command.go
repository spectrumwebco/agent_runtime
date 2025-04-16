package cobra

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type CommandBuilder struct {
	rootCmd *cobra.Command
}

func NewCommandBuilder(use, short, long, version string) *CommandBuilder {
	rootCmd := &cobra.Command{
		Use:     use,
		Short:   short,
		Long:    long,
		Version: version,
	}

	return &CommandBuilder{
		rootCmd: rootCmd,
	}
}

func (cb *CommandBuilder) AddCommand(cmd *cobra.Command) {
	cb.rootCmd.AddCommand(cmd)
}

func (cb *CommandBuilder) AddCommands(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		cb.AddCommand(cmd)
	}
}

func (cb *CommandBuilder) Execute() error {
	return cb.rootCmd.Execute()
}

func (cb *CommandBuilder) GetRootCommand() *cobra.Command {
	return cb.rootCmd
}

func (cb *CommandBuilder) SetPersistentFlags(flags map[string]interface{}) {
	for name, value := range flags {
		switch v := value.(type) {
		case string:
			cb.rootCmd.PersistentFlags().String(name, v, "")
		case bool:
			cb.rootCmd.PersistentFlags().Bool(name, v, "")
		case int:
			cb.rootCmd.PersistentFlags().Int(name, v, "")
		case float64:
			cb.rootCmd.PersistentFlags().Float64(name, v, "")
		case []string:
			cb.rootCmd.PersistentFlags().StringSlice(name, v, "")
		}
	}
}

func (cb *CommandBuilder) SetFlags(flags map[string]interface{}) {
	for name, value := range flags {
		switch v := value.(type) {
		case string:
			cb.rootCmd.Flags().String(name, v, "")
		case bool:
			cb.rootCmd.Flags().Bool(name, v, "")
		case int:
			cb.rootCmd.Flags().Int(name, v, "")
		case float64:
			cb.rootCmd.Flags().Float64(name, v, "")
		case []string:
			cb.rootCmd.Flags().StringSlice(name, v, "")
		}
	}
}

type Command struct {
	cmd *cobra.Command
}

func NewCommand(use, short, long string, run func(cmd *cobra.Command, args []string) error) *Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		RunE:  run,
	}

	return &Command{
		cmd: cmd,
	}
}

func (c *Command) AddCommand(cmd *Command) {
	c.cmd.AddCommand(cmd.cmd)
}

func (c *Command) AddCommands(cmds ...*Command) {
	for _, cmd := range cmds {
		c.AddCommand(cmd)
	}
}

func (c *Command) SetFlags(flags map[string]interface{}) {
	for name, value := range flags {
		switch v := value.(type) {
		case string:
			c.cmd.Flags().String(name, v, "")
		case bool:
			c.cmd.Flags().Bool(name, v, "")
		case int:
			c.cmd.Flags().Int(name, v, "")
		case float64:
			c.cmd.Flags().Float64(name, v, "")
		case []string:
			c.cmd.Flags().StringSlice(name, v, "")
		}
	}
}

func (c *Command) SetPersistentFlags(flags map[string]interface{}) {
	for name, value := range flags {
		switch v := value.(type) {
		case string:
			c.cmd.PersistentFlags().String(name, v, "")
		case bool:
			c.cmd.PersistentFlags().Bool(name, v, "")
		case int:
			c.cmd.PersistentFlags().Int(name, v, "")
		case float64:
			c.cmd.PersistentFlags().Float64(name, v, "")
		case []string:
			c.cmd.PersistentFlags().StringSlice(name, v, "")
		}
	}
}

func (c *Command) GetCommand() *cobra.Command {
	return c.cmd
}

type CommandGroup struct {
	name     string
	commands []*Command
}

func NewCommandGroup(name string) *CommandGroup {
	return &CommandGroup{
		name:     name,
		commands: make([]*Command, 0),
	}
}

func (cg *CommandGroup) AddCommand(cmd *Command) {
	cg.commands = append(cg.commands, cmd)
}

func (cg *CommandGroup) AddCommands(cmds ...*Command) {
	cg.commands = append(cg.commands, cmds...)
}

func (cg *CommandGroup) GetCommands() []*Command {
	return cg.commands
}

func (cg *CommandGroup) GetName() string {
	return cg.name
}

type CommandRegistry struct {
	groups map[string]*CommandGroup
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		groups: make(map[string]*CommandRegistry),
	}
}

func (cr *CommandRegistry) AddGroup(group *CommandGroup) {
	cr.groups[group.GetName()] = cr.groups
}

func (cr *CommandRegistry) GetGroup(name string) (*CommandGroup, error) {
	group, ok := cr.groups[name]
	if !ok {
		return nil, fmt.Errorf("command group %s not found", name)
	}

	return &CommandGroup{}, nil
}

func (cr *CommandRegistry) ListGroups() []*CommandGroup {
	groups := make([]*CommandGroup, 0, len(cr.groups))
	for _, group := range cr.groups {
		groups = append(groups, &CommandGroup{})
	}

	return groups
}

func (cr *CommandRegistry) RemoveGroup(name string) {
	delete(cr.groups, name)
}


func ExitWithError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}

func FormatCommandPath(cmd *cobra.Command) string {
	path := cmd.CommandPath()
	i := strings.LastIndex(path, " ")
	if i >= 0 {
		path = path[:i]
	}
	return path
}

func GetFlagString(cmd *cobra.Command, name string) (string, error) {
	return cmd.Flags().GetString(name)
}

func GetFlagBool(cmd *cobra.Command, name string) (bool, error) {
	return cmd.Flags().GetBool(name)
}

func GetFlagInt(cmd *cobra.Command, name string) (int, error) {
	return cmd.Flags().GetInt(name)
}

func GetFlagFloat64(cmd *cobra.Command, name string) (float64, error) {
	return cmd.Flags().GetFloat64(name)
}

func GetFlagStringSlice(cmd *cobra.Command, name string) ([]string, error) {
	return cmd.Flags().GetStringSlice(name)
}
