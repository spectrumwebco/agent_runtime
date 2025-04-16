package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/webscraper/colly"
	collyModule "github.com/spectrumwebco/agent_runtime/pkg/modules/colly"
)

func NewCollyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "colly",
		Short: "Colly web scraping framework",
		Long:  `Colly provides a fast and elegant web scraping framework for extracting structured data from websites.`,
	}

	cmd.AddCommand(newCollyScrapeCommand())
	cmd.AddCommand(newCollyExampleCommand())
	cmd.AddCommand(newCollyExtractCommand())

	return cmd
}

func newCollyScrapeCommand() *cobra.Command {
	var (
		timeout        int
		maxDepth       int
		allowedDomains string
		followLinks    bool
		randomUA       bool
		outputFile     string
		cacheDir       string
	)

	cmd := &cobra.Command{
		Use:   "scrape [url]",
		Short: "Scrape a website",
		Long:  `Scrape a website and extract links, images, and other content.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := collyModule.NewModule(cfg)
			
			if timeout > 0 {
				module.WithTimeout(time.Duration(timeout) * time.Second)
			}
			
			if randomUA {
				module.WithRandomUserAgent()
			}
			
			if cacheDir != "" {
				if err := module.EnableCaching(cacheDir); err != nil {
					return fmt.Errorf("failed to enable caching: %w", err)
				}
			}
			
			config := &colly.ScrapeConfig{
				MaxDepth:        maxDepth,
				FollowLinks:     followLinks,
				RandomUserAgent: randomUA,
				Timeout:         time.Duration(timeout) * time.Second,
			}
			
			if allowedDomains != "" {
				config.AllowedDomains = strings.Split(allowedDomains, ",")
				module.AllowedDomains(config.AllowedDomains...)
			}
			
			config.TextSelectors = map[string]string{
				"title":       "title",
				"h1":          "h1",
				"description": "meta[name=description]",
				"keywords":    "meta[name=keywords]",
			}
			
			fmt.Printf("Scraping %s...\n", args[0])
			
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()
			
			result, err := module.Scrape(ctx, args[0], config)
			if err != nil {
				return fmt.Errorf("failed to scrape: %w", err)
			}
			
			fmt.Printf("Scrape completed in %s\n", result.Duration)
			fmt.Printf("Found %d links\n", len(result.Links))
			fmt.Printf("Found %d images\n", len(result.Images))
			
			if len(result.Errors) > 0 {
				fmt.Printf("Encountered %d errors\n", len(result.Errors))
				for _, err := range result.Errors {
					fmt.Printf("  - %s\n", err)
				}
			}
			
			if outputFile != "" {
				data, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal result: %w", err)
				}
				
				if err := os.WriteFile(outputFile, data, 0644); err != nil {
					return fmt.Errorf("failed to write output file: %w", err)
				}
				
				fmt.Printf("Results saved to %s\n", outputFile)
			}
			
			return nil
		},
	}

	cmd.Flags().IntVar(&timeout, "timeout", 30, "Timeout in seconds")
	cmd.Flags().IntVar(&maxDepth, "depth", 1, "Maximum crawl depth")
	cmd.Flags().StringVar(&allowedDomains, "domains", "", "Comma-separated list of allowed domains")
	cmd.Flags().BoolVar(&followLinks, "follow", false, "Follow links")
	cmd.Flags().BoolVar(&randomUA, "random-ua", false, "Use random user agent")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file for results (JSON)")
	cmd.Flags().StringVar(&cacheDir, "cache", "", "Cache directory")

	return cmd
}

func newCollyExampleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "example",
		Short: "Run a Colly example",
		Long:  `Run a simple Colly example to demonstrate the framework.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := collyModule.NewModule(cfg)
			module.RunExample()

			return nil
		},
	}

	return cmd
}

