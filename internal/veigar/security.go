package veigar

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

type SecurityReview struct {
	ID            string                 `json:"id"`
	Repository    string                 `json:"repository"`
	Branch        string                 `json:"branch"`
	PRID          string                 `json:"pr_id"`
	PRTitle       string                 `json:"pr_title"`
	PRAuthor      string                 `json:"pr_author"`
	Status        string                 `json:"status"`
	SeverityLevel string                 `json:"severity_level"`
	Summary       string                 `json:"summary"`
	CreatedAt     time.Time              `json:"created_at"`
	CompletedAt   *time.Time             `json:"completed_at"`
	Vulnerabilities []SecurityVulnerability `json:"vulnerabilities"`
	ComplianceIssues []ComplianceIssue   `json:"compliance_issues"`
	Metadata      map[string]interface{} `json:"metadata"`
}

type SecurityVulnerability struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"`
	File        string    `json:"file"`
	LineNumber  int       `json:"line_number"`
	Tool        string    `json:"tool"`
	CreatedAt   time.Time `json:"created_at"`
}

type ComplianceIssue struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"`
	Framework   string    `json:"framework"`
	File        string    `json:"file"`
	LineNumber  int       `json:"line_number"`
	CreatedAt   time.Time `json:"created_at"`
}

type SecurityIntegration struct {
	Framework *Framework

	reviews map[string]*SecurityReview

	mutex sync.RWMutex

	Context map[string]interface{}
}

var _ Integration = (*SecurityIntegration)(nil)

type SecurityIntegrationConfig struct {
	BlockMergeSeverities []string `json:"block_merge_severities"`
	WarnSeverities       []string `json:"warn_severities"`
	InfoSeverities       []string `json:"info_severities"`
	ComplianceFrameworks []string `json:"compliance_frameworks"`
}

func NewSecurityIntegration(framework *Framework, config SecurityIntegrationConfig) (*SecurityIntegration, error) {
	if framework == nil {
		return nil, fmt.Errorf("framework cannot be nil")
	}

	integration := &SecurityIntegration{
		Framework: framework,
		reviews:   make(map[string]*SecurityReview),
		Context:   make(map[string]interface{}),
	}

	if err := integration.registerSecurityTools(); err != nil {
		return nil, fmt.Errorf("failed to register security tools: %w", err)
	}

	if err := framework.RegisterIntegration(integration); err != nil {
		return nil, fmt.Errorf("failed to register security integration: %w", err)
	}

	if framework.Config.Debug {
		log.Println("Security integration created")
	}

	return integration, nil
}

func (s *SecurityIntegration) Name() string {
	return "security"
}

func (s *SecurityIntegration) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Framework.Config.Debug {
		log.Println("Security integration started")
	}

	return nil
}

func (s *SecurityIntegration) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Framework.Config.Debug {
		log.Println("Security integration stopped")
	}

	return nil
}

func (s *SecurityIntegration) registerSecurityTools() error {
	if err := s.Framework.RegisterSecurityTool("static_analysis", s.staticAnalysisTool); err != nil {
		return fmt.Errorf("failed to register static analysis tool: %w", err)
	}

	if err := s.Framework.RegisterSecurityTool("vulnerability_scan", s.vulnerabilityScanTool); err != nil {
		return fmt.Errorf("failed to register vulnerability scanning tool: %w", err)
	}

	if err := s.Framework.RegisterSecurityTool("compliance_check", s.complianceCheckTool); err != nil {
		return fmt.Errorf("failed to register compliance checking tool: %w", err)
	}

	if err := s.Framework.RegisterSecurityTool("security_review", s.securityReviewTool); err != nil {
		return fmt.Errorf("failed to register security review tool: %w", err)
	}

	if err := s.Framework.RegisterSecurityTool("security_status", s.securityStatusTool); err != nil {
		return fmt.Errorf("failed to register security status tool: %w", err)
	}

	return nil
}

