package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/langgraph"
	langgraphModule "github.com/spectrumwebco/agent_runtime/pkg/modules/langgraph"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

var langgraphCmd = &cobra.Command{
	Use:   "langgraph",
	Short: "LangGraph-Go commands for agent workflows",
	Long:  `Commands for working with LangGraph-Go agent workflows.`,
}

var runWorkflowCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a LangGraph-Go agent workflow",
	Long:  `Run a LangGraph-Go agent workflow with the specified input and tools.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile, _ := cmd.Flags().GetString("input")
		toolsFile, _ := cmd.Flags().GetString("tools")
		djangoBaseURL, _ := cmd.Flags().GetString("django-url")
		
		inputData, err := ioutil.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}
		
		toolsData, err := ioutil.ReadFile(toolsFile)
		if err != nil {
			return fmt.Errorf("failed to read tools file: %w", err)
		}
		
		var tools []map[string]interface{}
		err = json.Unmarshal(toolsData, &tools)
		if err != nil {
			return fmt.Errorf("failed to parse tools file: %w", err)
		}
		
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		
		module := langgraphModule.NewModule(cfg, djangoBaseURL)
		err = module.Initialize()
		if err != nil {
			return fmt.Errorf("failed to initialize LangGraph module: %w", err)
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		
		result, err := module.ExecuteAgentWorkflow(ctx, string(inputData), tools)
		if err != nil {
			return fmt.Errorf("failed to execute agent workflow: %w", err)
		}
		
		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal result: %w", err)
		}
		
		fmt.Println(string(resultJSON))
		
		return nil
	},
}

var createGraphCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a LangGraph-Go agent workflow graph",
	Long:  `Create a LangGraph-Go agent workflow graph from a configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configFile, _ := cmd.Flags().GetString("config")
		outputFile, _ := cmd.Flags().GetString("output")
		
		configData, err := ioutil.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}
		
		var graphConfig struct {
			EntryPoint string                 `json:"entry_point"`
			ExitPoints []string               `json:"exit_points"`
			Nodes      map[string]interface{} `json:"nodes"`
			Edges      []struct {
				Source    string `json:"source"`
				Target    string `json:"target"`
				Condition string `json:"condition"`
			} `json:"edges"`
			StateSchema map[string]string `json:"state_schema"`
		}
		
		err = json.Unmarshal(configData, &graphConfig)
		if err != nil {
			return fmt.Errorf("failed to parse config file: %w", err)
		}
		
		graph := langgraph.NewGraph(langgraph.NodeID(graphConfig.EntryPoint))
		
		graph.DefineStateSchema(graphConfig.StateSchema)
		
		for id, _ := range graphConfig.Nodes {
			err := graph.AddNode(langgraph.NodeID(id), func(ctx context.Context, state map[string]interface{}) (map[string]interface{}, error) {
				return state, nil
			}, map[string]interface{}{
				"description": fmt.Sprintf("Node %s", id),
			})
			if err != nil {
				return fmt.Errorf("failed to add node %s: %w", id, err)
			}
		}
		
		for _, edge := range graphConfig.Edges {
			err := graph.AddEdge(langgraph.NodeID(edge.Source), langgraph.NodeID(edge.Target), langgraph.DefaultEdge(langgraph.NodeID(edge.Target)))
			if err != nil {
				return fmt.Errorf("failed to add edge from %s to %s: %w", edge.Source, edge.Target, err)
			}
		}
		
		for _, exitPoint := range graphConfig.ExitPoints {
			err := graph.SetExitPoint(langgraph.NodeID(exitPoint))
			if err != nil {
				return fmt.Errorf("failed to set exit point %s: %w", exitPoint, err)
			}
		}
		
		err = graph.Compile()
		if err != nil {
			return fmt.Errorf("failed to compile graph: %w", err)
		}
		
		if outputFile != "" {
			graphJSON, err := json.MarshalIndent(graphConfig, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal graph config: %w", err)
			}
			
			err = ioutil.WriteFile(outputFile, graphJSON, 0644)
			if err != nil {
				return fmt.Errorf("failed to write graph config: %w", err)
			}
			
			fmt.Printf("Graph configuration saved to %s\n", outputFile)
		}
		
		fmt.Println("Graph compiled successfully")
		
		return nil
	},
}

