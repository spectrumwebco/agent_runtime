package kled

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type DjangoModel interface {
	TableName() string
	Fields() map[string]string
}

type DjangoView interface {
	Handle(ctx context.Context, params map[string]interface{}) (interface{}, error)
}

type DjangoAdmin interface {
	ModelName() string
	Fields() []string
}

type Django struct {
	Models map[string]DjangoModel
	Views  map[string]DjangoView
	Admins map[string]DjangoAdmin
}

type DjangoConfig struct {
	AppDir       string `json:"app_dir"`
	TemplatesDir string `json:"templates_dir"`
	StaticDir    string `json:"static_dir"`
	Port         int    `json:"port"`
	Debug        bool   `json:"debug"`
}

type DjangoIntegration struct {
	Framework *Framework

	Router *gin.Engine

	Server *http.Server

	Django *Django

	mutex sync.RWMutex

	Context map[string]interface{}
}

var _ Integration = (*DjangoIntegration)(nil)

type DjangoIntegrationConfig struct {
	Address string `json:"address"`

	Routes []DjangoRoute `json:"routes"`

	DjangoConfig DjangoConfig `json:"django_config"`
}

type DjangoRoute struct {
	Path string `json:"path"`

	Method string `json:"method"`

	Handler string `json:"handler"`
}

func NewDjangoIntegration(framework *Framework, config DjangoIntegrationConfig) (*DjangoIntegration, error) {
	if framework == nil {
		return nil, fmt.Errorf("framework cannot be nil")
	}

	router := gin.Default()

	django := &Django{
		Models: make(map[string]DjangoModel),
		Views:  make(map[string]DjangoView),
		Admins: make(map[string]DjangoAdmin),
	}

	integration := &DjangoIntegration{
		Framework: framework,
		Router:    router,
		Django:    django,
		Context:   make(map[string]interface{}),
	}

	for _, route := range config.Routes {
		handler := integration.createHandler(route.Handler)

		switch route.Method {
		case "GET":
			router.GET(route.Path, handler)
		case "POST":
			router.POST(route.Path, handler)
		case "PUT":
			router.PUT(route.Path, handler)
		case "DELETE":
			router.DELETE(route.Path, handler)
		default:
			return nil, fmt.Errorf("unsupported method: %s", route.Method)
		}
	}

	server := &http.Server{
		Addr:    config.Address,
		Handler: router,
	}

	integration.Server = server

	if err := framework.RegisterIntegration(integration); err != nil {
		return nil, fmt.Errorf("failed to register Django integration: %w", err)
	}

	if framework.Config.Debug {
		log.Printf("Django integration created with %d routes\n", len(config.Routes))
	}

	return integration, nil
}

func (d *DjangoIntegration) Name() string {
	return "django"
}

func (d *DjangoIntegration) Start() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	go func() {
		if err := d.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting Django integration HTTP server: %v\n", err)
		}
	}()

	if d.Framework.Config.Debug {
		log.Printf("Django integration started on %s\n", d.Server.Addr)
	}

	return nil
}

func (d *DjangoIntegration) Stop() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.Server.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("error shutting down Django integration HTTP server: %w", err)
	}

	if d.Framework.Config.Debug {
		log.Println("Django integration stopped")
	}

	return nil
}

func (d *DjangoIntegration) createHandler(handlerName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData map[string]interface{}
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := context.WithValue(c.Request.Context(), "django_integration", d)

		result, err := d.Framework.ExecuteTool(ctx, handlerName, requestData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func (d *DjangoIntegration) GetContext() map[string]interface{} {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	contextCopy := make(map[string]interface{})
	for k, v := range d.Context {
		contextCopy[k] = v
	}

	return contextCopy
}

func (d *DjangoIntegration) SetContext(ctx map[string]interface{}) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.Context = ctx
}

func (d *DjangoIntegration) UpdateContext(updates map[string]interface{}) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for k, v := range updates {
		d.Context[k] = v
	}
}

func (d *DjangoIntegration) SaveState() ([]byte, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	state := map[string]interface{}{
		"context": d.Context,
	}

	return json.Marshal(state)
}

func (d *DjangoIntegration) LoadState(data []byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	var state map[string]interface{}
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("failed to unmarshal state: %w", err)
	}

	if context, ok := state["context"].(map[string]interface{}); ok {
		d.Context = context
	}

	return nil
}

func (d *DjangoIntegration) RegisterModel(name string, fields map[string]string) error {
	model := &djangoModel{
		name:   name,
		fields: fields,
	}
	d.Django.Models[name] = model
	return nil
}

func (d *DjangoIntegration) RegisterRoute(path string, method string, handler func(ctx context.Context, params map[string]interface{}) (interface{}, error)) error {
	d.Django.Views[path] = &djangoView{
		handler: handler,
	}
	return nil
}

type djangoModel struct {
	name   string
	fields map[string]string
}

func (m *djangoModel) TableName() string {
	return m.name
}

func (m *djangoModel) Fields() map[string]string {
	return m.fields
}

type djangoView struct {
	handler func(ctx context.Context, params map[string]interface{}) (interface{}, error)
}

func (v *djangoView) Handle(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	return v.handler(ctx, params)
}

func (d *DjangoIntegration) Models() map[string]DjangoModel {
	return d.Django.Models
}

func (d *DjangoIntegration) Routes() map[string]DjangoView {
	return d.Django.Views
}
