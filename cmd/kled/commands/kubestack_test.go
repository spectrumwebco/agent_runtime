package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKubestackCommand(t *testing.T) {
	cmd := NewKubestackCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "kubestack", cmd.Use)
	assert.Equal(t, "Kubestack commands", cmd.Short)
}

func TestKubestackInitCommand(t *testing.T) {
	cmd := newKubestackInitCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "init [path]", cmd.Use)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"/tmp/test"})
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Initializing Kubestack in /tmp/test")
}

func TestKubestackApplyCommand(t *testing.T) {
	cmd := newKubestackApplyCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "apply [path]", cmd.Use)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"/tmp/test"})
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Applying Kubestack configuration in /tmp/test")
}

func TestKubestackDestroyCommand(t *testing.T) {
	cmd := newKubestackDestroyCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "destroy [path]", cmd.Use)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"/tmp/test"})
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Destroying Kubestack resources in /tmp/test")
}
