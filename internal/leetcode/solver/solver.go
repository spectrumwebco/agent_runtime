package solver

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/spectrumwebco/agent_runtime/internal/leetcode/structures"
)

type ProblemCategory string

const (
	Array            ProblemCategory = "Array"
	String           ProblemCategory = "String"
	TwoPointers      ProblemCategory = "TwoPointers"
	LinkedList       ProblemCategory = "LinkedList"
	Stack            ProblemCategory = "Stack"
	Tree             ProblemCategory = "Tree"
	DynamicProgramming ProblemCategory = "DynamicProgramming"
	Backtracking     ProblemCategory = "Backtracking"
	DepthFirstSearch ProblemCategory = "DepthFirstSearch"
	BreadthFirstSearch ProblemCategory = "BreadthFirstSearch"
	BinarySearch     ProblemCategory = "BinarySearch"
	Math             ProblemCategory = "Math"
	HashTable        ProblemCategory = "HashTable"
	Sort             ProblemCategory = "Sort"
	BitManipulation  ProblemCategory = "BitManipulation"
	UnionFind        ProblemCategory = "UnionFind"
	SlidingWindow    ProblemCategory = "SlidingWindow"
	SegmentTree      ProblemCategory = "SegmentTree"
	BinaryIndexedTree ProblemCategory = "BinaryIndexedTree"
)

type Difficulty string

const (
	Easy   Difficulty = "Easy"
	Medium Difficulty = "Medium"
	Hard   Difficulty = "Hard"
)

type Problem struct {
	ID          int
	Title       string
	Description string
	Difficulty  Difficulty
	Categories  []ProblemCategory
	Solution    interface{}
}

type SolverResult struct {
	Problem      *Problem
	Input        []interface{}
	Output       interface{}
	ExpectedOutput interface{}
	Correct      bool
	ExecutionTime time.Duration
}

type Solver struct {
	problems map[int]*Problem
}

func NewSolver() *Solver {
	return &Solver{
		problems: make(map[int]*Problem),
	}
}

func (s *Solver) RegisterProblem(problem *Problem) {
	s.problems[problem.ID] = problem
}

func (s *Solver) GetProblem(id int) (*Problem, error) {
	problem, exists := s.problems[id]
	if !exists {
		return nil, fmt.Errorf("problem with ID %d not found", id)
	}
	return problem, nil
}

func (s *Solver) GetProblemsByCategory(category ProblemCategory) []*Problem {
	var result []*Problem
	for _, problem := range s.problems {
		for _, cat := range problem.Categories {
			if cat == category {
				result = append(result, problem)
				break
			}
		}
	}
	return result
}

func (s *Solver) GetProblemsByDifficulty(difficulty Difficulty) []*Problem {
	var result []*Problem
	for _, problem := range s.problems {
		if problem.Difficulty == difficulty {
			result = append(result, problem)
		}
	}
	return result
}

func (s *Solver) Solve(id int, inputs ...interface{}) (*SolverResult, error) {
	problem, err := s.GetProblem(id)
	if err != nil {
		return nil, err
	}

	if problem.Solution == nil {
		return nil, errors.New("problem has no solution function")
	}

	solutionValue := reflect.ValueOf(problem.Solution)
	if solutionValue.Kind() != reflect.Func {
		return nil, errors.New("solution is not a function")
	}

	solutionType := solutionValue.Type()
	if solutionType.NumIn() != len(inputs) {
		return nil, fmt.Errorf("solution function expects %d arguments, but got %d", solutionType.NumIn(), len(inputs))
	}

	args := make([]reflect.Value, len(inputs))
	for i, input := range inputs {
		inputValue := reflect.ValueOf(input)
		expectedType := solutionType.In(i)
		
		if !inputValue.Type().AssignableTo(expectedType) {
			if expectedType.Kind() == reflect.Ptr && expectedType.Elem().Name() == "TreeNode" {
				if inputValue.Kind() == reflect.Slice && inputValue.Type().Elem().Kind() == reflect.Int {
					intSlice := make([]int, inputValue.Len())
					for j := 0; j < inputValue.Len(); j++ {
						intSlice[j] = int(inputValue.Index(j).Int())
					}
					treeNode := structures.Ints2TreeNode(intSlice)
					args[i] = reflect.ValueOf(treeNode)
					continue
				}
			} else if expectedType.Kind() == reflect.Ptr && expectedType.Elem().Name() == "ListNode" {
				if inputValue.Kind() == reflect.Slice && inputValue.Type().Elem().Kind() == reflect.Int {
					intSlice := make([]int, inputValue.Len())
					for j := 0; j < inputValue.Len(); j++ {
						intSlice[j] = int(inputValue.Index(j).Int())
					}
					listNode := structures.Ints2ListNode(intSlice)
					args[i] = reflect.ValueOf(listNode)
					continue
				}
			}
			
			return nil, fmt.Errorf("argument %d: cannot convert %v to %v", i, inputValue.Type(), expectedType)
		}
		
		args[i] = inputValue
	}

	startTime := time.Now()
	results := solutionValue.Call(args)
	executionTime := time.Since(startTime)

	var output interface{}
	if len(results) > 0 {
		output = results[0].Interface()
	}

	return &SolverResult{
		Problem:      problem,
		Input:        inputs,
		Output:       output,
		ExecutionTime: executionTime,
	}, nil
}