func (s *SecurityIntegration) staticAnalysisTool(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	files, ok := params["files"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("files parameter must be an array")
	}

	depth, _ := params["depth"].(string)
	if depth == "" {
		depth = "standard"
	}

	findings := []map[string]interface{}{
		{
			"id":          "SA001",
			"title":       "Insecure function call",
			"description": "Use of insecure function that could lead to buffer overflow",
			"severity":    "high",
			"file":        "src/main.c",
			"line_number": 42,
			"tool":        "semgrep",
		},
		{
			"id":          "SA002",
			"title":       "SQL Injection vulnerability",
			"description": "Potential SQL injection in database query",
			"severity":    "critical",
			"file":        "src/database.py",
			"line_number": 123,
			"tool":        "bandit",
		},
	}

	summary := map[string]interface{}{
		"total_findings": len(findings),
		"critical":       1,
		"high":           1,
		"medium":         0,
		"low":            0,
		"info":           0,
	}

	return map[string]interface{}{
		"findings": findings,
		"summary":  summary,
	}, nil
}

func (s *SecurityIntegration) vulnerabilityScanTool(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	files, ok := params["files"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("files parameter must be an array")
	}

	depth, _ := params["depth"].(string)
	if depth == "" {
		depth = "standard"
	}

	vulnerabilities := []map[string]interface{}{
		{
			"id":          "VS001",
			"title":       "Outdated dependency",
			"description": "Using an outdated version of a dependency with known vulnerabilities",
			"severity":    "high",
			"file":        "requirements.txt",
			"line_number": 15,
			"tool":        "snyk",
		},
		{
			"id":          "VS002",
			"title":       "Insecure cryptographic algorithm",
			"description": "Using MD5 for cryptographic purposes",
			"severity":    "medium",
			"file":        "src/crypto.py",
			"line_number": 78,
			"tool":        "trivy",
		},
	}

	summary := map[string]interface{}{
		"total_vulnerabilities": len(vulnerabilities),
		"critical":              0,
		"high":                  1,
		"medium":                1,
		"low":                   0,
		"info":                  0,
	}

	return map[string]interface{}{
		"vulnerabilities": vulnerabilities,
		"summary":         summary,
	}, nil
}

func (s *SecurityIntegration) complianceCheckTool(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	files, ok := params["files"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("files parameter must be an array")
	}

	frameworksParam, _ := params["frameworks"].([]interface{})
	frameworks := make([]string, 0)
	for _, f := range frameworksParam {
		if framework, ok := f.(string); ok {
			frameworks = append(frameworks, framework)
		}
	}
	if len(frameworks) == 0 {
		frameworks = []string{"e8", "nist", "owasp"}
	}

	issues := []map[string]interface{}{
		{
			"id":          "C001",
			"title":       "Missing access control",
			"description": "Endpoint lacks proper access control checks",
			"severity":    "high",
			"framework":   "e8",
			"file":        "src/api.py",
			"line_number": 56,
		},
		{
			"id":          "C002",
			"title":       "Insecure data storage",
			"description": "Sensitive data stored without encryption",
			"severity":    "medium",
			"framework":   "nist",
			"file":        "src/storage.py",
			"line_number": 92,
		},
	}

	frameworkResults := map[string]interface{}{
		"e8": map[string]interface{}{
			"total_issues": 1,
			"critical":     0,
			"high":         1,
			"medium":       0,
			"low":          0,
			"info":         0,
		},
		"nist": map[string]interface{}{
			"total_issues": 1,
			"critical":     0,
			"high":         0,
			"medium":       1,
			"low":          0,
			"info":         0,
		},
		"owasp": map[string]interface{}{
			"total_issues": 0,
			"critical":     0,
			"high":         0,
			"medium":       0,
			"low":          0,
			"info":         0,
		},
	}

	summary := map[string]interface{}{
		"total_issues": len(issues),
		"critical":     0,
		"high":         1,
		"medium":       1,
		"low":          0,
		"info":         0,
	}

	return map[string]interface{}{
		"issues":     issues,
		"frameworks": frameworkResults,
		"summary":    summary,
	}, nil
}

