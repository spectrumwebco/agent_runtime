package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/pkg/modules/gorm"
)

func NewGORMCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gorm",
		Short: "GORM database operations",
		Long:  `Perform database operations using GORM ORM.`,
	}

	cmd.AddCommand(newGORMMigrateCommand())
	cmd.AddCommand(newGORMStatusCommand())
	cmd.AddCommand(newGORMListModelsCommand())

	return cmd
}

func newGORMMigrateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Long:  `Run database migrations to create or update database schema.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := gorm.NewModule(cfg)
			err = module.Initialize()
			if err != nil {
				return fmt.Errorf("failed to initialize GORM module: %w", err)
			}

			err = module.RunMigrations()
			if err != nil {
				return fmt.Errorf("failed to run migrations: %w", err)
			}

			fmt.Println("Migrations completed successfully")
			return nil
		},
	}

	return cmd
}

func newGORMStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check database connection status",
		Long:  `Check the status of the database connection.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := gorm.NewModule(cfg)
			err = module.Initialize()
			if err != nil {
				return fmt.Errorf("failed to initialize GORM module: %w", err)
			}

			db := module.GetDatabase()
			err = db.Ping()
			if err != nil {
				return fmt.Errorf("database connection failed: %w", err)
			}

			fmt.Println("Database connection successful")
			fmt.Printf("Database type: %s\n", db.Config.Type)
			fmt.Printf("Database name: %s\n", db.Config.Database)
			if db.Config.Type != "sqlite" {
				fmt.Printf("Database host: %s\n", db.Config.Host)
				fmt.Printf("Database port: %d\n", db.Config.Port)
			}

			return nil
		},
	}

	return cmd
}

func newGORMListModelsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-models",
		Short: "List available database models",
		Long:  `List all available database models that can be migrated.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Available database models:")
			fmt.Println("- User")
			fmt.Println("- Role")
			fmt.Println("- Permission")
			fmt.Println("- Agent")
			fmt.Println("- Task")
			fmt.Println("- Tool")
			fmt.Println("- AgentTool")
			fmt.Println("- Workspace")
			fmt.Println("- Session")
			fmt.Println("- AuditLog")
			fmt.Println("- Setting")
			fmt.Println("- APIKey")
			fmt.Println("- Notification")

			return nil
		},
	}

	return cmd
}