var visualizeGraphCmd = &cobra.Command{
	Use:   "visualize",
	Short: "Visualize a LangGraph-Go agent workflow graph",
	Long:  `Visualize a LangGraph-Go agent workflow graph from a configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configFile, _ := cmd.Flags().GetString("config")
		outputFile, _ := cmd.Flags().GetString("output")
		
		configData, err := ioutil.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}
		
		var graphConfig struct {
			EntryPoint string                 `json:"entry_point"`
			ExitPoints []string               `json:"exit_points"`
			Nodes      map[string]interface{} `json:"nodes"`
			Edges      []struct {
				Source    string `json:"source"`
				Target    string `json:"target"`
				Condition string `json:"condition"`
			} `json:"edges"`
		}
		
		err = json.Unmarshal(configData, &graphConfig)
		if err != nil {
			return fmt.Errorf("failed to parse config file: %w", err)
		}
		
		dot := "digraph G {\n"
		dot += "  rankdir=LR;\n"
		dot += "  node [shape=box, style=filled, fillcolor=lightblue];\n"
		
		for id := range graphConfig.Nodes {
			if id == graphConfig.EntryPoint {
				dot += fmt.Sprintf("  \"%s\" [fillcolor=lightgreen];\n", id)
			} else {
				isExitPoint := false
				for _, exitPoint := range graphConfig.ExitPoints {
					if id == exitPoint {
						isExitPoint = true
						break
					}
				}
				
				if isExitPoint {
					dot += fmt.Sprintf("  \"%s\" [fillcolor=lightcoral];\n", id)
				} else {
					dot += fmt.Sprintf("  \"%s\";\n", id)
				}
			}
		}
		
		for _, edge := range graphConfig.Edges {
			dot += fmt.Sprintf("  \"%s\" -> \"%s\"", edge.Source, edge.Target)
			if edge.Condition != "" && edge.Condition != "default" {
				dot += fmt.Sprintf(" [label=\"%s\"]", edge.Condition)
			}
			dot += ";\n"
		}
		
		dot += "}\n"
		
		if outputFile != "" {
			err = ioutil.WriteFile(outputFile, []byte(dot), 0644)
			if err != nil {
				return fmt.Errorf("failed to write DOT file: %w", err)
			}
			
			fmt.Printf("Graph visualization saved to %s\n", outputFile)
			fmt.Println("You can visualize this file using Graphviz: dot -Tpng -o graph.png", outputFile)
		} else {
			fmt.Println(dot)
		}
		
		return nil
	},
}

func init() {
	RootCmd.AddCommand(langgraphCmd)
	
	langgraphCmd.AddCommand(runWorkflowCmd)
	langgraphCmd.AddCommand(createGraphCmd)
	langgraphCmd.AddCommand(visualizeGraphCmd)
	
	runWorkflowCmd.Flags().String("input", "", "Path to input file")
	runWorkflowCmd.Flags().String("tools", "", "Path to tools file")
	runWorkflowCmd.Flags().String("django-url", "http://localhost:8000", "Django base URL")
	runWorkflowCmd.MarkFlagRequired("input")
	runWorkflowCmd.MarkFlagRequired("tools")
	
	createGraphCmd.Flags().String("config", "", "Path to graph configuration file")
	createGraphCmd.Flags().String("output", "", "Path to output file")
	createGraphCmd.MarkFlagRequired("config")
	
	visualizeGraphCmd.Flags().String("config", "", "Path to graph configuration file")
	visualizeGraphCmd.Flags().String("output", "", "Path to output DOT file")
	visualizeGraphCmd.MarkFlagRequired("config")
}
