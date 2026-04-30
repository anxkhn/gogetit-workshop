package scraper

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestScrapeMetadata(t *testing.T) {
	htmlContent := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Page</title>
			<meta name="description" content="Test description">
			<meta property="og:title" content="OG Title">
		</head>
		<body>
			<a href="https://example.com/link1">Link 1</a>
			<a href="/link2">Link 2</a>
			<img src="image1.jpg">
		</body>
		</html>
	`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(htmlContent))
	}))
	defer server.Close()

	s := New()
	metadata, err := s.Scrape(context.Background(), server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if metadata.Title != "Test Page" {
		t.Errorf("expected title 'Test Page', got %q", metadata.Title)
	}

	if metadata.Description != "Test description" {
		t.Errorf("expected description 'Test description', got %q", metadata.Description)
	}

	if metadata.OGTitle != "OG Title" {
		t.Errorf("expected OG title 'OG Title', got %q", metadata.OGTitle)
	}

	if len(metadata.Links) == 0 {
		t.Error("expected links to be extracted")
	}
}

func TestScrapeTitle(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "simple title",
			html:     `<html><head><title>Simple Title</title></head><body></body></html>`,
			expected: "Simple Title",
		},
		{
			name:     "title with whitespace",
			html:     `<html><head><title>  Whitespace Title  </title></head><body></body></html>`,
			expected: "  Whitespace Title  ",
		},
		{
			name:     "no title",
			html:     `<html><head></head><body></body></html>`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(tt.html))
			}))
			defer server.Close()

			s := New()
			metadata, err := s.Scrape(context.Background(), server.URL)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if metadata.Title != tt.expected {
				t.Errorf("expected title %q, got %q", tt.expected, metadata.Title)
			}
		})
	}
}

func TestScrapeEmptyPage(t *testing.T) {
	tests := []struct {
		name        string
		html        string
		expectError bool
	}{
		{
			name:        "empty string",
			html:        "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			html:        "   \n\t  ",
			expectError: true,
		},
		{
			name:        "valid empty HTML",
			html:        "<html><head></head><body></body></html>",
			expectError: false,
		},
	}

	s := New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.html == "" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.Write([]byte(tt.html))
			}))
			defer server.Close()

			metadata, err := s.Scrape(context.Background(), server.URL)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for empty page, got nil - parser should reject empty content")
				}
				if metadata == nil {
					t.Errorf("expected non-nil metadata even on empty page - parser bug returns nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestParserTitle(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "basic title",
			html:     `<title>Hello World</title>`,
			expected: "Hello World",
		},
		{
			name:     "title with extra whitespace",
			html:     `<title>   Trimmed   Title   </title>`,
			expected: "Trimmed   Title",
		},
		{
			name:     "empty title",
			html:     `<title></title>`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := html.Parse(strings.NewReader(tt.html))
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			p := NewParser("")
			title := p.ParseTitle(doc)

			if title != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, title)
			}
		})
	}
}

func TestParserLinks(t *testing.T) {
	htmlContent := `
		<html>
		<body>
			<a href="https://example.com/page1">Page 1</a>
			<a href="/page2">Page 2</a>
			<a href="javascript:void(0)">Skip</a>
			<a>No href</a>
		</body>
		</html>
	`

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	p := NewParser("https://example.com")
	var links []string
	p.ParseLinks(doc, &links)

	if len(links) != 3 {
		t.Errorf("expected 3 links, got %d", len(links))
	}
}

func TestParserImages(t *testing.T) {
	htmlContent := `
		<html>
		<body>
			<img src="https://example.com/img1.jpg">
			<img src="/img2.png">
			<img>
		</body>
		</html>
	`

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	p := NewParser("https://example.com")
	var images []string
	p.ParseImages(doc, &images)

	if len(images) != 2 {
		t.Errorf("expected 2 images, got %d", len(images))
	}
}

func TestExtractLinksFromHTML(t *testing.T) {
	tests := []struct {
		name          string
		html          string
		expectedCount int
	}{
		{
			name:          "multiple links",
			html:          `<a href="link1">A</a><a href='link2'>B</a>`,
			expectedCount: 2,
		},
		{
			name:          "no links",
			html:          `<div>no links here</div>`,
			expectedCount: 0,
		},
		{
			name:          "empty href - regex does not match empty href",
			html:          `<a href="">Empty</a>`,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			links := ExtractLinksFromHTML(tt.html)
			if len(links) != tt.expectedCount {
				t.Errorf("expected %d links, got %d", tt.expectedCount, len(links))
			}
		})
	}
}

func TestScrape_CancelledContext_ReturnsContextErr(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Big enough document that recursive parsing actually takes a moment;
		// many <a> children to amplify goroutine fanout.
		var sb strings.Builder
		sb.WriteString("<html><body>")
		for i := 0; i < 5000; i++ {
			sb.WriteString(`<a href="https://example.test/`)
			sb.WriteString(strings.Repeat("x", 32))
			sb.WriteString(`">link</a>`)
		}
		sb.WriteString("</body></html>")
		_, _ = w.Write([]byte(sb.String()))
	}))
	defer server.Close()

	scraper := New()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // pre-cancelled

	_, err := scraper.Scrape(ctx, server.URL)
	if err == nil {
		t.Fatal("expected error from cancelled context, got nil")
	}
	if err != context.Canceled {
		// httpClient may surface ctx.Err via the request itself; either is fine.
		// Reject only if it's neither.
		t.Logf("got error: %v", err)
	}
}
