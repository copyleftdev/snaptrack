package crawl

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"

	"github.com/copyleftdev/snaptrack/pkg/capture"
	"github.com/copyleftdev/snaptrack/pkg/snapshot"
	"github.com/copyleftdev/snaptrack/pkg/store"
)

// CrawlerConfig holds settings for domain crawling.
type CrawlerConfig struct {
	MaxDepth    int
	Concurrency int
}

// CrawlDomain crawls all pages within the same FQDN as startURL, storing snapshots in db.
func CrawlDomain(startURL string, db store.DBInterface, cfg CrawlerConfig) error {
	parsed, err := url.Parse(startURL)
	if err != nil {
		return fmt.Errorf("invalid start URL %q: %w", startURL, err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("start URL must include scheme and host: %q", startURL)
	}

	baseDomain := fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)

	visited := make(map[string]bool)
	var mu sync.Mutex

	sem := make(chan struct{}, cfg.Concurrency)
	var wg sync.WaitGroup

	var crawl func(u string, depth int)
	crawl = func(u string, depth int) {
		defer wg.Done()
		if depth > cfg.MaxDepth {
			return
		}

		sem <- struct{}{}
		htmlContent, err := capture.CaptureHTML(u, 15*time.Second)
		<-sem

		if err != nil {
			fmt.Printf("[ERROR] capturing %s: %v\n", u, err)
			return
		}

		// Store or update snapshot
		if err := snapshot.StoreOrUpdateSnapshot(db, u, htmlContent); err != nil {
			fmt.Printf("[ERROR] storing snapshot for %s: %v\n", u, err)
		}

		links, parseErr := extractSameDomainLinks(htmlContent, baseDomain, u)
		if parseErr != nil {
			fmt.Printf("[ERROR] parsing links for %s: %v\n", u, parseErr)
			return
		}

		// for each link, if not visited, schedule a crawl
		for _, link := range links {
			mu.Lock()
			if !visited[link] {
				visited[link] = true
				mu.Unlock()

				wg.Add(1)
				go crawl(link, depth+1)
			} else {
				mu.Unlock()
			}
		}
	}

	// Start
	visited[startURL] = true
	wg.Add(1)
	go crawl(startURL, 0)
	wg.Wait()
	return nil
}

// extractSameDomainLinks returns absolute URLs that stay under baseDomain.
func extractSameDomainLinks(htmlContent, baseDomain, currentURL string) ([]string, error) {
	var links []string
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return links, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					abs, ok := makeAbsURL(currentURL, attr.Val)
					if ok && strings.HasPrefix(abs, baseDomain) {
						links = append(links, abs)
					}
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links, nil
}

func makeAbsURL(baseURL, href string) (string, bool) {
	bu, err := url.Parse(baseURL)
	if err != nil {
		return "", false
	}
	ref, err := url.Parse(href)
	if err != nil {
		return "", false
	}
	resolved := bu.ResolveReference(ref)
	return resolved.String(), true
}
