package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerraformOperatorCommand(t *testing.T) {
	cmd := NewTerraformOperatorCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "terraform-operator", cmd.Use)
	assert.Equal(t, "Terraform Operator commands", cmd.Short)
}

func TestTerraformOperatorDeployCommand(t *testing.T) {
	cmd := newTerraformOperatorDeployCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "deploy", cmd.Use)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Deploying Terraform Operator to the cluster")
}

func TestTerraformOperatorListCommand(t *testing.T) {
	cmd := newTerraformOperatorListCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "list", cmd.Use)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Listing Terraform resources managed by the operator")
}

func TestTerraformOperatorApplyCommand(t *testing.T) {
	cmd := newTerraformOperatorApplyCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "apply [name]", cmd.Use)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"test-resource"})
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Applying Terraform resource test-resource")
}
