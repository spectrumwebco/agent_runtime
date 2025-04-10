package python_test

import (
	"context"
	"testing"

	"github.com/spectrumwebco/agent_runtime/pkg/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPythonTool(t *testing.T) {
	pythonTool, err := tools.NewPythonTool("python", "Executes Python code")
	require.NoError(t, err)
	defer pythonTool.Close()
	
	ctx := context.Background()
	
	t.Run("simple_execution", func(t *testing.T) {
		params := map[string]interface{}{
			"code": `
print("Hello from Python tool!")
`,
		}
		
		result, err := pythonTool.Execute(ctx, params)
		require.NoError(t, err)
		assert.Contains(t, result.(string), "Hello from Python tool!")
	})
	
	t.Run("with_file", func(t *testing.T) {
		pythonCode := `
print("Hello from Python file!")
`
		
		require.NoError(t, tools.WriteFile("/tmp/test_script.py", pythonCode))
		defer tools.DeleteFile("/tmp/test_script.py")
		
		params := map[string]interface{}{
			"file": "/tmp/test_script.py",
		}
		
		result, err := pythonTool.Execute(ctx, params)
		require.NoError(t, err)
		assert.Contains(t, result.(string), "Hello from Python file!")
	})
	
	t.Run("with_function", func(t *testing.T) {
		moduleCode := `
def add(a, b):
    return a + b
`
		
		require.NoError(t, tools.WriteFile("/tmp/math_module.py", moduleCode))
		defer tools.DeleteFile("/tmp/math_module.py")
		
		params := map[string]interface{}{
			"code": `
import sys
sys.path.append("/tmp")
import math_module
print(math_module.add(2, 3))
`,
		}
		
		result, err := pythonTool.Execute(ctx, params)
		require.NoError(t, err)
		assert.Contains(t, result.(string), "5")
	})
	
	t.Run("syntax_error", func(t *testing.T) {
		params := map[string]interface{}{
			"code": `
print("Missing closing parenthesis"
`,
		}
		
		_, err := pythonTool.Execute(ctx, params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to execute Python code")
	})
}
