package gods

import (
	"fmt"
	"strings"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/lists/doublylinkedlist"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/emirpasic/gods/maps/hashbidimap"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/emirpasic/gods/maps/treebidimap"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/queues/arrayqueue"
	"github.com/emirpasic/gods/queues/circularbuffer"
	"github.com/emirpasic/gods/queues/linkedlistqueue"
	"github.com/emirpasic/gods/queues/priorityqueue"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/emirpasic/gods/sets/linkedhashset"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/emirpasic/gods/stacks/arraystack"
	"github.com/emirpasic/gods/stacks/linkedliststack"
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/emirpasic/gods/trees/binaryheap"
	"github.com/emirpasic/gods/trees/btree"
	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/emirpasic/gods/utils"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Client struct {
	config *config.Config
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
	}
}

type ListType string

const (
	ArrayList         ListType = "arraylist"
	SinglyLinkedList  ListType = "singlylinkedlist"
	DoublyLinkedList  ListType = "doublylinkedlist"
)

type MapType string

const (
	HashMap        MapType = "hashmap"
	TreeMap        MapType = "treemap"
	LinkedHashMap  MapType = "linkedhashmap"
	HashBidiMap    MapType = "hashbidimap"
	TreeBidiMap    MapType = "treebidimap"
)

type SetType string

const (
	HashSet        SetType = "hashset"
	TreeSet        SetType = "treeset"
	LinkedHashSet  SetType = "linkedhashset"
)

type StackType string

const (
	ArrayStack        StackType = "arraystack"
	LinkedListStack   StackType = "linkedliststack"
)

type QueueType string

const (
	ArrayQueue        QueueType = "arrayqueue"
	LinkedListQueue   QueueType = "linkedlistqueue"
	CircularBuffer    QueueType = "circularbuffer"
	PriorityQueue     QueueType = "priorityqueue"
)

type TreeType string

const (
	RedBlackTree  TreeType = "redblacktree"
	AVLTree       TreeType = "avltree"
	BTree         TreeType = "btree"
	BinaryHeap    TreeType = "binaryheap"
)

func (c *Client) CreateList(listType ListType) (interface{}, error) {
	switch listType {
	case ArrayList:
		return arraylist.New(), nil
	case SinglyLinkedList:
		return singlylinkedlist.New(), nil
	case DoublyLinkedList:
		return doublylinkedlist.New(), nil
	default:
		return nil, fmt.Errorf("unknown list type: %s", listType)
	}
}

func (c *Client) CreateMap(mapType MapType) (interface{}, error) {
	switch mapType {
	case HashMap:
		return hashmap.New(), nil
	case TreeMap:
		return treemap.NewWithStringComparator(), nil
	case LinkedHashMap:
		return linkedhashmap.New(), nil
	case HashBidiMap:
		return hashbidimap.New(), nil
	case TreeBidiMap:
		return treebidimap.NewWithStringComparator(), nil
	default:
		return nil, fmt.Errorf("unknown map type: %s", mapType)
	}
}

func (c *Client) CreateSet(setType SetType) (interface{}, error) {
	switch setType {
	case HashSet:
		return hashset.New(), nil
	case TreeSet:
		return treeset.NewWithStringComparator(), nil
	case LinkedHashSet:
		return linkedhashset.New(), nil
	default:
		return nil, fmt.Errorf("unknown set type: %s", setType)
	}
}

func (c *Client) CreateStack(stackType StackType) (interface{}, error) {
	switch stackType {
	case ArrayStack:
		return arraystack.New(), nil
	case LinkedListStack:
		return linkedliststack.New(), nil
	default:
		return nil, fmt.Errorf("unknown stack type: %s", stackType)
	}
}

func (c *Client) CreateQueue(queueType QueueType) (interface{}, error) {
	switch queueType {
	case ArrayQueue:
		return arrayqueue.New(), nil
	case LinkedListQueue:
		return linkedlistqueue.New(), nil
	case CircularBuffer:
		return circularbuffer.New(10), nil // Default capacity of 10
	case PriorityQueue:
		return priorityqueue.New(utils.IntComparator), nil
	default:
		return nil, fmt.Errorf("unknown queue type: %s", queueType)
	}
}

func (c *Client) CreateTree(treeType TreeType) (interface{}, error) {
	switch treeType {
	case RedBlackTree:
		return redblacktree.NewWithStringComparator(), nil
	case AVLTree:
		return avltree.NewWithStringComparator(), nil
	case BTree:
		return btree.NewWithStringComparator(3), nil // Default order of 3
	case BinaryHeap:
		return binaryheap.NewWithIntComparator(), nil
	default:
		return nil, fmt.Errorf("unknown tree type: %s", treeType)
	}
}

func (c *Client) GetComparator(typeName string) utils.Comparator {
	switch strings.ToLower(typeName) {
	case "int":
		return utils.IntComparator
	case "string":
		return utils.StringComparator
	case "time":
		return utils.TimeComparator
	default:
		return utils.StringComparator
	}
}

