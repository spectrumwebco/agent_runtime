package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/actormodel"
	actormodelModule "github.com/spectrumwebco/agent_runtime/pkg/modules/actormodel"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

var actorModelCmd = &cobra.Command{
	Use:   "actormodel",
	Short: "Manage actor model systems",
	Long:  `Manage actor model systems for concurrent computation using the Ergo actor model.`,
}

var createSystemCmd = &cobra.Command{
	Use:   "create-system [name]",
	Short: "Create a new actor system",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		
		strategy, _ := cmd.Flags().GetString("strategy")
		maxRestarts, _ := cmd.Flags().GetInt("max-restarts")
		withinDuration, _ := cmd.Flags().GetInt("within-duration")
		
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		
		module := actormodelModule.NewModule(cfg)
		err = module.Initialize()
		if err != nil {
			return fmt.Errorf("failed to initialize actor model module: %w", err)
		}
		
		var supervisorStrategy actormodel.SupervisorStrategy
		switch strategy {
		case "one-for-one":
			supervisorStrategy = actormodel.OneForOne
		case "one-for-all":
			supervisorStrategy = actormodel.OneForAll
		case "rest-for-one":
			supervisorStrategy = actormodel.RestForOne
		default:
			return fmt.Errorf("invalid supervisor strategy: %s", strategy)
		}
		
		_, err = module.CreateActorSystem(name, actormodel.SupervisorOptions{
			ID:              "root",
			Strategy:        supervisorStrategy,
			MaxRestarts:     maxRestarts,
			WithinDuration:  time.Duration(withinDuration) * time.Second,
		})
		if err != nil {
			return fmt.Errorf("failed to create actor system: %w", err)
		}
		
		fmt.Printf("Actor system '%s' created successfully\n", name)
		return nil
	},
}

var startSystemCmd = &cobra.Command{
	Use:   "start-system [name]",
	Short: "Start an actor system",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		
		module := actormodelModule.NewModule(cfg)
		err = module.Initialize()
		if err != nil {
			return fmt.Errorf("failed to initialize actor model module: %w", err)
		}
		
		system, err := module.GetActorSystem(name)
		if err != nil {
			return fmt.Errorf("failed to get actor system: %w", err)
		}
		
		err = system.Start()
		if err != nil {
			return fmt.Errorf("failed to start actor system: %w", err)
		}
		
		fmt.Printf("Actor system '%s' started successfully\n", name)
		return nil
	},
}

var stopSystemCmd = &cobra.Command{
	Use:   "stop-system [name]",
	Short: "Stop an actor system",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		
		module := actormodelModule.NewModule(cfg)
		err = module.Initialize()
		if err != nil {
			return fmt.Errorf("failed to initialize actor model module: %w", err)
		}
		
		system, err := module.GetActorSystem(name)
		if err != nil {
			return fmt.Errorf("failed to get actor system: %w", err)
		}
		
		err = system.Stop()
		if err != nil {
			return fmt.Errorf("failed to stop actor system: %w", err)
		}
		
		fmt.Printf("Actor system '%s' stopped successfully\n", name)
		return nil
	},
}

var listSystemsCmd = &cobra.Command{
	Use:   "list-systems",
	Short: "List all actor systems",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		
		module := actormodelModule.NewModule(cfg)
		err = module.Initialize()
		if err != nil {
			return fmt.Errorf("failed to initialize actor model module: %w", err)
		}
		
		systems := module.ListActorSystems()
		
		if len(systems) == 0 {
			fmt.Println("No actor systems found")
			return nil
		}
		
		fmt.Println("Actor systems:")
		for _, system := range systems {
			fmt.Printf("- %s\n", system)
		}
		
		return nil
	},
}

