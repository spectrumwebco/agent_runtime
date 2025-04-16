package leetcode

import (
	"context"
	"fmt"

	"github.com/spectrumwebco/agent_runtime/internal/leetcode/solver"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Module struct {
	solver *solver.Solver
	config *config.Config
}

func NewModule(cfg *config.Config) *Module {
	s := solver.NewSolver()
	s.InitializeProblems()
	
	return &Module{
		solver: s,
		config: cfg,
	}
}

func (m *Module) Name() string {
	return "leetcode"
}

func (m *Module) Description() string {
	return "LeetCode problem solver module for algorithm practice and interview preparation"
}

func (m *Module) Initialize(ctx context.Context) error {
	return nil
}

func (m *Module) Shutdown(ctx context.Context) error {
	return nil
}

func (m *Module) GetSolver() *solver.Solver {
	return m.solver
}

func (m *Module) SolveProblem(id int, inputs ...interface{}) (*solver.SolverResult, error) {
	return m.solver.Solve(id, inputs...)
}

func (m *Module) VerifySolution(result *solver.SolverResult, expectedOutput interface{}) *solver.SolverResult {
	return m.solver.VerifySolution(result, expectedOutput)
}

func (m *Module) GetProblem(id int) (*solver.Problem, error) {
	return m.solver.GetProblem(id)
}

func (m *Module) GetProblemsByCategory(category solver.ProblemCategory) []*solver.Problem {
	return m.solver.GetProblemsByCategory(category)
}

func (m *Module) GetProblemsByDifficulty(difficulty solver.Difficulty) []*solver.Problem {
	return m.solver.GetProblemsByDifficulty(difficulty)
}

func (m *Module) RunExample() {
	fmt.Println("Running LeetCode solver example...")
	
	problem1, err := m.solver.GetProblem(1)
	if err != nil {
		fmt.Printf("Error getting problem: %v\n", err)
		return
	}
	
	fmt.Printf("Problem: %s (ID: %d)\n", problem1.Title, problem1.ID)
	fmt.Printf("Difficulty: %s\n", problem1.Difficulty)
	fmt.Printf("Description: %s\n\n", problem1.Description)
	
	nums := []int{2, 7, 11, 15}
	target := 9
	
	fmt.Printf("Input: nums = %v, target = %d\n", nums, target)
	
	result, err := m.solver.Solve(1, nums, target)
	if err != nil {
		fmt.Printf("Error solving problem: %v\n", err)
		return
	}
	
	fmt.Printf("Output: %v\n", result.Output)
	fmt.Printf("Execution time: %v\n\n", result.ExecutionTime)
	
	expectedOutput := []int{0, 1}
	result = m.solver.VerifySolution(result, expectedOutput)
	
	if result.Correct {
		fmt.Println("Solution is correct!")
	} else {
		fmt.Printf("Solution is incorrect. Expected: %v, Got: %v\n", expectedOutput, result.Output)
	}
	
	problem3, err := m.solver.GetProblem(3)
	if err != nil {
		fmt.Printf("Error getting problem: %v\n", err)
		return
	}
	
	fmt.Printf("\nProblem: %s (ID: %d)\n", problem3.Title, problem3.ID)
	fmt.Printf("Difficulty: %s\n", problem3.Difficulty)
	fmt.Printf("Description: %s\n\n", problem3.Description)
	
	s := "abcabcbb"
	
	fmt.Printf("Input: s = %s\n", s)
	
	result, err = m.solver.Solve(3, s)
	if err != nil {
		fmt.Printf("Error solving problem: %v\n", err)
		return
	}
	
	fmt.Printf("Output: %v\n", result.Output)
	fmt.Printf("Execution time: %v\n\n", result.ExecutionTime)
	
	expectedOutput = 3
	result = m.solver.VerifySolution(result, expectedOutput)
	
	if result.Correct {
		fmt.Println("Solution is correct!")
	} else {
		fmt.Printf("Solution is incorrect. Expected: %v, Got: %v\n", expectedOutput, result.Output)
	}
}
