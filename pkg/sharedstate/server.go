package sharedstate

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spectrumwebco/agent_runtime/internal/eventstream"
	"github.com/spectrumwebco/agent_runtime/internal/statemanager"
	"github.com/spectrumwebco/agent_runtime/internal/statemanager/supabase"
	"github.com/spectrumwebco/agent_runtime/pkg/ai/vercel"
)

type Server struct {
	router             *gin.Engine
	sharedStateManager *SharedStateManager
	integration        *Integration
	aiClient           *vercel.VercelAIClient
	actionAdapter      *vercel.ServerActionsAdapter
	port               int
}

type ServerConfig struct {
	Port               int
	EventStream        *eventstream.Stream
	StateManager       *statemanager.StateManager
	SupabaseClient     *supabase.Client
	SharedStateManager *SharedStateManager
	AIClient           *vercel.VercelAIClient
	ActionAdapter      *vercel.ServerActionsAdapter
}

func NewServer(cfg ServerConfig) (*Server, error) {
	if cfg.EventStream == nil {
		return nil, fmt.Errorf("event stream is required")
	}
	if cfg.StateManager == nil {
		return nil, fmt.Errorf("state manager is required")
	}
	if cfg.SharedStateManager == nil {
		sharedStateCfg := SharedStateConfig{
			EventStream:  cfg.EventStream,
			StateManager: cfg.StateManager,
		}
		var err error
		cfg.SharedStateManager, err = NewSharedStateManager(sharedStateCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create shared state manager: %w", err)
		}
	}

	if cfg.AIClient == nil {
		cfg.AIClient = vercel.NewVercelAIClient(cfg.SharedStateManager, cfg.EventStream)
	}

	if cfg.ActionAdapter == nil {
		cfg.ActionAdapter = vercel.NewServerActionsAdapter(cfg.SharedStateManager, cfg.AIClient)
	}

	integrationCfg := IntegrationConfig{
		EventStream:        cfg.EventStream,
		StateManager:       cfg.StateManager,
		SupabaseClient:     cfg.SupabaseClient,
		SharedStateManager: cfg.SharedStateManager,
	}
	integration, err := NewIntegration(integrationCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create integration: %w", err)
	}

	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	server := &Server{
		router:             router,
		sharedStateManager: cfg.SharedStateManager,
		integration:        integration,
		aiClient:           cfg.AIClient,
		actionAdapter:      cfg.ActionAdapter,
		port:               cfg.Port,
	}

	server.registerRoutes()

	return server, nil
}

func (s *Server) registerRoutes() {
	s.sharedStateManager.RegisterWebSocketHandler(s.router)
	
	s.router.Use(vercel.StreamingMiddleware(s.aiClient))
	s.router.Use(vercel.ActionMiddleware(s.actionAdapter))

	api := s.router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		})

		api.GET("/state/:type/:id", s.getState)

		api.POST("/state/:type/:id", s.updateState)
		
		api.GET("/actions", s.listActions)
	}
}

func (s *Server) getState(c *gin.Context) {
	stateType := c.Param("type")
	stateID := c.Param("id")

	state, err := s.sharedStateManager.GetState(StateType(stateType), stateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to get state: %v", err),
		})
		return
	}

	if state == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "state not found",
		})
		return
	}

	c.JSON(http.StatusOK, state)
}

func (s *Server) updateState(c *gin.Context) {
	stateType := c.Param("type")
	stateID := c.Param("id")

	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("failed to parse request body: %v", err),
		})
		return
	}

	if err := s.sharedStateManager.UpdateState(StateType(stateType), stateID, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to update state: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Starting shared state server on %s", addr)
	return s.router.Run(addr)
}

func (s *Server) listActions(c *gin.Context) {
	actions := s.actionAdapter.ListActions()
	
	c.JSON(http.StatusOK, gin.H{
		"actions": actions,
	})
}

func (s *Server) Close() error {
	if s.integration != nil {
		if err := s.integration.Close(); err != nil {
			log.Printf("Error closing integration: %v", err)
		}
	}

	if s.sharedStateManager != nil {
		if err := s.sharedStateManager.Close(); err != nil {
			log.Printf("Error closing shared state manager: %v", err)
		}
	}
	
	if s.aiClient != nil {
		s.aiClient.Close()
	}

	return nil
}
