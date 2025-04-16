package migrations

import (
	"github.com/spectrumwebco/agent_runtime/internal/database"
	"github.com/spectrumwebco/agent_runtime/internal/database/models"
)

func RunMigrations(db *database.DB) error {
	err := db.Migrate(
		&models.User{},
		&models.Workspace{},
		&models.ApiKey{},
		&models.Agent{},
		&models.Tool{},
		&models.Execution{},
		&models.ExecutionStep{},
	)
	if err != nil {
		return err
	}

	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user_workspaces (
			user_id INTEGER NOT NULL,
			workspace_id INTEGER NOT NULL,
			PRIMARY KEY (user_id, workspace_id),
			CONSTRAINT fk_user_workspaces_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			CONSTRAINT fk_user_workspaces_workspace FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
		)
	`).Error
	if err != nil {
		return err
	}

	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS agent_tools (
			agent_id INTEGER NOT NULL,
			tool_id INTEGER NOT NULL,
			PRIMARY KEY (agent_id, tool_id),
			CONSTRAINT fk_agent_tools_agent FOREIGN KEY (agent_id) REFERENCES agents(id) ON DELETE CASCADE,
			CONSTRAINT fk_agent_tools_tool FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE
		)
	`).Error
	if err != nil {
		return err
	}

	return nil
}