func (s *Solver) VerifySolution(result *SolverResult, expectedOutput interface{}) *SolverResult {
	result.ExpectedOutput = expectedOutput
	
	outputValue := reflect.ValueOf(result.Output)
	expectedValue := reflect.ValueOf(expectedOutput)
	
	if outputValue.Kind() == reflect.Ptr && outputValue.Type().Elem().Name() == "TreeNode" {
		outputInts := structures.TreeNode2Ints(outputValue.Interface().(*structures.TreeNode))
		
		if expectedValue.Kind() == reflect.Slice && expectedValue.Type().Elem().Kind() == reflect.Int {
			expectedInts := make([]int, expectedValue.Len())
			for i := 0; i < expectedValue.Len(); i++ {
				expectedInts[i] = int(expectedValue.Index(i).Int())
			}
			
			result.Correct = reflect.DeepEqual(outputInts, expectedInts)
		} else if expectedValue.Kind() == reflect.Ptr && expectedValue.Type().Elem().Name() == "TreeNode" {
			expectedInts := structures.TreeNode2Ints(expectedValue.Interface().(*structures.TreeNode))
			result.Correct = reflect.DeepEqual(outputInts, expectedInts)
		}
	} else if outputValue.Kind() == reflect.Ptr && outputValue.Type().Elem().Name() == "ListNode" {
		outputInts := structures.ListNode2Ints(outputValue.Interface().(*structures.ListNode))
		
		if expectedValue.Kind() == reflect.Slice && expectedValue.Type().Elem().Kind() == reflect.Int {
			expectedInts := make([]int, expectedValue.Len())
			for i := 0; i < expectedValue.Len(); i++ {
				expectedInts[i] = int(expectedValue.Index(i).Int())
			}
			
			result.Correct = reflect.DeepEqual(outputInts, expectedInts)
		} else if expectedValue.Kind() == reflect.Ptr && expectedValue.Type().Elem().Name() == "ListNode" {
			expectedInts := structures.ListNode2Ints(expectedValue.Interface().(*structures.ListNode))
			result.Correct = reflect.DeepEqual(outputInts, expectedInts)
		}
	} else {
		result.Correct = reflect.DeepEqual(result.Output, expectedOutput)
	}
	
	return result
}

