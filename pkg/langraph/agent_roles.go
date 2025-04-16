package langraph

import (
	"context"
	"fmt"
)

type AgentRoleDefinition struct {
	Role               AgentRole                `json:"role"`
	Name               string                   `json:"name"`
	Description        string                   `json:"description"`
	PrimaryCapabilities []AgentCapability       `json:"primary_capabilities"`
	SecondaryCapabilities []AgentCapability     `json:"secondary_capabilities"`
	ResponsibilityAreas []string                `json:"responsibility_areas"`
	DefaultTools       []string                 `json:"default_tools"`
	CommunicationPatterns []CommunicationPattern `json:"communication_patterns"`
}

type CommunicationPattern struct {
	TargetRole      AgentRole `json:"target_role"`
	Pattern         string    `json:"pattern"`
	Description     string    `json:"description"`
	TriggerCondition string    `json:"trigger_condition"`
}

var (
	FrontendAgentRole = AgentRoleDefinition{
		Role:        AgentRoleFrontend,
		Name:        "Frontend Agent",
		Description: "Specializes in UI/UX development, component design, and frontend architecture",
		PrimaryCapabilities: []AgentCapability{
			AgentCapabilityUIDesign,
			AgentCapabilityCodeGeneration,
			AgentCapabilityDocumentation,
		},
		SecondaryCapabilities: []AgentCapability{
			AgentCapabilityTesting,
			AgentCapabilityCodeReview,
		},
		ResponsibilityAreas: []string{
			"UI component design and implementation",
			"Frontend architecture planning",
			"User experience optimization",
			"Responsive design implementation",
			"Frontend performance optimization",
			"Frontend testing and validation",
			"Frontend documentation",
		},
		DefaultTools: []string{
			"ui_component_generator",
			"css_optimizer",
			"accessibility_checker",
			"responsive_design_validator",
			"frontend_performance_analyzer",
		},
		CommunicationPatterns: []CommunicationPattern{
			{
				TargetRole:      AgentRoleAppBuilder,
				Pattern:         "request_api_specification",
				Description:     "Request API specifications for frontend integration",
				TriggerCondition: "When implementing UI components that need to interact with backend APIs",
			},
			{
				TargetRole:      AgentRoleCodegen,
				Pattern:         "request_code_generation",
				Description:     "Request code generation for complex UI components",
				TriggerCondition: "When implementing complex UI components that require specialized code",
			},
			{
				TargetRole:      AgentRoleEngineering,
				Pattern:         "request_integration_support",
				Description:     "Request support for integrating frontend with backend systems",
				TriggerCondition: "When facing challenges with frontend-backend integration",
			},
		},
	}

	AppBuilderAgentRole = AgentRoleDefinition{
		Role:        AgentRoleAppBuilder,
		Name:        "App Builder Agent",
		Description: "Specializes in application assembly, API design, and system architecture",
		PrimaryCapabilities: []AgentCapability{
			AgentCapabilityAPIDesign,
			AgentCapabilityDatabaseDesign,
			AgentCapabilityDeployment,
		},
		SecondaryCapabilities: []AgentCapability{
			AgentCapabilityDocumentation,
			AgentCapabilityTesting,
		},
		ResponsibilityAreas: []string{
			"API design and implementation",
			"Database schema design",
			"System architecture planning",
			"Service integration",
			"Application deployment",
			"Performance optimization",
			"Scalability planning",
		},
		DefaultTools: []string{
			"api_designer",
			"database_schema_generator",
			"system_architecture_planner",
			"service_integrator",
			"deployment_manager",
		},
		CommunicationPatterns: []CommunicationPattern{
			{
				TargetRole:      AgentRoleFrontend,
				Pattern:         "provide_api_specification",
				Description:     "Provide API specifications for frontend integration",
				TriggerCondition: "When API specifications are requested or updated",
			},
			{
				TargetRole:      AgentRoleCodegen,
				Pattern:         "request_code_generation",
				Description:     "Request code generation for complex API implementations",
				TriggerCondition: "When implementing complex APIs that require specialized code",
			},
			{
				TargetRole:      AgentRoleEngineering,
				Pattern:         "request_integration_support",
				Description:     "Request support for integrating with backend systems",
				TriggerCondition: "When facing challenges with backend integration",
			},
		},
	}

	CodegenAgentRole = AgentRoleDefinition{
		Role:        AgentRoleCodegen,
		Name:        "Codegen Agent",
		Description: "Specializes in code generation, code review, and code optimization",
		PrimaryCapabilities: []AgentCapability{
			AgentCapabilityCodeGeneration,
			AgentCapabilityCodeReview,
			AgentCapabilityTesting,
		},
		SecondaryCapabilities: []AgentCapability{
			AgentCapabilityDocumentation,
			AgentCapabilityDeployment,
		},
		ResponsibilityAreas: []string{
			"Code generation for various languages",
			"Code review and quality assurance",
			"Code optimization",
			"Test generation",
			"Code documentation",
			"Code refactoring",
			"Code migration",
		},
		DefaultTools: []string{
			"code_generator",
			"code_reviewer",
			"code_optimizer",
			"test_generator",
			"documentation_generator",
		},
		CommunicationPatterns: []CommunicationPattern{
			{
				TargetRole:      AgentRoleFrontend,
				Pattern:         "provide_generated_code",
				Description:     "Provide generated code for frontend components",
				TriggerCondition: "When code generation for frontend components is requested",
			},
			{
				TargetRole:      AgentRoleAppBuilder,
				Pattern:         "provide_generated_code",
				Description:     "Provide generated code for API implementations",
				TriggerCondition: "When code generation for API implementations is requested",
			},
			{
				TargetRole:      AgentRoleEngineering,
				Pattern:         "request_code_review",
				Description:     "Request code review for generated code",
				TriggerCondition: "When code generation is complete and needs review",
			},
		},
	}

	EngineeringAgentRole = AgentRoleDefinition{
		Role:        AgentRoleEngineering,
		Name:        "Engineering Agent",
		Description: "Specializes in core development tasks, system integration, and technical problem-solving",
		PrimaryCapabilities: []AgentCapability{
			AgentCapabilityTesting,
			AgentCapabilityDeployment,
			AgentCapabilityCodeReview,
		},
		SecondaryCapabilities: []AgentCapability{
			AgentCapabilityCodeGeneration,
			AgentCapabilityDocumentation,
			AgentCapabilityDatabaseDesign,
		},
		ResponsibilityAreas: []string{
			"System integration",
			"Technical problem-solving",
			"Performance optimization",
			"Security implementation",
			"Code quality assurance",
			"DevOps implementation",
			"Technical architecture",
		},
		DefaultTools: []string{
			"system_integrator",
			"performance_optimizer",
			"security_analyzer",
			"code_quality_checker",
			"devops_manager",
		},
		CommunicationPatterns: []CommunicationPattern{
			{
				TargetRole:      AgentRoleFrontend,
				Pattern:         "provide_integration_support",
				Description:     "Provide support for frontend-backend integration",
				TriggerCondition: "When frontend-backend integration support is requested",
			},
			{
				TargetRole:      AgentRoleAppBuilder,
				Pattern:         "provide_integration_support",
				Description:     "Provide support for backend integration",
				TriggerCondition: "When backend integration support is requested",
			},
			{
				TargetRole:      AgentRoleCodegen,
				Pattern:         "provide_code_review",
				Description:     "Provide code review for generated code",
				TriggerCondition: "When code review is requested",
			},
		},
	}

	OrchestratorAgentRole = AgentRoleDefinition{
		Role:        AgentRoleOrchestrator,
		Name:        "Orchestrator Agent",
		Description: "Coordinates the activities of other agents, manages task allocation, and ensures overall system coherence",
		PrimaryCapabilities: []AgentCapability{
			AgentCapabilityOrchestration,
			AgentCapabilityPlanning,
		},
		SecondaryCapabilities: []AgentCapability{
			AgentCapabilityDocumentation,
		},
		ResponsibilityAreas: []string{
			"Task allocation and coordination",
			"Progress monitoring",
			"Conflict resolution",
			"Resource allocation",
			"System coherence maintenance",
			"Project planning",
			"Task prioritization",
		},
		DefaultTools: []string{
			"task_allocator",
			"progress_monitor",
			"conflict_resolver",
			"resource_allocator",
			"system_coherence_checker",
		},
		CommunicationPatterns: []CommunicationPattern{
			{
				TargetRole:      AgentRoleFrontend,
				Pattern:         "assign_task",
				Description:     "Assign frontend development tasks",
				TriggerCondition: "When frontend development tasks need to be assigned",
			},
			{
				TargetRole:      AgentRoleAppBuilder,
				Pattern:         "assign_task",
				Description:     "Assign application building tasks",
				TriggerCondition: "When application building tasks need to be assigned",
			},
			{
				TargetRole:      AgentRoleCodegen,
				Pattern:         "assign_task",
				Description:     "Assign code generation tasks",
				TriggerCondition: "When code generation tasks need to be assigned",
			},
			{
				TargetRole:      AgentRoleEngineering,
				Pattern:         "assign_task",
				Description:     "Assign engineering tasks",
				TriggerCondition: "When engineering tasks need to be assigned",
			},
		},
	}
)

