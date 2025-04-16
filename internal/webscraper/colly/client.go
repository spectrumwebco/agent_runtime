package colly

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Client struct {
	config     *config.Config
	collector  *colly.Collector
	httpClient *http.Client
	userAgent  string
	timeout    time.Duration
}

type ClientOption func(*Client)

func NewClient(cfg *config.Config, opts ...ClientOption) *Client {
	client := &Client{
		config:    cfg,
		timeout:   30 * time.Second,
		userAgent: "Kled.io Framework Scraper/1.0",
	}

	for _, opt := range opts {
		opt(client)
	}

	client.collector = colly.NewCollector(
		colly.UserAgent(client.userAgent),
		colly.MaxDepth(3),
		colly.Async(true),
	)

	client.collector.SetRequestTimeout(client.timeout)

	client.httpClient = &http.Client{
		Timeout: client.timeout,
	}

	return client
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = timeout
	}
}

func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

func WithRandomUserAgent() ClientOption {
	return func(c *Client) {
		c.userAgent = ""
	}
}

func (c *Client) GetCollector() *colly.Collector {
	return c.collector
}

func (c *Client) Clone() *colly.Collector {
	return c.collector.Clone()
}

func (c *Client) EnableCaching(cacheDir string) error {
	return c.collector.CacheDir(cacheDir)
}

func (c *Client) EnableRandomUserAgent() {
	extensions.RandomUserAgent(c.collector)
}

func (c *Client) EnableRandomMobileUserAgent() {
	extensions.RandomMobileUserAgent(c.collector)
}

func (c *Client) AllowedDomains(domains ...string) {
	c.collector.AllowedDomains = domains
}

func (c *Client) DisallowedDomains(domains ...string) {
	c.collector.DisallowedDomains = domains
}

func (c *Client) AllowURLRevisit() {
	c.collector.AllowURLRevisit = true
}

func (c *Client) DisableKeepAlive() {
	c.collector.DisableKeepAlive = true
}

func (c *Client) IgnoreRobotsTxt() {
	c.collector.IgnoreRobotsTxt = true
}

func (c *Client) SetMaxDepth(depth int) {
	c.collector.MaxDepth = depth
}

func (c *Client) SetMaxBodySize(size int) {
	c.collector.MaxBodySize = size
}

func (c *Client) SetParallelism(n int) {
	c.collector.Async = true
	c.collector.Limit(&colly.LimitRule{
		Parallelism: n,
	})
}

func (c *Client) LimitDomainGlob(pattern string, delay time.Duration, randomize bool, parallelism int) {
	c.collector.Limit(&colly.LimitRule{
		DomainGlob:  pattern,
		Delay:       delay,
		RandomDelay: randomize,
		Parallelism: parallelism,
	})
}

func (c *Client) OnHTML(selector string, callback func(*colly.HTMLElement)) {
	c.collector.OnHTML(selector, callback)
}

func (c *Client) OnXML(selector string, callback func(*colly.XMLElement)) {
	c.collector.OnXML(selector, callback)
}

func (c *Client) OnRequest(callback func(*colly.Request)) {
	c.collector.OnRequest(callback)
}

func (c *Client) OnResponse(callback func(*colly.Response)) {
	c.collector.OnResponse(callback)
}

func (c *Client) OnError(callback func(*colly.Response, error)) {
	c.collector.OnError(callback)
}

func (c *Client) OnScraped(callback func(*colly.Response)) {
	c.collector.OnScraped(callback)
}

func (c *Client) Visit(url string) error {
	return c.collector.Visit(url)
}

func (c *Client) VisitWithContext(ctx context.Context, url string) error {
	return c.collector.Request("GET", url, nil, ctx, nil)
}

func (c *Client) Wait() {
	c.collector.Wait()
}

func (c *Client) Scrape(ctx context.Context, url string, config *ScrapeConfig) (*ScrapeResult, error) {
	result := &ScrapeResult{
		URL:         url,
		StartTime:   time.Now(),
		Links:       make([]string, 0),
		Images:      make([]string, 0),
		Text:        make(map[string]string),
		Data:        make(map[string]interface{}),
		StatusCodes: make(map[string]int),
	}

	collector := c.collector.Clone()

	if config != nil {
		if len(config.AllowedDomains) > 0 {
			collector.AllowedDomains = config.AllowedDomains
		}
		
		if config.MaxDepth > 0 {
			collector.MaxDepth = config.MaxDepth
		}
		
		if config.Timeout > 0 {
			collector.SetRequestTimeout(config.Timeout)
		}
		
		if config.RandomUserAgent {
			extensions.RandomUserAgent(collector)
		}
	}

	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link != "" {
			result.Links = append(result.Links, link)
		}
		
		if config != nil && config.FollowLinks {
			e.Request.Visit(link)
		}
	})

	collector.OnHTML("img[src]", func(e *colly.HTMLElement) {
		img := e.Request.AbsoluteURL(e.Attr("src"))
		if img != "" {
			result.Images = append(result.Images, img)
		}
	})

	if config != nil && len(config.TextSelectors) > 0 {
		for name, selector := range config.TextSelectors {
			selectorName := name
			collector.OnHTML(selector, func(e *colly.HTMLElement) {
				result.Text[selectorName] = e.Text
			})
		}
	}

	collector.OnResponse(func(r *colly.Response) {
		result.StatusCodes[r.Request.URL.String()] = r.StatusCode
	})

	collector.OnError(func(r *colly.Response, err error) {
		result.Errors = append(result.Errors, fmt.Sprintf("Error on %s: %v", r.Request.URL, err))
	})

	err := collector.Visit(url)
	if err != nil {
		return result, err
	}

	collector.Wait()

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	return result, nil
}

type ScrapeConfig struct {
	AllowedDomains []string
	MaxDepth       int
	Timeout        time.Duration
	FollowLinks    bool
	RandomUserAgent bool
	TextSelectors  map[string]string
}

type ScrapeResult struct {
	URL         string
	StartTime   time.Time
	EndTime     time.Time
	Duration    time.Duration
	Links       []string
	Images      []string
	Text        map[string]string
	Data        map[string]interface{}
	StatusCodes map[string]int
	Errors      []string
}
