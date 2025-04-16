package examples

import (
	"fmt"
	"time"

	"github.com/spectrumwebco/agent_runtime/internal/leetcode/solver"
	"github.com/spectrumwebco/agent_runtime/internal/leetcode/structures"
)

func RunTwoSumExample() {
	fmt.Println("Running Two Sum Example")
	fmt.Println("=======================")

	s := solver.NewSolver()
	s.InitializeProblems()

	problem, err := s.GetProblem(1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Problem: %s (ID: %d)\n", problem.Title, problem.ID)
	fmt.Printf("Difficulty: %s\n", problem.Difficulty)
	fmt.Printf("Description: %s\n\n", problem.Description)

	nums := []int{2, 7, 11, 15}
	target := 9

	fmt.Printf("Input: nums = %v, target = %d\n", nums, target)

	startTime := time.Now()
	result, err := s.Solve(1, nums, target)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	executionTime := time.Since(startTime)

	fmt.Printf("Output: %v\n", result.Output)
	fmt.Printf("Execution time: %v\n\n", executionTime)

	expectedOutput := []int{0, 1}
	result = s.VerifySolution(result, expectedOutput)

	if result.Correct {
		fmt.Println("Solution is correct!")
	} else {
		fmt.Printf("Solution is incorrect. Expected: %v, Got: %v\n", expectedOutput, result.Output)
	}
}

func RunAddTwoNumbersExample() {
	fmt.Println("Running Add Two Numbers Example")
	fmt.Println("==============================")

	s := solver.NewSolver()
	s.InitializeProblems()

	problem, err := s.GetProblem(2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Problem: %s (ID: %d)\n", problem.Title, problem.ID)
	fmt.Printf("Difficulty: %s\n", problem.Difficulty)
	fmt.Printf("Description: %s\n\n", problem.Description)

	l1 := structures.Ints2ListNode([]int{2, 4, 3})
	l2 := structures.Ints2ListNode([]int{5, 6, 4})

	fmt.Printf("Input: l1 = %v, l2 = %v\n", structures.ListNode2Ints(l1), structures.ListNode2Ints(l2))

	startTime := time.Now()
	result, err := s.Solve(2, l1, l2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	executionTime := time.Since(startTime)

	output := result.Output.(*structures.ListNode)
	fmt.Printf("Output: %v\n", structures.ListNode2Ints(output))
	fmt.Printf("Execution time: %v\n\n", executionTime)

	expectedOutput := structures.Ints2ListNode([]int{7, 0, 8})
	result = s.VerifySolution(result, expectedOutput)

	if result.Correct {
		fmt.Println("Solution is correct!")
	} else {
		fmt.Printf("Solution is incorrect. Expected: %v, Got: %v\n", 
			structures.ListNode2Ints(expectedOutput), 
			structures.ListNode2Ints(output))
	}
}

func RunLongestSubstringExample() {
	fmt.Println("Running Longest Substring Example")
	fmt.Println("================================")

	s := solver.NewSolver()
	s.InitializeProblems()

	problem, err := s.GetProblem(3)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Problem: %s (ID: %d)\n", problem.Title, problem.ID)
	fmt.Printf("Difficulty: %s\n", problem.Difficulty)
	fmt.Printf("Description: %s\n\n", problem.Description)

	input := "abcabcbb"

	fmt.Printf("Input: s = %s\n", input)

	startTime := time.Now()
	result, err := s.Solve(3, input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	executionTime := time.Since(startTime)

	fmt.Printf("Output: %v\n", result.Output)
	fmt.Printf("Execution time: %v\n\n", executionTime)

	expectedOutput := 3
	result = s.VerifySolution(result, expectedOutput)

	if result.Correct {
		fmt.Println("Solution is correct!")
	} else {
		fmt.Printf("Solution is incorrect. Expected: %v, Got: %v\n", expectedOutput, result.Output)
	}
}

func RunAllExamples() {
	RunTwoSumExample()
	fmt.Println()
	RunAddTwoNumbersExample()
	fmt.Println()
	RunLongestSubstringExample()
}