var createActorCmd = &cobra.Command{
	Use:   "create-actor [system] [id]",
	Short: "Create a new actor in a system",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemName := args[0]
		actorID := args[1]
		
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		
		module := actormodelModule.NewModule(cfg)
		err = module.Initialize()
		if err != nil {
			return fmt.Errorf("failed to initialize actor model module: %w", err)
		}
		
		system, err := module.GetActorSystem(systemName)
		if err != nil {
			return fmt.Errorf("failed to get actor system: %w", err)
		}
		
		behavior := func(ctx context.Context, msg actormodel.Message) error {
			fmt.Printf("Actor %s received message: %s\n", actorID, msg.Type)
			
			if msg.ReplyTo != nil {
				reply := actormodel.Message{
					Type:    "echo",
					Payload: msg.Payload,
				}
				
				select {
				case msg.ReplyTo <- reply:
					fmt.Printf("Actor %s sent reply\n", actorID)
				case <-time.After(5 * time.Second):
					fmt.Printf("Actor %s timed out sending reply\n", actorID)
				case <-ctx.Done():
					fmt.Printf("Actor %s context done\n", actorID)
				}
			}
			
			return nil
		}
		
		actor, err := system.SpawnActor(actorID, behavior, nil)
		if err != nil {
			return fmt.Errorf("failed to create actor: %w", err)
		}
		
		err = actor.Start()
		if err != nil {
			return fmt.Errorf("failed to start actor: %w", err)
		}
		
		fmt.Printf("Actor '%s' created and started in system '%s'\n", actorID, systemName)
		return nil
	},
}

var sendMessageCmd = &cobra.Command{
	Use:   "send-message [system] [actor] [message-type]",
	Short: "Send a message to an actor",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		systemName := args[0]
		actorID := args[1]
		messageType := args[2]
		
		payload, _ := cmd.Flags().GetString("payload")
		waitReply, _ := cmd.Flags().GetBool("wait-reply")
		
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		
		module := actormodelModule.NewModule(cfg)
		err = module.Initialize()
		if err != nil {
			return fmt.Errorf("failed to initialize actor model module: %w", err)
		}
		
		system, err := module.GetActorSystem(systemName)
		if err != nil {
			return fmt.Errorf("failed to get actor system: %w", err)
		}
		
		actor, err := system.GetActor(actorID)
		if err != nil {
			return fmt.Errorf("failed to get actor: %w", err)
		}
		
		if waitReply {
			reply, err := actor.SendAndWait(messageType, payload)
			if err != nil {
				return fmt.Errorf("failed to send message and wait for reply: %w", err)
			}
			
			fmt.Printf("Sent message '%s' to actor '%s' and received reply: %v\n", messageType, actorID, reply.Payload)
		} else {
			msg := actormodel.Message{
				Type:    messageType,
				Payload: payload,
			}
			
			err = actor.Send(msg)
			if err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}
			
			fmt.Printf("Sent message '%s' to actor '%s'\n", messageType, actorID)
		}
		
		return nil
	},
}

