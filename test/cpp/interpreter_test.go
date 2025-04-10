package cpp_test

import (
	"context"
	"os"
	"testing"

	"github.com/spectrumwebco/agent_runtime/internal/cpp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCppInterpreter(t *testing.T) {
	config := cpp.InterpreterConfig{
		Flags: []string{"-std=c++17", "-O2"},
	}
	
	interpreter, err := cpp.NewInterpreter(config)
	require.NoError(t, err)
	defer interpreter.Close()
	
	t.Run("simple_program", func(t *testing.T) {
		code := `
#include <iostream>

int main() {
    std::cout << "Hello from C++!" << std::endl;
    return 0;
}
`
		ctx := context.Background()
		output, err := interpreter.Exec(ctx, code)
		require.NoError(t, err)
		assert.Contains(t, output, "Hello from C++!")
	})
	
	t.Run("with_input", func(t *testing.T) {
		code := `
#include <iostream>
#include <string>

int main() {
    std::string name;
    std::getline(std::cin, name);
    std::cout << "Hello, " << name << "!" << std::endl;
    return 0;
}
`
		ctx := context.Background()
		output, err := interpreter.ExecWithInput(ctx, code, "Sam Sepiol")
		require.NoError(t, err)
		assert.Contains(t, output, "Hello, Sam Sepiol!")
	})
	
	t.Run("compilation_error", func(t *testing.T) {
		code := `
#include <iostream>

int main() {
    std::cout << "This will not compile" << std::endl
    return 0;
}
`
		ctx := context.Background()
		_, err := interpreter.Exec(ctx, code)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "compilation failed")
	})
	
	t.Run("library_compilation", func(t *testing.T) {
		code := `
#include <iostream>

extern "C" {
    int add(int a, int b) {
        return a + b;
    }
}
`
		ctx := context.Background()
		libPath, err := interpreter.CompileLibrary(ctx, code, "mathlib")
		require.NoError(t, err)
		assert.Contains(t, libPath, "libmathlib.so")
		
		defer func() {
			if err := os.Remove(libPath); err != nil {
				t.Logf("Failed to remove test library: %v", err)
			}
		}()
	})
}
