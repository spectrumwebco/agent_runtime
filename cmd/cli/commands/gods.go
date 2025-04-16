package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/datastructures/gods"
	godsModule "github.com/spectrumwebco/agent_runtime/pkg/modules/gods"
)

func NewGodsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gods",
		Short: "GoDS data structures and algorithms",
		Long:  `GoDS provides implementations of various data structures and algorithms in Go.`,
	}

	cmd.AddCommand(newGodsListCommand())
	cmd.AddCommand(newGodsInfoCommand())
	cmd.AddCommand(newGodsExampleCommand())
	cmd.AddCommand(newGodsCreateCommand())

	return cmd
}

func newGodsListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available data structures",
		Long:  `List all available data structures provided by GoDS.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := godsModule.NewModule(cfg)
			structures := module.ListAvailableStructures()

			fmt.Println("Available Data Structures:")
			fmt.Println("==========================")
			
			for category, items := range structures {
				fmt.Printf("\n%s:\n", category)
				for _, item := range items {
					fmt.Printf("  - %s\n", item)
					desc := module.GetStructureDescription(item)
					fmt.Printf("    %s\n", desc)
				}
			}

			return nil
		},
	}

	return cmd
}

func newGodsInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info [structure]",
		Short: "Get information about a data structure",
		Long:  `Get detailed information about a specific data structure, including description and use cases.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := godsModule.NewModule(cfg)
			structureName := strings.ToLower(args[0])
			
			description := module.GetStructureDescription(structureName)
			useCases := module.GetStructureUseCases(structureName)
			
			fmt.Printf("Data Structure: %s\n", structureName)
			fmt.Printf("Description: %s\n", description)
			
			fmt.Println("\nUse Cases:")
			for _, useCase := range useCases {
				fmt.Printf("  - %s\n", useCase)
			}
			
			return nil
		},
	}

	return cmd
}

func newGodsExampleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "example",
		Short: "Run a GoDS example",
		Long:  `Run a simple example demonstrating the usage of GoDS data structures.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := godsModule.NewModule(cfg)
			module.RunExample()

			return nil
		},
	}

	return cmd
}

func newGodsCreateCommand() *cobra.Command {
	var (
		dataType     string
		structureType string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a data structure",
		Long:  `Create a specific data structure and perform operations on it.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := godsModule.NewModule(cfg)
			
			dataType = strings.ToLower(dataType)
			structureType = strings.ToLower(structureType)
			
			switch dataType {
			case "list":
				return handleListCreation(module, gods.ListType(structureType))
			case "map":
				return handleMapCreation(module, gods.MapType(structureType))
			case "set":
				return handleSetCreation(module, gods.SetType(structureType))
			case "stack":
				return handleStackCreation(module, gods.StackType(structureType))
			case "queue":
				return handleQueueCreation(module, gods.QueueType(structureType))
			default:
				return fmt.Errorf("unknown data type: %s", dataType)
			}
		},
	}

	cmd.Flags().StringVarP(&dataType, "type", "t", "", "Type of data structure (list, map, set, stack, queue)")
	cmd.Flags().StringVarP(&structureType, "structure", "s", "", "Specific structure implementation")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("structure")

	return cmd
}

func handleListCreation(module *godsModule.Module, listType gods.ListType) error {
	list, err := module.CreateList(listType)
	if err != nil {
		return err
	}
	
	stringList, ok := list.(interface {
		Add(values ...interface{})
		Remove(index int)
		Contains(values ...interface{}) bool
		Get(index int) (interface{}, bool)
		Empty() bool
		Size() int
		Clear()
		Values() []interface{}
		String() string
	})
	if !ok {
		return fmt.Errorf("could not type assert list")
	}
	
	stringList.Add("apple", "banana", "cherry", "date", "elderberry")
	fmt.Printf("Created list with 5 elements: %s\n", stringList.String())
	
	fmt.Printf("Contains 'banana': %v\n", stringList.Contains("banana"))
	
	if val, found := stringList.Get(2); found {
		fmt.Printf("Element at index 2: %v\n", val)
	}
	
	stringList.Remove(1)
	fmt.Printf("After removing element at index 1: %s\n", stringList.String())
	
	fmt.Printf("Size: %d\n", stringList.Size())
	
	fmt.Printf("Values: %v\n", stringList.Values())
	
	return nil
}