func (s *Solver) InitializeProblems() {
	s.RegisterProblem(&Problem{
		ID:          1,
		Title:       "Two Sum",
		Description: "Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.",
		Difficulty:  Easy,
		Categories:  []ProblemCategory{Array, HashTable},
		Solution:    twoSum,
	})
	
	s.RegisterProblem(&Problem{
		ID:          2,
		Title:       "Add Two Numbers",
		Description: "You are given two non-empty linked lists representing two non-negative integers. The digits are stored in reverse order, and each of their nodes contains a single digit. Add the two numbers and return the sum as a linked list.",
		Difficulty:  Medium,
		Categories:  []ProblemCategory{LinkedList, Math},
		Solution:    addTwoNumbers,
	})
	
	s.RegisterProblem(&Problem{
		ID:          3,
		Title:       "Longest Substring Without Repeating Characters",
		Description: "Given a string s, find the length of the longest substring without repeating characters.",
		Difficulty:  Medium,
		Categories:  []ProblemCategory{String, HashTable, SlidingWindow},
		Solution:    lengthOfLongestSubstring,
	})
	
	s.RegisterProblem(&Problem{
		ID:          4,
		Title:       "Median of Two Sorted Arrays",
		Description: "Given two sorted arrays nums1 and nums2 of size m and n respectively, return the median of the two sorted arrays.",
		Difficulty:  Hard,
		Categories:  []ProblemCategory{Array, BinarySearch, DivideAndConquer},
		Solution:    findMedianSortedArrays,
	})
	
	s.RegisterProblem(&Problem{
		ID:          5,
		Title:       "Longest Palindromic Substring",
		Description: "Given a string s, return the longest palindromic substring in s.",
		Difficulty:  Medium,
		Categories:  []ProblemCategory{String, DynamicProgramming},
		Solution:    longestPalindrome,
	})
}


func twoSum(nums []int, target int) []int {
	m := make(map[int]int)
	for k, v := range nums {
		if idx, ok := m[target-v]; ok {
			return []int{idx, k}
		}
		m[v] = k
	}
	return nil
}

func addTwoNumbers(l1 *structures.ListNode, l2 *structures.ListNode) *structures.ListNode {
	dummy := &structures.ListNode{}
	curr := dummy
	carry := 0
	
	for l1 != nil || l2 != nil || carry > 0 {
		sum := carry
		
		if l1 != nil {
			sum += l1.Val
			l1 = l1.Next
		}
		
		if l2 != nil {
			sum += l2.Val
			l2 = l2.Next
		}
		
		carry = sum / 10
		curr.Next = &structures.ListNode{Val: sum % 10}
		curr = curr.Next
	}
	
	return dummy.Next
}

func lengthOfLongestSubstring(s string) int {
	charMap := make(map[byte]int)
	maxLength := 0
	start := 0
	
	for i := 0; i < len(s); i++ {
		if idx, found := charMap[s[i]]; found && idx >= start {
			start = idx + 1
		}
		
		charMap[s[i]] = i
		currLength := i - start + 1
		if currLength > maxLength {
			maxLength = currLength
		}
	}
	
	return maxLength
}

func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	if len(nums1) > len(nums2) {
		nums1, nums2 = nums2, nums1
	}
	
	x, y := len(nums1), len(nums2)
	low, high := 0, x
	
	for low <= high {
		partitionX := (low + high) / 2
		partitionY := (x + y + 1) / 2 - partitionX
		
		maxX := getMax(nums1, partitionX)
		minX := getMin(nums1, partitionX)
		
		maxY := getMax(nums2, partitionY)
		minY := getMin(nums2, partitionY)
		
		if maxX <= minY && maxY <= minX {
			if (x+y)%2 == 0 {
				return float64(max(maxX, maxY)+min(minX, minY)) / 2.0
			}
			return float64(max(maxX, maxY))
		} else if maxX > minY {
			high = partitionX - 1
		} else {
			low = partitionX + 1
		}
	}
	
	return 0.0
}

func getMax(nums []int, partition int) int {
	if partition == 0 {
		return -1 << 31
	}
	return nums[partition-1]
}

func getMin(nums []int, partition int) int {
	if partition == len(nums) {
		return 1 << 31
	}
	return nums[partition]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func longestPalindrome(s string) string {
	if len(s) < 2 {
		return s
	}
	
	start, maxLen := 0, 1
	
	for i := 0; i < len(s); {
		if len(s)-i <= maxLen/2 {
			break
		}
		
		j, k := i, i
		for k < len(s)-1 && s[k+1] == s[k] {
			k++
		}
		
		i = k + 1
		
		for k < len(s)-1 && j > 0 && s[k+1] == s[j-1] {
			k++
			j--
		}
		
		newLen := k - j + 1
		if newLen > maxLen {
			start = j
			maxLen = newLen
		}
	}
	
	return s[start : start+maxLen]
}

const DivideAndConquer ProblemCategory = "DivideAndConquer"
