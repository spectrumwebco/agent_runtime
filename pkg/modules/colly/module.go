package colly

import (
	"context"
	"fmt"
	"time"

	"github.com/spectrumwebco/agent_runtime/internal/webscraper/colly"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Module struct {
	config *config.Config
	client *colly.Client
}

func NewModule(cfg *config.Config) *Module {
	return &Module{
		config: cfg,
		client: colly.NewClient(cfg),
	}
}

func (m *Module) Name() string {
	return "colly"
}

func (m *Module) Description() string {
	return "Web scraping framework for extracting structured data from websites"
}

func (m *Module) Initialize(ctx context.Context) error {
	return nil
}

func (m *Module) Shutdown(ctx context.Context) error {
	return nil
}

func (m *Module) GetClient() *colly.Client {
	return m.client
}

func (m *Module) WithTimeout(timeout time.Duration) *Module {
	m.client = colly.NewClient(m.config, colly.WithTimeout(timeout))
	return m
}

func (m *Module) WithUserAgent(userAgent string) *Module {
	m.client = colly.NewClient(m.config, colly.WithUserAgent(userAgent))
	return m
}

func (m *Module) WithRandomUserAgent() *Module {
	m.client = colly.NewClient(m.config, colly.WithRandomUserAgent())
	m.client.EnableRandomUserAgent()
	return m
}

func (m *Module) EnableCaching(cacheDir string) error {
	return m.client.EnableCaching(cacheDir)
}

func (m *Module) AllowedDomains(domains ...string) *Module {
	m.client.AllowedDomains(domains...)
	return m
}

func (m *Module) DisallowedDomains(domains ...string) *Module {
	m.client.DisallowedDomains(domains...)
	return m
}

func (m *Module) SetMaxDepth(depth int) *Module {
	m.client.SetMaxDepth(depth)
	return m
}

func (m *Module) SetParallelism(n int) *Module {
	m.client.SetParallelism(n)
	return m
}

func (m *Module) Scrape(ctx context.Context, url string, config *colly.ScrapeConfig) (*colly.ScrapeResult, error) {
	return m.client.Scrape(ctx, url, config)
}

func (m *Module) RunExample() {
	fmt.Println("Running Colly example...")
	
	client := m.client.GetCollector()
	
	client.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
	})
	
	client.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	
	client.Visit("https://en.wikipedia.org/wiki/Web_scraping")
	
	client.Wait()
	
	fmt.Println("Example completed!")
}

func (m *Module) CreateScrapeConfig() *colly.ScrapeConfig {
	return &colly.ScrapeConfig{
		MaxDepth:        3,
		FollowLinks:     true,
		RandomUserAgent: true,
		TextSelectors:   make(map[string]string),
	}
}

func (m *Module) AddTextSelector(config *colly.ScrapeConfig, name, selector string) *colly.ScrapeConfig {
	if config.TextSelectors == nil {
		config.TextSelectors = make(map[string]string)
	}
	config.TextSelectors[name] = selector
	return config
}
