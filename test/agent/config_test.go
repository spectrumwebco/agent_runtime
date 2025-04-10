package agent_test

import (
	"testing"

	"github.com/spectrumwebco/agent_runtime/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestAgentConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	
	assert.Equal(t, "Sam Sepiol", cfg.Agent.Name)
	assert.Equal(t, "An autonomous software engineering agent", cfg.Agent.Description)
	assert.Equal(t, "0.1.0", cfg.Agent.Version)
	
	
}
