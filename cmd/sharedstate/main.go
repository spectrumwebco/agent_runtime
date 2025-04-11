package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spectrumwebco/agent_runtime/internal/config"
	"github.com/spectrumwebco/agent_runtime/internal/eventstream"
	"github.com/spectrumwebco/agent_runtime/internal/statemanager"
	"github.com/spectrumwebco/agent_runtime/internal/statemanager/supabase"
	"github.com/spectrumwebco/agent_runtime/pkg/sharedstate"
)

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	var cfg *config.Config
	var err error
	if *configPath != "" {
		cfg, err = config.LoadConfig(*configPath)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
	} else {
		cfg = config.DefaultConfig()
	}

	eventStream, err := eventstream.NewStream(cfg)
	if err != nil {
		log.Fatalf("Failed to create event stream: %v", err)
	}

	stateManager, err := statemanager.NewStateManager(cfg)
	if err != nil {
		log.Fatalf("Failed to create state manager: %v", err)
	}

	var supabaseStateManager *supabase.SupabaseStateManager
	if cfg.Supabase.Enabled {
		supabaseStateManager = supabase.NewSupabaseStateManager(supabase.SupabaseStateConfig{
			MainURL:     cfg.Supabase.MainURL,
			ReadonlyURL: cfg.Supabase.ReadonlyURL,
			RollbackURL: cfg.Supabase.RollbackURL,
			APIKey:      cfg.Supabase.APIKey,
			AuthToken:   cfg.Supabase.AuthToken,
		})
	}

	var supabaseClient *supabase.Client
	if cfg.Supabase.Enabled {
		supabaseClient = supabase.NewSupabaseClient(supabase.SupabaseConfig{
			URL:       cfg.Supabase.MainURL,
			APIKey:    cfg.Supabase.APIKey,
			AuthToken: cfg.Supabase.AuthToken,
		})
	}

	sharedStateManager, err := sharedstate.NewSharedStateManager(sharedstate.SharedStateConfig{
		EventStream:  eventStream,
		StateManager: stateManager,
	})
	if err != nil {
		log.Fatalf("Failed to create shared state manager: %v", err)
	}

	server, err := sharedstate.NewServer(sharedstate.ServerConfig{
		Port:               *port,
		EventStream:        eventStream,
		StateManager:       stateManager,
		SupabaseClient:     supabaseClient,
		SharedStateManager: sharedStateManager,
	})
	if err != nil {
		log.Fatalf("Failed to create shared state server: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down...")
		server.Close()
		stateManager.Close()
		os.Exit(0)
	}()

	log.Printf("Starting shared state server on port %d", *port)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
