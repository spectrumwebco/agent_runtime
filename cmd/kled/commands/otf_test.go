package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOTFCommand(t *testing.T) {
	cmd := NewOTFCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "otf", cmd.Use)
	assert.Equal(t, "Open Terraform Framework commands", cmd.Short)
}

func TestOTFInitCommand(t *testing.T) {
	cmd := newOTFInitCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "init [path]", cmd.Use)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"/tmp/test"})
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Initializing OTF in /tmp/test")
}

func TestOTFApplyCommand(t *testing.T) {
	cmd := newOTFApplyCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "apply [path]", cmd.Use)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"/tmp/test"})
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Applying OTF configuration in /tmp/test")
}

func TestOTFDestroyCommand(t *testing.T) {
	cmd := newOTFDestroyCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "destroy [path]", cmd.Use)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"/tmp/test"})
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Destroying OTF resources in /tmp/test")
}
