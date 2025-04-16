package framework

import (
	"context"
	"fmt"

	"github.com/spectrumwebco/agent_runtime/internal/actormodel"
	"github.com/spectrumwebco/agent_runtime/internal/auth/casbin"
	"github.com/spectrumwebco/agent_runtime/internal/database/gorm"
	"github.com/spectrumwebco/agent_runtime/internal/langgraph"
	"github.com/spectrumwebco/agent_runtime/internal/microservices/gomicro"
	"github.com/spectrumwebco/agent_runtime/internal/webframework/gin"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Framework struct {
	config          *config.Config
	router          *gin.Router
	enforcer        *casbin.Enforcer
	database        *gorm.Database
	serviceRegistry *gomicro.ServiceRegistry
	actorSystem     *actormodel.ActorSystem
	graphManager    *langgraph.GraphManager
}

type FrameworkOptions struct {
	ConfigPath string
	ModelPath  string
	PolicyPath string
	ModelText  string
}

func NewFramework(opts FrameworkOptions) (*Framework, error) {
	cfg, err := config.LoadConfig(opts.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	router := gin.NewRouter()

	enforcer, err := casbin.NewEnforcer(casbin.Config{
		ModelPath:  opts.ModelPath,
		PolicyPath: opts.PolicyPath,
		ModelText:  opts.ModelText,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create enforcer: %w", err)
	}

	dbConfig := &gorm.Config{
		Type:         cfg.GetString("database.type"),
		Host:         cfg.GetString("database.host"),
		Port:         cfg.GetInt("database.port"),
		Username:     cfg.GetString("database.username"),
		Password:     cfg.GetString("database.password"),
		Database:     cfg.GetString("database.name"),
		SSLMode:      cfg.GetString("database.ssl_mode"),
		MaxOpenConns: cfg.GetInt("database.max_open_conns"),
		MaxIdleConns: cfg.GetInt("database.max_idle_conns"),
		TablePrefix:  cfg.GetString("database.table_prefix"),
		Debug:        cfg.GetBool("database.debug"),
	}

	if dbConfig.Type == "" {
		dbConfig.Type = "sqlite"
	}

	if dbConfig.Database == "" && dbConfig.Type == "sqlite" {
		dbConfig.Database = "kled.db"
	}

	database, err := gorm.NewDatabase(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	serviceRegistry := gomicro.NewServiceRegistry()

	actorSystem := actormodel.NewActorSystem(actormodel.ActorSystemOptions{
		Name: "kled",
		SupervisorOptions: actormodel.SupervisorOptions{
			ID:              "root",
			Strategy:        actormodel.OneForOne,
			MaxRestarts:     10,
			WithinDuration:  60,
		},
	})

	graphManager := langgraph.NewGraphManager()

	return &Framework{
		config:          cfg,
		router:          router,
		enforcer:        enforcer,
		database:        database,
		serviceRegistry: serviceRegistry,
		actorSystem:     actorSystem,
		graphManager:    graphManager,
	}, nil
}

func (f *Framework) Initialize() error {
	err := f.database.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	err = f.enforcer.LoadPolicy()
	if err != nil {
		return fmt.Errorf("failed to load policy: %w", err)
	}

	err = f.actorSystem.Start()
	if err != nil {
		return fmt.Errorf("failed to start actor system: %w", err)
	}

	return nil
}

func (f *Framework) Start(ctx context.Context) error {
	go func() {
		err := f.router.Run(f.config.GetString("server.address"))
		if err != nil {
			fmt.Printf("HTTP server error: %s\n", err)
		}
	}()

	return nil
}

func (f *Framework) Stop(ctx context.Context) error {
	err := f.actorSystem.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop actor system: %w", err)
	}

	err = f.database.Close()
	if err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	return nil
}

func (f *Framework) Router() *gin.Router {
	return f.router
}

func (f *Framework) Enforcer() *casbin.Enforcer {
	return f.enforcer
}

func (f *Framework) Database() *gorm.Database {
	return f.database
}

func (f *Framework) ServiceRegistry() *gomicro.ServiceRegistry {
	return f.serviceRegistry
}

func (f *Framework) ActorSystem() *actormodel.ActorSystem {
	return f.actorSystem
}

func (f *Framework) GraphManager() *langgraph.GraphManager {
	return f.graphManager
}

func (f *Framework) Config() *config.Config {
	return f.config
}

func (f *Framework) CreateService(name, version, address string, metadata map[string]string) (*gomicro.Service, error) {
	service := gomicro.NewService(gomicro.ServiceOptions{
		Name:      name,
		Version:   version,
		Address:   address,
		Metadata:  metadata,
		Timeout:   f.config.GetDuration("microservices.timeout"),
		Retries:   f.config.GetInt("microservices.retries"),
		Namespace: f.config.GetString("microservices.namespace"),
	})

	f.serviceRegistry.Register(service)

	return service, nil
}

func (f *Framework) CreateGraph(name string) (*langgraph.Graph, error) {
	graph := langgraph.NewGraph(name)
	f.graphManager.RegisterGraph(graph)
	return graph, nil
}

func (f *Framework) CreateActor(id string, behavior actormodel.Behavior, state map[string]interface{}) (*actormodel.Actor, error) {
	return f.actorSystem.SpawnActor(id, behavior, state)
}

func (f *Framework) CreateSupervisor(opts actormodel.SupervisorOptions) (*actormodel.Supervisor, error) {
	return f.actorSystem.SpawnSupervisor(opts)
}

func (f *Framework) RunMigrations() error {
	return f.database.Migrate(gorm.GetModels()...)
}

func (f *Framework) WithAuth(group *gin.RouterGroup) *gin.RouterGroup {
	return group.WithAuth()
}

func (f *Framework) WithDomainAuth(group *gin.RouterGroup) *gin.RouterGroup {
	return group.WithDomainAuth()
}

func (f *Framework) Enforce(params ...interface{}) (bool, error) {
	return f.enforcer.Enforce(params...)
}
