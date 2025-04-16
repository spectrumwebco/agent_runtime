package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/security/clair"
	clairmodule "github.com/spectrumwebco/agent_runtime/pkg/modules/clair"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

func NewClairCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clair",
		Short: "Container vulnerability scanning",
		Long:  `Clair provides container vulnerability scanning capabilities for the Kled.io Framework.`,
	}

	cmd.AddCommand(newClairScanCommand())
	cmd.AddCommand(newClairReportCommand())
	cmd.AddCommand(newClairConfigCommand())

	return cmd
}

func newClairScanCommand() *cobra.Command {
	var outputFormat string
	var waitForResult bool

	cmd := &cobra.Command{
		Use:   "scan [image]",
		Short: "Scan a container image for vulnerabilities",
		Long:  `Scan a container image for vulnerabilities using Clair.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			imageRef := args[0]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load configuration: %v", err)
			}

			scanner, err := clair.NewScanner(cfg)
			if err != nil {
				return fmt.Errorf("failed to create scanner: %v", err)
			}

			if waitForResult {
				result, err := scanner.ScanImage(context.Background(), imageRef)
				if err != nil {
					return fmt.Errorf("failed to scan image: %v", err)
				}

				switch outputFormat {
				case "json":
					json, err := json.MarshalIndent(result, "", "  ")
					if err != nil {
						return fmt.Errorf("failed to marshal result: %v", err)
					}
					fmt.Println(string(json))
				case "summary":
					fmt.Printf("Image: %s\n", result.ImageRef)
					fmt.Printf("Digest: %s\n", result.Digest)
					fmt.Printf("Scan Time: %s\n", result.ScanTime.Format(time.RFC3339))
					fmt.Printf("Status: %s\n", result.Status)
					fmt.Printf("Vulnerabilities: %d total (%d critical, %d high, %d medium, %d low, %d unknown)\n",
						result.Summary.Total, result.Summary.Critical, result.Summary.High,
						result.Summary.Medium, result.Summary.Low, result.Summary.Unknown)
				default:
					return fmt.Errorf("unknown output format: %s", outputFormat)
				}
			} else {
				job, err := scanner.ScanImageAsync(context.Background(), imageRef)
				if err != nil {
					return fmt.Errorf("failed to scan image: %v", err)
				}

				switch outputFormat {
				case "json":
					json, err := json.MarshalIndent(job, "", "  ")
					if err != nil {
						return fmt.Errorf("failed to marshal job: %v", err)
					}
					fmt.Println(string(json))
				case "summary":
					fmt.Printf("Job ID: %s\n", job.ID)
					fmt.Printf("Image: %s\n", job.ImageRef)
					fmt.Printf("Status: %s\n", job.Status)
					fmt.Printf("Created At: %s\n", job.CreatedAt.Format(time.RFC3339))
				default:
					return fmt.Errorf("unknown output format: %s", outputFormat)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "summary", "Output format (json, summary)")
	cmd.Flags().BoolVarP(&waitForResult, "wait", "w", true, "Wait for scan result")

	return cmd
}

func newClairReportCommand() *cobra.Command {
	var outputFormat string
	var outputFile string

	cmd := &cobra.Command{
		Use:   "report [job-id]",
		Short: "Get a vulnerability report",
		Long:  `Get a vulnerability report for a scan job.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			jobID := args[0]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load configuration: %v", err)
			}

			scanner, err := clair.NewScanner(cfg)
			if err != nil {
				return fmt.Errorf("failed to create scanner: %v", err)
			}

			job, err := scanner.GetScanJob(context.Background(), jobID)
			if err != nil {
				return fmt.Errorf("failed to get job: %v", err)
			}

			if job.Status != "completed" {
				return fmt.Errorf("job is not completed: %s", job.Status)
			}

			if job.Result == nil {
				return fmt.Errorf("job has no result")
			}

			var output []byte
			switch outputFormat {
			case "json":
				output, err = json.MarshalIndent(job.Result, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal result: %v", err)
				}
			case "summary":
				summary := fmt.Sprintf("Image: %s\n", job.Result.ImageRef)
				summary += fmt.Sprintf("Digest: %s\n", job.Result.Digest)
				summary += fmt.Sprintf("Scan Time: %s\n", job.Result.ScanTime.Format(time.RFC3339))
				summary += fmt.Sprintf("Status: %s\n", job.Result.Status)
				summary += fmt.Sprintf("Vulnerabilities: %d total (%d critical, %d high, %d medium, %d low, %d unknown)\n",
					job.Result.Summary.Total, job.Result.Summary.Critical, job.Result.Summary.High,
					job.Result.Summary.Medium, job.Result.Summary.Low, job.Result.Summary.Unknown)
				output = []byte(summary)
			default:
				return fmt.Errorf("unknown output format: %s", outputFormat)
			}

			if outputFile == "" {
				fmt.Println(string(output))
			} else {
				if err := os.WriteFile(outputFile, output, 0644); err != nil {
					return fmt.Errorf("failed to write output file: %v", err)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "summary", "Output format (json, summary)")
	cmd.Flags().StringVarP(&outputFile, "file", "f", "", "Output file")

	return cmd
}

func newClairConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure Clair",
		Long:  `Configure Clair for container vulnerability scanning.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load configuration: %v", err)
			}

			fmt.Println("Current Clair Configuration:")
			fmt.Printf("URL: %s\n", cfg.Clair.URL)
			fmt.Printf("Timeout: %d seconds\n", cfg.Clair.Timeout)

			return nil
		},
	}

	return cmd
}
