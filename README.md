# colly-neo4j-storage

[![GitHub Action - Test Go Package](https://github.com/konsalex/colly-neo4j-storage/actions/workflows/workflow.yml/badge.svg)](https://github.com/konsalex/colly-neo4j-storage/actions/workflows/workflow.yml)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/konsalex/colly-neo4j-storage)
[![Go Report Card](https://goreportcard.com/badge/github.com/konsalex/colly-neo4j-storage)](https://goreportcard.com/report/github.com/konsalex/colly-neo4j-storage)

A Neo4j storage back end for the Colly web crawling/scraping framework https://go-colly.org

Example Usage:

```go
package main

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/konsalex/colly-neo4j-storage/colly/neo4j"
)

func main() {

	c := colly.NewCollector()

	storage := &neo4j.Storage{
        URI:      "bolt://localhost:7687",
        Username: "neo4j",
        Password: "password",
        Database: "colly" , // Override default database "neo4j" (optional)
	}

	if err := c.SetStorage(storage); err != nil {
		panic(err)
	}

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("http://go-colly.org/")
}
```
