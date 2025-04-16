package gods

import (
	"context"
	"fmt"

	"github.com/spectrumwebco/agent_runtime/internal/datastructures/gods"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Module struct {
	config *config.Config
	client *gods.Client
}

func NewModule(cfg *config.Config) *Module {
	return &Module{
		config: cfg,
		client: gods.NewClient(cfg),
	}
}

func (m *Module) Name() string {
	return "gods"
}

func (m *Module) Description() string {
	return "Implementation of various data structures and algorithms in Go"
}

func (m *Module) Initialize(ctx context.Context) error {
	return nil
}

func (m *Module) Shutdown(ctx context.Context) error {
	return nil
}

func (m *Module) GetClient() *gods.Client {
	return m.client
}

func (m *Module) ListAvailableStructures() map[string][]string {
	return m.client.ListAvailableStructures()
}

func (m *Module) GetStructureDescription(structureType string) string {
	return m.client.GetStructureDescription(structureType)
}

func (m *Module) GetStructureUseCases(structureType string) []string {
	return m.client.GetStructureUseCases(structureType)
}

func (m *Module) CreateList(listType gods.ListType) (interface{}, error) {
	return m.client.CreateList(listType)
}

func (m *Module) CreateMap(mapType gods.MapType) (interface{}, error) {
	return m.client.CreateMap(mapType)
}

func (m *Module) CreateSet(setType gods.SetType) (interface{}, error) {
	return m.client.CreateSet(setType)
}

func (m *Module) CreateStack(stackType gods.StackType) (interface{}, error) {
	return m.client.CreateStack(stackType)
}

func (m *Module) CreateQueue(queueType gods.QueueType) (interface{}, error) {
	return m.client.CreateQueue(queueType)
}

func (m *Module) CreateTree(treeType gods.TreeType) (interface{}, error) {
	return m.client.CreateTree(treeType)
}

func (m *Module) RunExample() {
	fmt.Println("Running GoDS example...")
	
	list, err := m.CreateList(gods.ArrayList)
	if err != nil {
		fmt.Printf("Error creating list: %v\n", err)
		return
	}
	
	arrayList, ok := list.(interface {
		Add(values ...interface{})
		Get(index int) (interface{}, bool)
		Size() int
		String() string
	})
	if !ok {
		fmt.Println("Error: could not type assert list")
		return
	}
	
	arrayList.Add("apple", "banana", "cherry")
	
	if val, found := arrayList.Get(1); found {
		fmt.Printf("Element at index 1: %v\n", val)
	}
	
	fmt.Printf("List size: %d\n", arrayList.Size())
	
	fmt.Printf("List contents: %s\n", arrayList.String())
	
	hashMap, err := m.CreateMap(gods.HashMap)
	if err != nil {
		fmt.Printf("Error creating map: %v\n", err)
		return
	}
	
	strMap, ok := hashMap.(interface {
		Put(key interface{}, value interface{})
		Get(key interface{}) (interface{}, bool)
		Size() int
		String() string
	})
	if !ok {
		fmt.Println("Error: could not type assert map")
		return
	}
	
	strMap.Put("one", "first")
	strMap.Put("two", "second")
	strMap.Put("three", "third")
	
	if val, found := strMap.Get("two"); found {
		fmt.Printf("Value for key 'two': %v\n", val)
	}
	
	fmt.Printf("Map size: %d\n", strMap.Size())
	
	fmt.Printf("Map contents: %s\n", strMap.String())
	
	fmt.Println("Example completed!")
}
