package cpp_test

import (
	"context"
	"testing"

	"github.com/spectrumwebco/agent_runtime/pkg/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCppTool(t *testing.T) {
	cppTool, err := tools.NewCppTool("cpp", "Executes C++ code")
	require.NoError(t, err)
	defer cppTool.Close()
	
	ctx := context.Background()
	
	t.Run("simple_execution", func(t *testing.T) {
		params := map[string]interface{}{
			"code": `
#include <iostream>

int main() {
    std::cout << "Hello from C++ tool!" << std::endl;
    return 0;
}
`,
		}
		
		result, err := cppTool.Execute(ctx, params)
		require.NoError(t, err)
		assert.Contains(t, result.(string), "Hello from C++ tool!")
	})
	
	t.Run("with_input", func(t *testing.T) {
		params := map[string]interface{}{
			"code": `
#include <iostream>
#include <string>

int main() {
    std::string name;
    std::getline(std::cin, name);
    std::cout << "Hello, " << name << "!" << std::endl;
    return 0;
}
`,
			"input": "Sam Sepiol",
		}
		
		result, err := cppTool.Execute(ctx, params)
		require.NoError(t, err)
		assert.Contains(t, result.(string), "Hello, Sam Sepiol!")
	})
	
	t.Run("with_custom_flags", func(t *testing.T) {
		params := map[string]interface{}{
			"code": `
#include <iostream>

int main() {
    #ifdef DEBUG
    std::cout << "Debug mode enabled" << std::endl;
    #else
    std::cout << "Debug mode disabled" << std::endl;
    #endif
    return 0;
}
`,
			"flags": []interface{}{"-std=c++17", "-DDEBUG"},
		}
		
		result, err := cppTool.Execute(ctx, params)
		require.NoError(t, err)
		assert.Contains(t, result.(string), "Debug mode enabled")
	})
	
	t.Run("with_include_dirs", func(t *testing.T) {
		headerCode := `
#ifndef TEST_HEADER_H
#define TEST_HEADER_H

#include <string>

inline std::string getGreeting() {
    return "Hello from custom header!";
}

#endif // TEST_HEADER_H
`
		
		params := map[string]interface{}{
			"code": `
#include <iostream>
#include "test_header.h"

int main() {
    std::cout << getGreeting() << std::endl;
    return 0;
}
`,
			"include_dirs": []interface{}{"/tmp"},
		}
		
		require.NoError(t, tools.WriteFile("/tmp/test_header.h", headerCode))
		defer tools.DeleteFile("/tmp/test_header.h")
		
		result, err := cppTool.Execute(ctx, params)
		require.NoError(t, err)
		assert.Contains(t, result.(string), "Hello from custom header!")
	})
	
	t.Run("library_compilation", func(t *testing.T) {
		params := map[string]interface{}{
			"code": `
#include <iostream>

extern "C" {
    int add(int a, int b) {
        return a + b;
    }
}
`,
			"library": "mathlib",
		}
		
		result, err := cppTool.Execute(ctx, params)
		require.NoError(t, err)
		assert.Contains(t, result.(string), "libmathlib.so")
	})
	
	t.Run("compilation_error", func(t *testing.T) {
		params := map[string]interface{}{
			"code": `
#include <iostream>

int main() {
    std::cout << "Missing semicolon" << std::endl
    return 0;
}
`,
		}
		
		_, err := cppTool.Execute(ctx, params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "compilation failed")
	})
}
