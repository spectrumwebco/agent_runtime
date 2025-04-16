package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/leetcode/solver"
	"github.com/spectrumwebco/agent_runtime/pkg/modules/leetcode"
)

func NewLeetCodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "leetcode",
		Short: "LeetCode problem solver",
		Long:  `LeetCode problem solver for algorithm practice and interview preparation.`,
	}

	cmd.AddCommand(newLeetCodeListCommand())
	cmd.AddCommand(newLeetCodeSolveCommand())
	cmd.AddCommand(newLeetCodeExampleCommand())
	cmd.AddCommand(newLeetCodeCategoryCommand())
	cmd.AddCommand(newLeetCodeDifficultyCommand())

	return cmd
}

func newLeetCodeListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available LeetCode problems",
		Long:  `List all available LeetCode problems with their IDs, titles, and difficulties.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := leetcode.NewModule(cfg)
			solver := module.GetSolver()

			fmt.Println("Available LeetCode problems:")
			fmt.Println("----------------------------")

			for i := 1; i <= 5; i++ {
				problem, err := solver.GetProblem(i)
				if err != nil {
					continue
				}

				fmt.Printf("%d. %s (Difficulty: %s)\n", problem.ID, problem.Title, problem.Difficulty)
				fmt.Printf("   Categories: %s\n", formatCategories(problem.Categories))
				fmt.Printf("   Description: %s\n\n", problem.Description)
			}

			return nil
		},
	}

	return cmd
}

func newLeetCodeSolveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "solve [problem_id] [inputs...]",
		Short: "Solve a LeetCode problem",
		Long:  `Solve a LeetCode problem with the given inputs.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := leetcode.NewModule(cfg)

			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid problem ID: %s", args[0])
			}

			problem, err := module.GetProblem(id)
			if err != nil {
				return err
			}

			fmt.Printf("Solving problem: %s (ID: %d)\n", problem.Title, problem.ID)
			fmt.Printf("Difficulty: %s\n", problem.Difficulty)
			fmt.Printf("Description: %s\n\n", problem.Description)

			var inputs []interface{}
			var expectedOutput interface{}

			switch id {
			case 1: // Two Sum
				if len(args) < 3 {
					return fmt.Errorf("two sum problem requires nums array and target")
				}

				numsStr := strings.Trim(args[1], "[]")
				numsParts := strings.Split(numsStr, ",")
				nums := make([]int, len(numsParts))
				for i, part := range numsParts {
					num, err := strconv.Atoi(strings.TrimSpace(part))
					if err != nil {
						return fmt.Errorf("invalid number in nums array: %s", part)
					}
					nums[i] = num
				}

				target, err := strconv.Atoi(args[2])
				if err != nil {
					return fmt.Errorf("invalid target: %s", args[2])
				}

				inputs = append(inputs, nums, target)
				expectedOutput = []int{0, 1} // Default expected output for example case

				if len(args) > 3 {
					expectedStr := strings.Trim(args[3], "[]")
					expectedParts := strings.Split(expectedStr, ",")
					expected := make([]int, len(expectedParts))
					for i, part := range expectedParts {
						num, err := strconv.Atoi(strings.TrimSpace(part))
						if err != nil {
							return fmt.Errorf("invalid number in expected output: %s", part)
						}
						expected[i] = num
					}
					expectedOutput = expected
				}

			case 3: // Longest Substring Without Repeating Characters
				if len(args) < 2 {
					return fmt.Errorf("longest substring problem requires a string input")
				}

				s := args[1]
				inputs = append(inputs, s)
				expectedOutput = 3 // Default expected output for example case

				if len(args) > 2 {
					expected, err := strconv.Atoi(args[2])
					if err != nil {
						return fmt.Errorf("invalid expected output: %s", args[2])
					}
					expectedOutput = expected
				}

			default:
				return fmt.Errorf("solving problem %d is not yet implemented in the CLI", id)
			}

			fmt.Printf("Input: %v\n", inputs)
			startTime := time.Now()
			result, err := module.SolveProblem(id, inputs...)
			if err != nil {
				return fmt.Errorf("error solving problem: %v", err)
			}
			executionTime := time.Since(startTime)

			fmt.Printf("Output: %v\n", result.Output)
			fmt.Printf("Execution time: %v\n\n", executionTime)

			result = module.VerifySolution(result, expectedOutput)

			if result.Correct {
				fmt.Println("Solution is correct!")
			} else {
				fmt.Printf("Solution is incorrect. Expected: %v, Got: %v\n", expectedOutput, result.Output)
			}

			return nil
		},
	}

	return cmd
}

func newLeetCodeExampleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "example",
		Short: "Run LeetCode solver examples",
		Long:  `Run examples to demonstrate the LeetCode solver functionality.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := leetcode.NewModule(cfg)
			module.RunExample()

			return nil
		},
	}

	return cmd
}

func newLeetCodeCategoryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "category [category_name]",
		Short: "List problems by category",
		Long:  `List all LeetCode problems in a specific category.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := leetcode.NewModule(cfg)

			categoryName := args[0]
			var category solver.ProblemCategory

			switch strings.ToLower(categoryName) {
			case "array":
				category = solver.Array
			case "string":
				category = solver.String
			case "twopointers":
				category = solver.TwoPointers
			case "linkedlist":
				category = solver.LinkedList
			case "stack":
				category = solver.Stack
			case "tree":
				category = solver.Tree
			case "dynamicprogramming":
				category = solver.DynamicProgramming
			case "backtracking":
				category = solver.Backtracking
			case "depthfirstsearch":
				category = solver.DepthFirstSearch
			case "breadthfirstsearch":
				category = solver.BreadthFirstSearch
			case "binarysearch":
				category = solver.BinarySearch
			case "math":
				category = solver.Math
			case "hashtable":
				category = solver.HashTable
			case "sort":
				category = solver.Sort
			case "bitmanipulation":
				category = solver.BitManipulation
			case "unionfind":
				category = solver.UnionFind
			case "slidingwindow":
				category = solver.SlidingWindow
			case "segmenttree":
				category = solver.SegmentTree
			case "binaryindexedtree":
				category = solver.BinaryIndexedTree
			default:
				return fmt.Errorf("unknown category: %s", categoryName)
			}

			problems := module.GetProblemsByCategory(category)

			fmt.Printf("LeetCode problems in category '%s':\n", categoryName)
			fmt.Println("-----------------------------------")

			if len(problems) == 0 {
				fmt.Println("No problems found in this category.")
				return nil
			}

			for _, problem := range problems {
				fmt.Printf("%d. %s (Difficulty: %s)\n", problem.ID, problem.Title, problem.Difficulty)
				fmt.Printf("   Description: %s\n\n", problem.Description)
			}

			return nil
		},
	}

	return cmd
}

func newLeetCodeDifficultyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "difficulty [difficulty_level]",
		Short: "List problems by difficulty",
		Long:  `List all LeetCode problems with a specific difficulty level (easy, medium, hard).`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := leetcode.NewModule(cfg)

			difficultyName := args[0]
			var difficulty solver.Difficulty

			switch strings.ToLower(difficultyName) {
			case "easy":
				difficulty = solver.Easy
			case "medium":
				difficulty = solver.Medium
			case "hard":
				difficulty = solver.Hard
			default:
				return fmt.Errorf("unknown difficulty: %s (must be easy, medium, or hard)", difficultyName)
			}

			problems := module.GetProblemsByDifficulty(difficulty)

			fmt.Printf("LeetCode problems with difficulty '%s':\n", difficultyName)
			fmt.Println("---------------------------------------")

			if len(problems) == 0 {
				fmt.Println("No problems found with this difficulty.")
				return nil
			}

			for _, problem := range problems {
				fmt.Printf("%d. %s\n", problem.ID, problem.Title)
				fmt.Printf("   Categories: %s\n", formatCategories(problem.Categories))
				fmt.Printf("   Description: %s\n\n", problem.Description)
			}

			return nil
		},
	}

	return cmd
}

func formatCategories(categories []solver.ProblemCategory) string {
	var categoryStrings []string
	for _, category := range categories {
		categoryStrings = append(categoryStrings, string(category))
	}
	return strings.Join(categoryStrings, ", ")
}