func handleMapCreation(module *godsModule.Module, mapType gods.MapType) error {
	m, err := module.CreateMap(mapType)
	if err != nil {
		return err
	}
	
	stringMap, ok := m.(interface {
		Put(key, value interface{})
		Get(key interface{}) (interface{}, bool)
		Remove(key interface{})
		Empty() bool
		Size() int
		Keys() []interface{}
		Values() []interface{}
		Clear()
		String() string
	})
	if !ok {
		return fmt.Errorf("could not type assert map")
	}
	
	stringMap.Put("one", "first")
	stringMap.Put("two", "second")
	stringMap.Put("three", "third")
	stringMap.Put("four", "fourth")
	stringMap.Put("five", "fifth")
	fmt.Printf("Created map with 5 elements: %s\n", stringMap.String())
	
	if val, found := stringMap.Get("three"); found {
		fmt.Printf("Value for key 'three': %v\n", val)
	}
	
	stringMap.Remove("two")
	fmt.Printf("After removing key 'two': %s\n", stringMap.String())
	
	fmt.Printf("Size: %d\n", stringMap.Size())
	
	fmt.Printf("Keys: %v\n", stringMap.Keys())
	fmt.Printf("Values: %v\n", stringMap.Values())
	
	return nil
}

func handleSetCreation(module *godsModule.Module, setType gods.SetType) error {
	set, err := module.CreateSet(setType)
	if err != nil {
		return err
	}
	
	stringSet, ok := set.(interface {
		Add(items ...interface{})
		Remove(items ...interface{})
		Contains(items ...interface{}) bool
		Empty() bool
		Size() int
		Clear()
		Values() []interface{}
		String() string
	})
	if !ok {
		return fmt.Errorf("could not type assert set")
	}
	
	stringSet.Add("apple", "banana", "cherry", "apple", "date")
	fmt.Printf("Created set with unique elements: %s\n", stringSet.String())
	
	fmt.Printf("Contains 'banana': %v\n", stringSet.Contains("banana"))
	
	stringSet.Remove("cherry")
	fmt.Printf("After removing 'cherry': %s\n", stringSet.String())
	
	fmt.Printf("Size: %d\n", stringSet.Size())
	
	fmt.Printf("Values: %v\n", stringSet.Values())
	
	return nil
}

func handleStackCreation(module *godsModule.Module, stackType gods.StackType) error {
	stack, err := module.CreateStack(stackType)
	if err != nil {
		return err
	}
	
	stringStack, ok := stack.(interface {
		Push(value interface{})
		Pop() (interface{}, bool)
		Peek() (interface{}, bool)
		Empty() bool
		Size() int
		Clear()
		Values() []interface{}
		String() string
	})
	if !ok {
		return fmt.Errorf("could not type assert stack")
	}
	
	stringStack.Push("first")
	stringStack.Push("second")
	stringStack.Push("third")
	fmt.Printf("Created stack with 3 elements: %s\n", stringStack.String())
	
	if val, found := stringStack.Peek(); found {
		fmt.Printf("Top element: %v\n", val)
	}
	
	if val, found := stringStack.Pop(); found {
		fmt.Printf("Popped element: %v\n", val)
	}
	
	fmt.Printf("Size after pop: %d\n", stringStack.Size())
	
	fmt.Printf("Remaining values: %v\n", stringStack.Values())
	
	return nil
}

func handleQueueCreation(module *godsModule.Module, queueType gods.QueueType) error {
	queue, err := module.CreateQueue(queueType)
	if err != nil {
		return err
	}
	
	stringQueue, ok := queue.(interface {
		Enqueue(value interface{})
		Dequeue() (interface{}, bool)
		Peek() (interface{}, bool)
		Empty() bool
		Size() int
		Clear()
		Values() []interface{}
		String() string
	})
	if !ok {
		return fmt.Errorf("could not type assert queue")
	}
	
	stringQueue.Enqueue("first")
	stringQueue.Enqueue("second")
	stringQueue.Enqueue("third")
	fmt.Printf("Created queue with 3 elements: %s\n", stringQueue.String())
	
	if val, found := stringQueue.Peek(); found {
		fmt.Printf("Front element: %v\n", val)
	}
	
	if val, found := stringQueue.Dequeue(); found {
		fmt.Printf("Dequeued element: %v\n", val)
	}
	
	fmt.Printf("Size after dequeue: %d\n", stringQueue.Size())
	
	fmt.Printf("Remaining values: %v\n", stringQueue.Values())
	
	return nil
}