func (c *Client) ListAvailableStructures() map[string][]string {
	structures := make(map[string][]string)
	
	structures["Lists"] = []string{
		string(ArrayList),
		string(SinglyLinkedList),
		string(DoublyLinkedList),
	}
	
	structures["Maps"] = []string{
		string(HashMap),
		string(TreeMap),
		string(LinkedHashMap),
		string(HashBidiMap),
		string(TreeBidiMap),
	}
	
	structures["Sets"] = []string{
		string(HashSet),
		string(TreeSet),
		string(LinkedHashSet),
	}
	
	structures["Stacks"] = []string{
		string(ArrayStack),
		string(LinkedListStack),
	}
	
	structures["Queues"] = []string{
		string(ArrayQueue),
		string(LinkedListQueue),
		string(CircularBuffer),
		string(PriorityQueue),
	}
	
	structures["Trees"] = []string{
		string(RedBlackTree),
		string(AVLTree),
		string(BTree),
		string(BinaryHeap),
	}
	
	return structures
}

func (c *Client) GetStructureDescription(structureType string) string {
	descriptions := map[string]string{
		string(ArrayList):         "Dynamic array implementation of the list interface",
		string(SinglyLinkedList):  "Linked list implementation of the list interface with single links",
		string(DoublyLinkedList):  "Linked list implementation of the list interface with double links",
		string(HashMap):           "Map implementation using hash tables",
		string(TreeMap):           "Map implementation using red-black tree",
		string(LinkedHashMap):     "Map implementation that preserves insertion order",
		string(HashBidiMap):       "Bidirectional map implementation using hash tables",
		string(TreeBidiMap):       "Bidirectional map implementation using red-black tree",
		string(HashSet):           "Set implementation using hash tables",
		string(TreeSet):           "Set implementation using red-black tree",
		string(LinkedHashSet):     "Set implementation that preserves insertion order",
		string(ArrayStack):        "Stack implementation using dynamic array",
		string(LinkedListStack):   "Stack implementation using linked list",
		string(ArrayQueue):        "Queue implementation using dynamic array",
		string(LinkedListQueue):   "Queue implementation using linked list",
		string(CircularBuffer):    "Queue implementation with fixed size circular buffer",
		string(PriorityQueue):     "Queue implementation where elements are ordered by priority",
		string(RedBlackTree):      "Self-balancing binary search tree with red-black properties",
		string(AVLTree):           "Self-balancing binary search tree with height-balanced property",
		string(BTree):             "Self-balancing tree data structure that maintains sorted data",
		string(BinaryHeap):        "Tree-based data structure that satisfies the heap property",
	}
	
	if desc, ok := descriptions[strings.ToLower(structureType)]; ok {
		return desc
	}
	
	return "Unknown data structure type"
}

func (c *Client) GetStructureUseCases(structureType string) []string {
	useCases := map[string][]string{
		string(ArrayList): {
			"Fast random access to elements",
			"Efficient appending of elements",
			"When frequent modifications are not needed",
		},
		string(SinglyLinkedList): {
			"Efficient insertion and removal at the beginning",
			"Memory-efficient for large datasets",
			"When random access is not required",
		},
		string(DoublyLinkedList): {
			"Efficient insertion and removal at both ends",
			"Bidirectional traversal",
			"Implementation of deques",
		},
		string(HashMap): {
			"Fast key-value lookups",
			"Caching data",
			"Counting occurrences",
		},
		string(TreeMap): {
			"Ordered key-value pairs",
			"Range queries",
			"When keys need to be sorted",
		},
		string(LinkedHashMap): {
			"LRU caches",
			"Maintaining insertion order",
			"Predictable iteration order",
		},
		string(HashBidiMap): {
			"Two-way lookups",
			"Inverse mappings",
			"When values are unique",
		},
		string(TreeBidiMap): {
			"Ordered two-way lookups",
			"Sorted inverse mappings",
			"When values are unique and order matters",
		},
		string(HashSet): {
			"Removing duplicates",
			"Fast membership testing",
			"Set operations (union, intersection)",
		},
		string(TreeSet): {
			"Ordered set operations",
			"Range queries on sets",
			"When elements need to be sorted",
		},
		string(LinkedHashSet): {
			"Ordered duplicate removal",
			"Predictable iteration order",
			"When insertion order matters",
		},
		string(ArrayStack): {
			"Function call stack",
			"Expression evaluation",
			"Backtracking algorithms",
		},
		string(LinkedListStack): {
			"Memory-efficient stack implementation",
			"Undo functionality",
			"Syntax parsing",
		},
		string(ArrayQueue): {
			"Job scheduling",
			"Breadth-first search",
			"Buffering",
		},
		string(LinkedListQueue): {
			"Memory-efficient queue implementation",
			"Message queues",
			"Task scheduling",
		},
		string(CircularBuffer): {
			"Fixed-size buffers",
			"Streaming data processing",
			"Implementing ring buffers",
		},
		string(PriorityQueue): {
			"Task scheduling by priority",
			"Dijkstra's algorithm",
			"Event-driven simulation",
		},
		string(RedBlackTree): {
			"Implementing maps and sets",
			"Database indexing",
			"Computational geometry",
		},
		string(AVLTree): {
			"Strict balancing requirements",
			"Frequent lookups",
			"When tree height matters",
		},
		string(BTree): {
			"Database and file system indexing",
			"Multi-level indexing",
			"When working with disk storage",
		},
		string(BinaryHeap): {
			"Priority queues",
			"Heap sort",
			"Graph algorithms like Prim's and Dijkstra's",
		},
	}
	
	if cases, ok := useCases[strings.ToLower(structureType)]; ok {
		return cases
	}
	
	return []string{"No specific use cases available"}
}
