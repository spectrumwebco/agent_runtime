package modules

type Module interface {
	Name() string
	Description() string
	Initialize(/* TODO: Add context/config */) error
	Execute(/* TODO: Add input/context */) ( /* TODO: Add output */ interface{}, error)
	Cleanup() error
}
