package clair

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/quay/claircore"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	config     *config.Config
}

func NewClient(cfg *config.Config) (*Client, error) {
	baseURL, err := url.Parse(cfg.Clair.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Clair URL: %v", err)
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
		config: cfg,
	}, nil
}

func (c *Client) ScanImage(ctx context.Context, imageRef string) (*VulnerabilityReport, error) {
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image reference: %v", err)
	}

	img, err := remote.Image(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %v", err)
	}

	digest, err := img.Digest()
	if err != nil {
		return nil, fmt.Errorf("failed to get image digest: %v", err)
	}

	indexReport := &claircore.IndexReport{
		Hash:    digest.String(),
		State:   claircore.IndexState{},
		Packages: []claircore.Package{},
		Distributions: []claircore.Distribution{},
		Repositories: []claircore.Repository{},
		Environments: []claircore.Environment{},
	}

	vulnReport, err := c.scanImage(ctx, indexReport)
	if err != nil {
		return nil, fmt.Errorf("failed to scan image: %v", err)
	}

	return convertVulnerabilityReport(vulnReport), nil
}

func (c *Client) scanImage(ctx context.Context, indexReport *claircore.IndexReport) (*claircore.VulnerabilityReport, error) {
	return &claircore.VulnerabilityReport{
		Hash:          indexReport.Hash,
		Vulnerabilities: []claircore.Vulnerability{},
		PackageVulnerabilities: map[string][]string{},
		DistributionVulnerabilities: map[string][]string{},
		RepositoryVulnerabilities: map[string][]string{},
	}, nil
}

type VulnerabilityReport struct {
	ImageRef        string                `json:"image_ref"`
	Digest          string                `json:"digest"`
	Vulnerabilities []Vulnerability       `json:"vulnerabilities"`
	Summary         VulnerabilitySummary  `json:"summary"`
	ScanTime        time.Time             `json:"scan_time"`
}

type Vulnerability struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"`
	FixedIn     string    `json:"fixed_in,omitempty"`
	Link        string    `json:"link,omitempty"`
	Package     string    `json:"package"`
	Version     string    `json:"version"`
	Layer       string    `json:"layer,omitempty"`
}

type VulnerabilitySummary struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
	Unknown  int `json:"unknown"`
	Total    int `json:"total"`
}

func convertVulnerabilityReport(report *claircore.VulnerabilityReport) *VulnerabilityReport {
	return &VulnerabilityReport{
		Digest:          report.Hash,
		Vulnerabilities: []Vulnerability{},
		Summary: VulnerabilitySummary{
			Total: len(report.Vulnerabilities),
		},
		ScanTime: time.Now(),
	}
}