func (s *SecurityIntegration) securityReviewTool(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	repository, _ := params["repository"].(string)
	branch, _ := params["branch"].(string)
	prID, _ := params["pr_id"].(string)
	files, _ := params["files"].([]interface{})

	if repository == "" || branch == "" || prID == "" {
		return nil, fmt.Errorf("repository, branch, and pr_id are required parameters")
	}

	reviewID := fmt.Sprintf("SR-%d", time.Now().Unix())

	review := &SecurityReview{
		ID:         reviewID,
		Repository: repository,
		Branch:     branch,
		PRID:       prID,
		PRTitle:    params["pr_title"].(string),
		PRAuthor:   params["pr_author"].(string),
		Status:     "running",
		CreatedAt:  time.Now(),
		Metadata:   make(map[string]interface{}),
	}

	s.mutex.Lock()
	s.reviews[reviewID] = review
	s.mutex.Unlock()

	staticAnalysisResult, err := s.staticAnalysisTool(ctx, map[string]interface{}{
		"files": files,
		"depth": "deep",
	})
	if err != nil {
		review.Status = "failed"
		review.Summary = fmt.Sprintf("Static analysis failed: %v", err)
		return map[string]interface{}{
			"status":  "error",
			"error":   err.Error(),
			"id":      reviewID,
		}, nil
	}

	vulnerabilityScanResult, err := s.vulnerabilityScanTool(ctx, map[string]interface{}{
		"files": files,
		"depth": "deep",
	})
	if err != nil {
		review.Status = "failed"
		review.Summary = fmt.Sprintf("Vulnerability scanning failed: %v", err)
		return map[string]interface{}{
			"status":  "error",
			"error":   err.Error(),
			"id":      reviewID,
		}, nil
	}

	complianceCheckResult, err := s.complianceCheckTool(ctx, map[string]interface{}{
		"files":      files,
		"frameworks": []string{"e8", "nist", "owasp"},
	})
	if err != nil {
		review.Status = "failed"
		review.Summary = fmt.Sprintf("Compliance checking failed: %v", err)
		return map[string]interface{}{
			"status":  "error",
			"error":   err.Error(),
			"id":      reviewID,
		}, nil
	}

	staticAnalysisFindings, _ := staticAnalysisResult.(map[string]interface{})["findings"].([]map[string]interface{})
	
	vulnerabilityScanVulnerabilities, _ := vulnerabilityScanResult.(map[string]interface{})["vulnerabilities"].([]map[string]interface{})
	
	complianceCheckIssues, _ := complianceCheckResult.(map[string]interface{})["issues"].([]map[string]interface{})

	for _, finding := range staticAnalysisFindings {
		vulnerability := SecurityVulnerability{
			ID:          finding["id"].(string),
			Title:       finding["title"].(string),
			Description: finding["description"].(string),
			Severity:    finding["severity"].(string),
			File:        finding["file"].(string),
			LineNumber:  finding["line_number"].(int),
			Tool:        finding["tool"].(string),
			CreatedAt:   time.Now(),
		}
		review.Vulnerabilities = append(review.Vulnerabilities, vulnerability)
	}

	for _, v := range vulnerabilityScanVulnerabilities {
		vulnerability := SecurityVulnerability{
			ID:          v["id"].(string),
			Title:       v["title"].(string),
			Description: v["description"].(string),
			Severity:    v["severity"].(string),
			File:        v["file"].(string),
			LineNumber:  v["line_number"].(int),
			Tool:        v["tool"].(string),
			CreatedAt:   time.Now(),
		}
		review.Vulnerabilities = append(review.Vulnerabilities, vulnerability)
	}

	for _, issue := range complianceCheckIssues {
		complianceIssue := ComplianceIssue{
			ID:          issue["id"].(string),
			Title:       issue["title"].(string),
			Description: issue["description"].(string),
			Severity:    issue["severity"].(string),
			Framework:   issue["framework"].(string),
			File:        issue["file"].(string),
			LineNumber:  issue["line_number"].(int),
			CreatedAt:   time.Now(),
		}
		review.ComplianceIssues = append(review.ComplianceIssues, complianceIssue)
	}

	severityLevel := "none"
	for _, vulnerability := range review.Vulnerabilities {
		if vulnerability.Severity == "critical" {
			severityLevel = "critical"
			break
		} else if vulnerability.Severity == "high" && severityLevel != "critical" {
			severityLevel = "high"
		} else if vulnerability.Severity == "medium" && severityLevel != "critical" && severityLevel != "high" {
			severityLevel = "medium"
		} else if vulnerability.Severity == "low" && severityLevel != "critical" && severityLevel != "high" && severityLevel != "medium" {
			severityLevel = "low"
		}
	}

	for _, issue := range review.ComplianceIssues {
		if issue.Severity == "critical" && severityLevel != "critical" {
			severityLevel = "critical"
			break
		} else if issue.Severity == "high" && severityLevel != "critical" {
			severityLevel = "high"
		} else if issue.Severity == "medium" && severityLevel != "critical" && severityLevel != "high" {
			severityLevel = "medium"
		} else if issue.Severity == "low" && severityLevel != "critical" && severityLevel != "high" && severityLevel != "medium" {
			severityLevel = "low"
		}
	}

	review.Status = "completed"
	review.SeverityLevel = severityLevel
	completedAt := time.Now()
	review.CompletedAt = &completedAt

	summary := fmt.Sprintf("Security review completed with severity level: %s\n", severityLevel)
	summary += fmt.Sprintf("Vulnerabilities: %d (Critical: %d, High: %d, Medium: %d, Low: %d)\n",
		len(review.Vulnerabilities),
		countVulnerabilitiesBySeverity(review.Vulnerabilities, "critical"),
		countVulnerabilitiesBySeverity(review.Vulnerabilities, "high"),
		countVulnerabilitiesBySeverity(review.Vulnerabilities, "medium"),
		countVulnerabilitiesBySeverity(review.Vulnerabilities, "low"))
	summary += fmt.Sprintf("Compliance Issues: %d (Critical: %d, High: %d, Medium: %d, Low: %d)\n",
		len(review.ComplianceIssues),
		countComplianceIssuesBySeverity(review.ComplianceIssues, "critical"),
		countComplianceIssuesBySeverity(review.ComplianceIssues, "high"),
		countComplianceIssuesBySeverity(review.ComplianceIssues, "medium"),
		countComplianceIssuesBySeverity(review.ComplianceIssues, "low"))

	review.Summary = summary

	s.Framework.PublishSecurityEvent("security_review_completed", "veigar", map[string]interface{}{
		"review_id":      reviewID,
		"severity_level": severityLevel,
		"repository":     repository,
		"branch":         branch,
		"pr_id":          prID,
	}, map[string]interface{}{
		"vulnerabilities_count": len(review.Vulnerabilities),
		"compliance_issues_count": len(review.ComplianceIssues),
	})

	return map[string]interface{}{
		"status":         "success",
		"id":             reviewID,
		"severity_level": severityLevel,
		"security_report": summary,
		"vulnerabilities": review.Vulnerabilities,
		"compliance_issues": review.ComplianceIssues,
	}, nil
}

