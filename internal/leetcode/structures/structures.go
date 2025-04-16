package structures

import (
	"fmt"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

type ListNode struct {
	Val  int
	Next *ListNode
}

var NULL = -1 << 63

func Ints2TreeNode(ints []int) *TreeNode {
	n := len(ints)
	if n == 0 {
		return nil
	}

	root := &TreeNode{
		Val: ints[0],
	}

	queue := make([]*TreeNode, 1, n*2)
	queue[0] = root

	i := 1
	for i < n {
		node := queue[0]
		queue = queue[1:]

		if i < n && ints[i] != NULL {
			node.Left = &TreeNode{Val: ints[i]}
			queue = append(queue, node.Left)
		}
		i++

		if i < n && ints[i] != NULL {
			node.Right = &TreeNode{Val: ints[i]}
			queue = append(queue, node.Right)
		}
		i++
	}

	return root
}

func TreeNode2Ints(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}

	res := make([]int, 0)
	queue := []*TreeNode{root}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if node == nil {
			res = append(res, NULL)
		} else {
			res = append(res, node.Val)
			queue = append(queue, node.Left, node.Right)
		}
	}

	for i := len(res) - 1; i >= 0; i-- {
		if res[i] != NULL {
			return res[:i+1]
		}
	}

	return res
}

func Ints2ListNode(ints []int) *ListNode {
	if len(ints) == 0 {
		return nil
	}

	head := &ListNode{Val: ints[0]}
	current := head

	for i := 1; i < len(ints); i++ {
		current.Next = &ListNode{Val: ints[i]}
		current = current.Next
	}

	return head
}

func ListNode2Ints(head *ListNode) []int {
	limit := 100
	times := 0

	res := []int{}
	for head != nil {
		times++
		if times > limit {
			msg := fmt.Sprintf("List depth exceeds %d, possible cycle detected", limit)
			panic(msg)
		}

		res = append(res, head.Val)
		head = head.Next
	}

	return res
}

type Stack struct {
	elements []interface{}
}

func NewStack() *Stack {
	return &Stack{elements: make([]interface{}, 0)}
}

func (s *Stack) Push(v interface{}) {
	s.elements = append(s.elements, v)
}

func (s *Stack) Pop() interface{} {
	if len(s.elements) == 0 {
		return nil
	}
	
	lastIndex := len(s.elements) - 1
	value := s.elements[lastIndex]
	s.elements = s.elements[:lastIndex]
	return value
}

func (s *Stack) Peek() interface{} {
	if len(s.elements) == 0 {
		return nil
	}
	
	return s.elements[len(s.elements)-1]
}

func (s *Stack) IsEmpty() bool {
	return len(s.elements) == 0
}

func (s *Stack) Size() int {
	return len(s.elements)
}

type Queue struct {
	elements []interface{}
}

func NewQueue() *Queue {
	return &Queue{elements: make([]interface{}, 0)}
}

func (q *Queue) Enqueue(v interface{}) {
	q.elements = append(q.elements, v)
}

func (q *Queue) Dequeue() interface{} {
	if len(q.elements) == 0 {
		return nil
	}
	
	value := q.elements[0]
	q.elements = q.elements[1:]
	return value
}

func (q *Queue) Front() interface{} {
	if len(q.elements) == 0 {
		return nil
	}
	
	return q.elements[0]
}

func (q *Queue) IsEmpty() bool {
	return len(q.elements) == 0
}

func (q *Queue) Size() int {
	return len(q.elements)
}
