package veigar

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type SecurityModel interface {
	TableName() string
	Fields() map[string]string
}

type SecurityView interface {
	Handle(ctx context.Context, params map[string]interface{}) (interface{}, error)
}

type SecurityAdmin interface {
	ModelName() string
	Fields() []string
}

type SecurityDjango struct {
	Models map[string]SecurityModel
	Views  map[string]SecurityView
	Admins map[string]SecurityAdmin
}

type SecurityDjangoConfig struct {
	AppDir       string `json:"app_dir"`
	TemplatesDir string `json:"templates_dir"`
	StaticDir    string `json:"static_dir"`
	Port         int    `json:"port"`
	Debug        bool   `json:"debug"`
}

type SecurityDjangoIntegration struct {
	Framework *Framework

	Router *gin.Engine

	Server *http.Server

	Django *SecurityDjango

	mutex sync.RWMutex

	Context map[string]interface{}
}

var _ Integration = (*SecurityDjangoIntegration)(nil)

type SecurityDjangoIntegrationConfig struct {
	Address string `json:"address"`

	Routes []SecurityDjangoRoute `json:"routes"`

	DjangoConfig SecurityDjangoConfig `json:"django_config"`
}

type SecurityDjangoRoute struct {
	Path string `json:"path"`

	Method string `json:"method"`

	Handler string `json:"handler"`
}

func NewSecurityDjangoIntegration(framework *Framework, config SecurityDjangoIntegrationConfig) (*SecurityDjangoIntegration, error) {
	if framework == nil {
		return nil, fmt.Errorf("framework cannot be nil")
	}

	router := gin.Default()

	django := &SecurityDjango{
		Models: make(map[string]SecurityModel),
		Views:  make(map[string]SecurityView),
		Admins: make(map[string]SecurityAdmin),
	}

	integration := &SecurityDjangoIntegration{
		Framework: framework,
		Router:    router,
		Django:    django,
		Context:   make(map[string]interface{}),
	}

	router.Use(integration.securityMiddleware())

	for _, route := range config.Routes {
		handler := integration.createSecurityHandler(route.Handler)

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

	integration.registerSecurityAPIRoutes(router)

	server := &http.Server{
		Addr:    config.Address,
		Handler: router,
	}

	integration.Server = server

	if err := framework.RegisterIntegration(integration); err != nil {
		return nil, fmt.Errorf("failed to register Security Django integration: %w", err)
	}

	if framework.Config.Debug {
		log.Printf("Security Django integration created with %d routes\n", len(config.Routes))
	}

	return integration, nil
}

func (d *SecurityDjangoIntegration) Name() string {
	return "security_django"
}

func (d *SecurityDjangoIntegration) Start() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	go func() {
		if err := d.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting Security Django integration HTTP server: %v\n", err)
		}
	}()

	if d.Framework.Config.Debug {
		log.Printf("Security Django integration started on %s\n", d.Server.Addr)
	}

	return nil
}

func (d *SecurityDjangoIntegration) Stop() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.Server.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("error shutting down Security Django integration HTTP server: %w", err)
	}

	if d.Framework.Config.Debug {
		log.Println("Security Django integration stopped")
	}

	return nil
}

func (d *SecurityDjangoIntegration) securityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Content-Security-Policy", "default-src 'self'")
		
		c.Next()
	}
}

func (d *SecurityDjangoIntegration) registerSecurityAPIRoutes(router *gin.Engine) {
	router.POST("/api/security/review", func(c *gin.Context) {
		var requestData map[string]interface{}
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := context.WithValue(c.Request.Context(), "security_django_integration", d)
		result, err := d.Framework.ExecuteSecurityTool(ctx, "security_review", requestData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	})

	router.POST("/api/security/scan", func(c *gin.Context) {
		var requestData map[string]interface{}
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := context.WithValue(c.Request.Context(), "security_django_integration", d)
		result, err := d.Framework.ExecuteSecurityTool(ctx, "vulnerability_scan", requestData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	})

	router.POST("/api/security/compliance", func(c *gin.Context) {
		var requestData map[string]interface{}
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := context.WithValue(c.Request.Context(), "security_django_integration", d)
		result, err := d.Framework.ExecuteSecurityTool(ctx, "compliance_check", requestData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	})

	router.GET("/api/security/status/:id", func(c *gin.Context) {
		id := c.Param("id")
		
		ctx := context.WithValue(c.Request.Context(), "security_django_integration", d)
		result, err := d.Framework.ExecuteSecurityTool(ctx, "security_status", map[string]interface{}{
			"id": id,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	})
}

func (d *SecurityDjangoIntegration) createSecurityHandler(handlerName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData map[string]interface{}
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := context.WithValue(c.Request.Context(), "security_django_integration", d)

		result, err := d.Framework.ExecuteSecurityTool(ctx, handlerName, requestData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func (d *SecurityDjangoIntegration) GetContext() map[string]interface{} {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	contextCopy := make(map[string]interface{})
	for k, v := range d.Context {
		contextCopy[k] = v
	}

	return contextCopy
}

func (d *SecurityDjangoIntegration) SetContext(ctx map[string]interface{}) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.Context = ctx
}

func (d *SecurityDjangoIntegration) UpdateContext(updates map[string]interface{}) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for k, v := range updates {
		d.Context[k] = v
	}
}

func (d *SecurityDjangoIntegration) SaveState() ([]byte, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	state := map[string]interface{}{
		"context": d.Context,
	}

	return json.Marshal(state)
}

func (d *SecurityDjangoIntegration) LoadState(data []byte) error {
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

func (d *SecurityDjangoIntegration) RegisterSecurityModel(name string, fields map[string]string) error {
	model := &securityModel{
		name:   name,
		fields: fields,
	}
	d.Django.Models[name] = model
	return nil
}

func (d *SecurityDjangoIntegration) RegisterSecurityRoute(path string, method string, handler func(ctx context.Context, params map[string]interface{}) (interface{}, error)) error {
	d.Django.Views[path] = &securityView{
		handler: handler,
	}
	return nil
}

type securityModel struct {
	name   string
	fields map[string]string
}

func (m *securityModel) TableName() string {
	return m.name
}

func (m *securityModel) Fields() map[string]string {
	return m.fields
}

type securityView struct {
	handler func(ctx context.Context, params map[string]interface{}) (interface{}, error)
}

func (v *securityView) Handle(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	return v.handler(ctx, params)
}

func (d *SecurityDjangoIntegration) SecurityModels() map[string]SecurityModel {
	return d.Django.Models
}

func (d *SecurityDjangoIntegration) SecurityRoutes() map[string]SecurityView {
	return d.Django.Views
}