func (s *SecurityIntegration) securityStatusTool(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	id, ok := params["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id parameter is required")
	}

	s.mutex.RLock()
	review, exists := s.reviews[id]
	s.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("security review with id %s not found", id)
	}

	return map[string]interface{}{
		"status":         "success",
		"id":             review.ID,
		"review_status":  review.Status,
		"severity_level": review.SeverityLevel,
		"summary":        review.Summary,
		"created_at":     review.CreatedAt,
		"completed_at":   review.CompletedAt,
		"vulnerabilities_count": len(review.Vulnerabilities),
		"compliance_issues_count": len(review.ComplianceIssues),
	}, nil
}

func (s *SecurityIntegration) GetSecurityReview(id string) (*SecurityReview, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	review, exists := s.reviews[id]
	if !exists {
		return nil, fmt.Errorf("security review with id %s not found", id)
	}

	return review, nil
}

func (s *SecurityIntegration) ListSecurityReviews() []*SecurityReview {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	reviews := make([]*SecurityReview, 0, len(s.reviews))
	for _, review := range s.reviews {
		reviews = append(reviews, review)
	}

	return reviews
}

func (s *SecurityIntegration) SaveState() ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	state := map[string]interface{}{
		"context": s.Context,
		"reviews": s.reviews,
	}

	return json.Marshal(state)
}

func (s *SecurityIntegration) LoadState(data []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var state map[string]interface{}
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("failed to unmarshal state: %w", err)
	}

	if context, ok := state["context"].(map[string]interface{}); ok {
		s.Context = context
	}

	if reviewsData, ok := state["reviews"].(map[string]interface{}); ok {
		for id, reviewData := range reviewsData {
			if reviewBytes, err := json.Marshal(reviewData); err == nil {
				var review SecurityReview
				if err := json.Unmarshal(reviewBytes, &review); err == nil {
					s.reviews[id] = &review
				}
			}
		}
	}

	return nil
}

func countVulnerabilitiesBySeverity(vulnerabilities []SecurityVulnerability, severity string) int {
	count := 0
	for _, v := range vulnerabilities {
		if v.Severity == severity {
			count++
		}
	}
	return count
}

func countComplianceIssuesBySeverity(issues []ComplianceIssue, severity string) int {
	count := 0
	for _, i := range issues {
		if i.Severity == severity {
			count++
		}
	}
	return count
}
