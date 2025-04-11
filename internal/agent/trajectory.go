package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func (a *DefaultAgent) saveTrajectory() error {
	if a.trajPath == "" {
		return fmt.Errorf("trajectory path is not set")
	}

	type SerializableTrajectory struct {
		Info       AgentInfo        `json:"info"`
		Trajectory []TrajectoryStep `json:"trajectory"`
		Timestamp  int64            `json:"timestamp"`
	}

	serializableTrajectory := SerializableTrajectory{
		Info:       a.Info,
		Trajectory: a.trajectory,
		Timestamp:  time.Now().Unix(),
	}

	data, err := json.MarshalIndent(serializableTrajectory, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal trajectory: %w", err)
	}

	err = os.WriteFile(a.trajPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write trajectory to file: %w", err)
	}

	return nil
}

func LoadTrajectory(trajPath string) (*AgentRunResult, error) {
	data, err := os.ReadFile(trajPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read trajectory file: %w", err)
	}

	var serializableTrajectory struct {
		Info       AgentInfo        `json:"info"`
		Trajectory []TrajectoryStep `json:"trajectory"`
		Timestamp  int64            `json:"timestamp"`
	}

	err = json.Unmarshal(data, &serializableTrajectory)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal trajectory: %w", err)
	}

	result := &AgentRunResult{
		Info:       serializableTrajectory.Info,
		Trajectory: serializableTrajectory.Trajectory,
	}

	return result, nil
}
