# Colly Web Scraping Framework Integration

This document describes the integration of [Colly](https://github.com/gocolly/colly) into the Kled.io Framework.

## Overview

Colly is a fast and elegant web scraping framework for Go that provides a clean interface to write any kind of crawler, scraper, or spider. This integration enables the Kled.io Framework to extract structured data from websites for a wide range of applications, including data mining, data processing, and archiving.

## Directory Structure

The integration follows the established pattern for the agent_runtime repository:

```
agent_runtime/
├── cmd/cli/commands/
│   └── colly.go               # CLI commands for Colly
├── internal/webscraper/colly/
│   └── client.go              # Client wrapper for Colly
└── pkg/modules/colly/
    └── module.go              # Module integration for the framework
```

## Features

The Colly integration provides the following capabilities:

1. **Web Scraping**
   - Extract links, images, and text from websites
   - Follow links to crawl entire websites
   - Filter content based on CSS selectors

2. **Data Extraction**
   - Extract structured data from HTML and XML
   - Convert extracted data to various formats (JSON, CSV, TXT)
   - Customize extraction with CSS selectors

3. **Scraping Control**
   - Limit crawl depth and concurrency
   - Set request delays and timeouts
   - Filter domains and URLs

4. **CLI Commands**
   - `kled colly scrape [url]`: Scrape a website for links, images, and content
   - `kled colly extract [url]`: Extract specific data using CSS selectors
   - `kled colly example`: Run a simple example to demonstrate the framework

## Integration with Multiple Container Runtimes

This integration is designed to work with Spectrum Web Co's infrastructure that supports multiple container runtimes including LXC, Podman, Docker, and Kata Containers. The web scraping capabilities can be deployed in any of these container environments, providing flexibility for different deployment scenarios.

## Usage

### Basic Example

```bash
# Scrape a website
kled colly scrape https://example.com --depth 2 --follow --random-ua -o results.json

# Extract specific data
kled colly extract https://example.com -s "h1" -o headings.json

# Extract links
kled colly extract https://example.com -s "a[href]" -a "href" -o links.json

# Run a simple example
kled colly example
```

### Configuration Options

The Colly integration supports various configuration options:

- `--depth`: Maximum crawl depth
- `--timeout`: Request timeout in seconds
- `--domains`: Comma-separated list of allowed domains
- `--follow`: Follow links during scraping
- `--random-ua`: Use random user agent
- `--output, -o`: Output file for results
- `--format, -f`: Output format (json, csv, txt)
- `--selector, -s`: CSS selector for extraction
- `--attribute, -a`: Attribute to extract
- `--limit, -l`: Limit number of results
- `--cache`: Cache directory for requests

## Integration with Other Framework Components

The Colly integration works seamlessly with other components of the Kled.io Framework:

1. **LangChain Integration**: Use extracted data as input for LLM processing
2. **Database Integration**: Store scraped data in databases using GORM
3. **Microservices**: Distribute scraping tasks across microservices
4. **Terminal UI**: Visualize scraping progress with BubbleTea

## Dependencies

- github.com/gocolly/colly/v2
- github.com/PuerkitoBio/goquery (used by Colly for HTML parsing)

## Future Enhancements

1. Add support for distributed scraping
2. Implement proxy rotation for large-scale scraping
3. Add integration with data processing pipelines
4. Enhance error handling and retry mechanisms
5. Add support for JavaScript rendering with headless browsers