var runExampleCmd = &cobra.Command{
	Use:   "run-example",
	Short: "Run an example actor system",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		
		module := actormodelModule.NewModule(cfg)
		err = module.Initialize()
		if err != nil {
			return fmt.Errorf("failed to initialize actor model module: %w", err)
		}
		
		system, err := module.CreateActorSystem("example", actormodel.SupervisorOptions{
			ID:              "root",
			Strategy:        actormodel.OneForOne,
			MaxRestarts:     10,
			WithinDuration:  60 * time.Second,
		})
		if err != nil {
			return fmt.Errorf("failed to create actor system: %w", err)
		}
		
		err = system.Start()
		if err != nil {
			return fmt.Errorf("failed to start actor system: %w", err)
		}
		
		fmt.Println("Actor system 'example' created and started")
		
		workerBehavior := func(ctx context.Context, msg actormodel.Message) error {
			fmt.Printf("Worker received message: %s\n", msg.Type)
			
			if msg.Type == "work" {
				time.Sleep(1 * time.Second)
				
				if msg.ReplyTo != nil {
					reply := actormodel.Message{
						Type:    "work_done",
						Payload: "Work completed successfully",
					}
					
					select {
					case msg.ReplyTo <- reply:
						fmt.Println("Worker sent reply")
					case <-time.After(5 * time.Second):
						fmt.Println("Worker timed out sending reply")
					case <-ctx.Done():
						fmt.Println("Worker context done")
					}
				}
			}
			
			return nil
		}
		
		supervisorBehavior := func(ctx context.Context, msg actormodel.Message) error {
			fmt.Printf("Supervisor received message: %s\n", msg.Type)
			
			if msg.Type == "create_workers" {
				count, ok := msg.Payload.(int)
				if !ok {
					count = 3
				}
				
				for i := 0; i < count; i++ {
					workerID := fmt.Sprintf("worker-%d", i+1)
					
					actor, ok := msg.Sender.(*actormodel.Actor)
					if !ok {
						return fmt.Errorf("sender is not an actor")
					}
					
					worker, err := actor.Spawn(workerID, workerBehavior, nil)
					if err != nil {
						return fmt.Errorf("failed to create worker: %w", err)
					}
					
					err = worker.Start()
					if err != nil {
						return fmt.Errorf("failed to start worker: %w", err)
					}
					
					fmt.Printf("Created and started worker '%s'\n", workerID)
				}
				
				if msg.ReplyTo != nil {
					reply := actormodel.Message{
						Type:    "workers_created",
						Payload: count,
					}
					
					select {
					case msg.ReplyTo <- reply:
						fmt.Println("Supervisor sent reply")
					case <-time.After(5 * time.Second):
						fmt.Println("Supervisor timed out sending reply")
					case <-ctx.Done():
						fmt.Println("Supervisor context done")
					}
				}
			}
			
			return nil
		}
		
		supervisor, err := system.SpawnSupervisor(actormodel.SupervisorOptions{
			ID:              "supervisor",
			Strategy:        actormodel.OneForOne,
			MaxRestarts:     5,
			WithinDuration:  30 * time.Second,
		})
		if err != nil {
			return fmt.Errorf("failed to create supervisor: %w", err)
		}
		
		err = supervisor.Start()
		if err != nil {
			return fmt.Errorf("failed to start supervisor: %w", err)
		}
		
		fmt.Println("Supervisor created and started")
		
		msg := actormodel.Message{
			Type:    "create_workers",
			Payload: 3,
			ReplyTo: make(chan actormodel.Message, 1),
		}
		
		err = supervisor.Send(msg)
		if err != nil {
			return fmt.Errorf("failed to send message to supervisor: %w", err)
		}
		
		select {
		case reply := <-msg.ReplyTo:
			fmt.Printf("Received reply from supervisor: %s - %v\n", reply.Type, reply.Payload)
		case <-time.After(10 * time.Second):
			fmt.Println("Timed out waiting for reply from supervisor")
		}
		
		worker, err := system.GetActor("worker-1")
		if err != nil {
			return fmt.Errorf("failed to get worker: %w", err)
		}
		
		workMsg := actormodel.Message{
			Type:    "work",
			Payload: "Do some work",
			ReplyTo: make(chan actormodel.Message, 1),
		}
		
		err = worker.Send(workMsg)
		if err != nil {
			return fmt.Errorf("failed to send message to worker: %w", err)
		}
		
		select {
		case reply := <-workMsg.ReplyTo:
			fmt.Printf("Received reply from worker: %s - %v\n", reply.Type, reply.Payload)
		case <-time.After(10 * time.Second):
			fmt.Println("Timed out waiting for reply from worker")
		}
		
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		
		fmt.Println("Press Ctrl+C to stop the example")
		
		<-sigCh
		
		err = system.Stop()
		if err != nil {
			return fmt.Errorf("failed to stop actor system: %w", err)
		}
		
		fmt.Println("Actor system stopped")
		
		return nil
	},
}

func NewActorModelCommand() *cobra.Command {
	createSystemCmd.Flags().String("strategy", "one-for-one", "Supervisor strategy (one-for-one, one-for-all, rest-for-one)")
	createSystemCmd.Flags().Int("max-restarts", 10, "Maximum number of restarts")
	createSystemCmd.Flags().Int("within-duration", 60, "Time window for restarts in seconds")
	
	sendMessageCmd.Flags().String("payload", "", "Message payload")
	sendMessageCmd.Flags().Bool("wait-reply", false, "Wait for a reply")
	
	actorModelCmd.AddCommand(createSystemCmd)
	actorModelCmd.AddCommand(startSystemCmd)
	actorModelCmd.AddCommand(stopSystemCmd)
	actorModelCmd.AddCommand(listSystemsCmd)
	actorModelCmd.AddCommand(createActorCmd)
	actorModelCmd.AddCommand(sendMessageCmd)
	actorModelCmd.AddCommand(runExampleCmd)
	
	return actorModelCmd
}