func newCollyExtractCommand() *cobra.Command {
	var (
		selector    string
		attribute   string
		outputFile  string
		limit       int
		format      string
		randomUA    bool
		followLinks bool
	)

	cmd := &cobra.Command{
		Use:   "extract [url]",
		Short: "Extract data from a website",
		Long:  `Extract specific data from a website using CSS selectors.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if selector == "" {
				return fmt.Errorf("selector is required")
			}
			
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := collyModule.NewModule(cfg)
			
			if randomUA {
				module.WithRandomUserAgent()
			}
			
			client := module.GetClient().GetCollector()
			
			type Result struct {
				URL       string `json:"url"`
				Text      string `json:"text,omitempty"`
				Attribute string `json:"attribute,omitempty"`
				Index     int    `json:"index"`
			}
			
			results := make([]Result, 0)
			count := 0
			
			client.OnHTML(selector, func(e *colly.HTMLElement) {
				if limit > 0 && count >= limit {
					return
				}
				
				result := Result{
					URL:   e.Request.URL.String(),
					Text:  e.Text,
					Index: count,
				}
				
				if attribute != "" {
					result.Attribute = e.Attr(attribute)
				}
				
				results = append(results, result)
				count++
				
				fmt.Printf("Found match %d: %s\n", count, strings.TrimSpace(e.Text))
				
				if attribute != "" {
					fmt.Printf("  Attribute %s: %s\n", attribute, e.Attr(attribute))
				}
				
				if followLinks && attribute == "href" {
					link := e.Request.AbsoluteURL(e.Attr("href"))
					if link != "" {
						e.Request.Visit(link)
					}
				}
			})
			
			client.OnRequest(func(r *colly.Request) {
				fmt.Println("Visiting", r.URL.String())
			})
			
			client.Visit(args[0])
			
			client.Wait()
			
			fmt.Printf("Extraction completed. Found %d matches.\n", len(results))
			
			if outputFile != "" {
				var data []byte
				var err error
				
				switch format {
				case "json":
					data, err = json.MarshalIndent(results, "", "  ")
				case "csv":
					var lines []string
					lines = append(lines, "url,text,attribute,index")
					for _, r := range results {
						line := fmt.Sprintf("%s,%s,%s,%d", 
							r.URL, 
							strings.ReplaceAll(r.Text, ",", " "), 
							strings.ReplaceAll(r.Attribute, ",", " "), 
							r.Index)
						lines = append(lines, line)
					}
					data = []byte(strings.Join(lines, "\n"))
				case "txt":
					var lines []string
					for _, r := range results {
						lines = append(lines, r.Text)
					}
					data = []byte(strings.Join(lines, "\n"))
				default:
					return fmt.Errorf("unsupported format: %s", format)
				}
				
				if err != nil {
					return fmt.Errorf("failed to marshal result: %w", err)
				}
				
				dir := filepath.Dir(outputFile)
				if dir != "." {
					if err := os.MkdirAll(dir, 0755); err != nil {
						return fmt.Errorf("failed to create directory: %w", err)
					}
				}
				
				if err := os.WriteFile(outputFile, data, 0644); err != nil {
					return fmt.Errorf("failed to write output file: %w", err)
				}
				
				fmt.Printf("Results saved to %s\n", outputFile)
			}
			
			return nil
		},
	}

	cmd.Flags().StringVarP(&selector, "selector", "s", "", "CSS selector to extract (required)")
	cmd.Flags().StringVarP(&attribute, "attribute", "a", "", "Attribute to extract (e.g., href, src)")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file for results")
	cmd.Flags().IntVarP(&limit, "limit", "l", 0, "Limit number of results (0 for unlimited)")
	cmd.Flags().StringVarP(&format, "format", "f", "json", "Output format (json, csv, txt)")
	cmd.Flags().BoolVar(&randomUA, "random-ua", false, "Use random user agent")
	cmd.Flags().BoolVar(&followLinks, "follow", false, "Follow links (only works with href attribute)")

	cmd.MarkFlagRequired("selector")

	return cmd
}
