package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLynxCommand(t *testing.T) {
	cmd := NewLynxCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "lynx", cmd.Use)
	assert.Equal(t, "Lynx commands", cmd.Short)
}

func TestLynxInitCommand(t *testing.T) {
	cmd := newLynxInitCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "init [path]", cmd.Use)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"/tmp/test"})
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Initializing Lynx in /tmp/test")
}

func TestLynxDeployCommand(t *testing.T) {
	cmd := newLynxDeployCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "deploy [path]", cmd.Use)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"/tmp/test"})
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Deploying Lynx configuration from /tmp/test")
}

func TestLynxStatusCommand(t *testing.T) {
	cmd := newLynxStatusCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "status", cmd.Use)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Checking Lynx deployment status")
}
