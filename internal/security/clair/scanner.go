package clair

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Scanner struct {
	client *Client
	config *config.Config
}

func NewScanner(cfg *config.Config) (*Scanner, error) {
	client, err := NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Clair client: %v", err)
	}

	return &Scanner{
		client: client,
		config: cfg,
	}, nil
}

func (s *Scanner) ScanImage(ctx context.Context, imageRef string) (*ScanResult, error) {
	if _, err := name.ParseReference(imageRef); err != nil {
		return nil, fmt.Errorf("invalid image reference: %v", err)
	}

	report, err := s.client.ScanImage(ctx, imageRef)
	if err != nil {
		return nil, fmt.Errorf("failed to scan image: %v", err)
	}

	result := &ScanResult{
		ImageRef:        imageRef,
		Digest:          report.Digest,
		Vulnerabilities: report.Vulnerabilities,
		Summary:         report.Summary,
		ScanTime:        report.ScanTime,
		Status:          "completed",
	}

	return result, nil
}

func (s *Scanner) ScanImageAsync(ctx context.Context, imageRef string) (*ScanJob, error) {
	job := &ScanJob{
		ID:        generateJobID(),
		ImageRef:  imageRef,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	go func() {
		asyncCtx := context.Background()

		result, err := s.ScanImage(asyncCtx, imageRef)
		if err != nil {
			job.Status = "failed"
			job.Error = err.Error()
			return
		}

		job.Status = "completed"
		job.Result = result
		job.CompletedAt = time.Now()
	}()

	return job, nil
}

func (s *Scanner) GetScanJob(ctx context.Context, jobID string) (*ScanJob, error) {
	return &ScanJob{
		ID:       jobID,
		Status:   "unknown",
		ImageRef: "",
	}, nil
}

type ScanResult struct {
	ImageRef        string               `json:"image_ref"`
	Digest          string               `json:"digest"`
	Vulnerabilities []Vulnerability      `json:"vulnerabilities"`
	Summary         VulnerabilitySummary `json:"summary"`
	ScanTime        time.Time            `json:"scan_time"`
	Status          string               `json:"status"`
}

type ScanJob struct {
	ID          string      `json:"id"`
	ImageRef    string      `json:"image_ref"`
	Status      string      `json:"status"`
	Error       string      `json:"error,omitempty"`
	Result      *ScanResult `json:"result,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	CompletedAt time.Time   `json:"completed_at,omitempty"`
}

func generateJobID() string {
	return fmt.Sprintf("scan-%d", time.Now().UnixNano())
}