func GetAgentRoleDefinition(role AgentRole) (AgentRoleDefinition, error) {
	switch role {
	case AgentRoleFrontend:
		return FrontendAgentRole, nil
	case AgentRoleAppBuilder:
		return AppBuilderAgentRole, nil
	case AgentRoleCodegen:
		return CodegenAgentRole, nil
	case AgentRoleEngineering:
		return EngineeringAgentRole, nil
	case AgentRoleOrchestrator:
		return OrchestratorAgentRole, nil
	default:
		return AgentRoleDefinition{}, fmt.Errorf("unknown agent role: %s", role)
	}
}

func CreateAgentWithRoleDefinition(roleDef AgentRoleDefinition, eventStream EventStream) (*Agent, error) {
	config := AgentConfig{
		Name:         roleDef.Name,
		Description:  roleDef.Description,
		Role:         roleDef.Role,
		Capabilities: append(roleDef.PrimaryCapabilities, roleDef.SecondaryCapabilities...),
		Metadata: map[string]interface{}{
			"responsibility_areas":    roleDef.ResponsibilityAreas,
			"primary_capabilities":    roleDef.PrimaryCapabilities,
			"secondary_capabilities":  roleDef.SecondaryCapabilities,
			"communication_patterns":  roleDef.CommunicationPatterns,
		},
	}

	agent := NewAgent(config, func(ctx context.Context, agent *Agent, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"status": "processed",
			"role":   string(agent.Config.Role),
			"input":  input,
		}, nil
	}, eventStream)

	for _, toolName := range roleDef.DefaultTools {
		agent.AddTool(toolName, fmt.Sprintf("Default tool for %s", toolName), func(ctx context.Context, agent *Agent, params map[string]interface{}) (map[string]interface{}, error) {
			return map[string]interface{}{
				"status":   "executed",
				"tool":     toolName,
				"params":   params,
				"agent_id": agent.Config.ID,
			}, nil
		}, map[string]interface{}{})
	}

	return agent, nil
}

