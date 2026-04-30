package scraper

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type Metadata struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Links       []string `json:"links"`
	Images      []string `json:"images"`
	OGTitle     string   `json:"og_title,omitempty"`
	OGImage     string   `json:"og_image,omitempty"`
	OGURL       string   `json:"og_url,omitempty"`
}

type Scraper struct {
	client *http.Client
}

type ScraperError struct {
	URL string
	Err error
}

func (e *ScraperError) Error() string {
	return fmt.Sprintf("scraper error for %s: %v", e.URL, e.Err)
}

func New() *Scraper {
	return &Scraper{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *Scraper) Scrape(ctx context.Context, url string) (*Metadata, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "GoGetIt/0.1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	metadata := &Metadata{
		Links:  []string{},
		Images: []string{},
	}

	var parseWg sync.WaitGroup
	var mu sync.Mutex

	var parseNode func(*html.Node)
	parseNode = func(n *html.Node) {
		// Bail early if the caller's context has been cancelled (timeout
		// or interrupt). Without this, recursive goroutines keep walking
		// the tree long after the request has been abandoned.
		if ctx.Err() != nil {
			return
		}

		if n.Type == html.ElementNode {
			switch n.Data {
			case "title":
				if n.FirstChild != nil {
					mu.Lock()
					metadata.Title = n.FirstChild.Data
					mu.Unlock()
				}
			case "meta":
				s.parseMeta(n, metadata)
			case "a":
				if href := s.getAttribute(n, "href"); href != "" {
					mu.Lock()
					metadata.Links = append(metadata.Links, href)
					mu.Unlock()
				}
			case "img":
				if src := s.getAttribute(n, "src"); src != "" {
					mu.Lock()
					metadata.Images = append(metadata.Images, src)
					mu.Unlock()
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if ctx.Err() != nil {
				return
			}
			parseWg.Add(1)
			go func(child *html.Node) {
				defer parseWg.Done()
				parseNode(child)
			}(c)
		}
	}

	parseWg.Add(1)
	go func() {
		defer parseWg.Done()
		parseNode(doc)
	}()

	parseWg.Wait()

	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return metadata, nil
}

func (s *Scraper) parseMeta(n *html.Node, metadata *Metadata) {
	name := s.getAttribute(n, "name")
	property := s.getAttribute(n, "property")
	content := s.getAttribute(n, "content")

	if content == "" {
		return
	}

	switch {
	case name == "description":
		metadata.Description = content
	case property == "og:title":
		metadata.OGTitle = content
	case property == "og:image":
		metadata.OGImage = content
	case property == "og:url":
		metadata.OGURL = content
	}
}

func (s *Scraper) getAttribute(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func (s *Scraper) FetchLinkedResources(ctx context.Context, baseURL string, links []string) ([]*Metadata, error) {
	results := make([]*Metadata, len(links))
	errChan := make(chan error, len(links))

	var wg sync.WaitGroup

	for i, link := range links {
		wg.Add(1)
		go func(idx int, l string) {
			defer wg.Done()

			if !strings.HasPrefix(l, "http") {
				results[idx] = nil
				return
			}

			m, err := s.Scrape(ctx, l)
			if err != nil {
				errChan <- &ScraperError{URL: l, Err: err}
				return
			}
			results[idx] = m
		}(i, link)
	}

	wg.Wait()
	close(errChan)

	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return results, errors[0]
	}

	return results, nil
}