func GetCommunicationPatterns(role AgentRole) ([]CommunicationPattern, error) {
	roleDef, err := GetAgentRoleDefinition(role)
	if err != nil {
		return nil, err
	}
	return roleDef.CommunicationPatterns, nil
}

func GetCommunicationPatternForTarget(role AgentRole, targetRole AgentRole) (CommunicationPattern, error) {
	patterns, err := GetCommunicationPatterns(role)
	if err != nil {
		return CommunicationPattern{}, err
	}

	for _, pattern := range patterns {
		if pattern.TargetRole == targetRole {
			return pattern, nil
		}
	}

	return CommunicationPattern{}, fmt.Errorf("no communication pattern found for target role: %s", targetRole)
}

func EstablishCommunicationChannel(sourceAgent, targetAgent *Agent) error {
	if sourceAgent == nil || targetAgent == nil {
		return fmt.Errorf("source or target agent is nil")
	}

	pattern, err := GetCommunicationPatternForTarget(sourceAgent.Config.Role, targetAgent.Config.Role)
	if err != nil {
		pattern = CommunicationPattern{
			TargetRole:      targetAgent.Config.Role,
			Pattern:         "default_communication",
			Description:     "Default communication pattern",
			TriggerCondition: "When communication is needed",
		}
	}

	sourceAgent.UpdateState(map[string]interface{}{
		fmt.Sprintf("communication_channel_%s", targetAgent.Config.ID): map[string]interface{}{
			"target_agent_id":   targetAgent.Config.ID,
			"target_agent_role": string(targetAgent.Config.Role),
			"pattern":           pattern.Pattern,
			"description":       pattern.Description,
			"trigger_condition": pattern.TriggerCondition,
		},
	})

	return nil
}
