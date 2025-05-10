package parser

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"
)

// prepareRequest creates a new Http Request and sets its User-Agent header
func prepareRequest(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error preparing request:", err)
	}
	req.Header.Set("User-Agent", "wikipath-finder")
	return req
}

// GetResponse prepares request and gets response from Wikipedia.
// If server returns invalid response status code or throws an error,
// then fails with fatal log. Otherwise, returns response data
func GetResponse(url string) *http.Response {
	req := prepareRequest(url)
	client := http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error getting response:", err)
	} else if resp.StatusCode != 200 {
		log.Fatal("Invalid response status code:", resp.StatusCode)
	}
	return resp
}

// filterLinks removes duplicates and validates each link via ValidateURL method.
// If URL is valid, then adds it to resulting slice after formatting the value with FormatURL method
func filterLinks(links *[]string) *[]string {
	var filteredLinks []string
	for _, link := range *links {
		// Filter out duplicates
		if slices.Contains(filteredLinks, link) {
			continue
		}
		// Validate new link
		validated, err := ValidateURL(link)
		if err != nil {
			continue
		}
		filteredLinks = append(filteredLinks, FormatURL(validated))
	}
	return &filteredLinks
}

// ProcessPage gets the page data, parsers HTML to extract all links,
// then clears the slice of URLs from duplicates and returns the resulting slice
func ProcessPage(url string) *[]string {
	resp := GetResponse(url)
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal("Error parsing HTML:", err)
	}

	var links []string
	bodyContent := FindBodyContent(doc)
	FindAllLinks(bodyContent, &links)

	return filterLinks(&links)
}

// WideSearch continually requests pages until it finds the one that has a target URL
// or reaches max depth
func WideSearch(initialUrl *string, targetUrl *string, maxConcurrency *int, maxDepth *int, ctx context.Context, cancel context.CancelFunc) (*PageReport, error) {
	sem := make(chan struct{}, *maxConcurrency)
	var wg sync.WaitGroup

	var (
		reports       []*PageReport
		processedUrls = map[string]bool{}
		mu            sync.Mutex
	)
	queue := []*PageReport{{Url: *initialUrl, Depth: 0}}

	var finalReport *PageReport
	var pageCount int

	for len(queue) > 0 && finalReport == nil {
		currentLevel := queue
		queue = nil

		for _, page := range currentLevel {
			mu.Lock()
			// Skip already processed URLs
			if processedUrls[page.Url] || page.Depth > *maxDepth {
				mu.Unlock()
				continue
			}
			// Mark new URL as processed
			processedUrls[page.Url] = true
			mu.Unlock()

			// Block concurrency channel
			wg.Add(1)
			sem <- struct{}{}

			go page.Process(ctx, sem, &wg, &mu, processedUrls, &reports, &queue, &finalReport, targetUrl, &pageCount, cancel)
		}

		// Wait until all processes are finished
		wg.Wait()
	}
	fmt.Println("")

	if finalReport != nil {
		return finalReport, nil
	} else {
		return finalReport, errors.New("path not found. Try to increase max depth")
	}
}

// restorePath restores URL that lead to target URL.
// The return value is address of a string slice with URL
// ordered from starting to target URLs
func restorePath(report *PageReport) *[]string {
	var path []string
	for r := report; r != nil; r = r.Parent {
		path = append([]string{r.Url}, path...)
	}
	return &path
}

// PrintReport gets full path to target URL and prepares string report
func PrintReport(report *PageReport) string {
	path := *restorePath(report)
	reportParts := []string{
		fmt.Sprintf("It takes %d clicks:\n", len(path)-1),
		path[0],
	}
	for i, url := range path[1:] {
		reportParts = append(reportParts, fmt.Sprintf("-> (%d) %s", i+1, url))
	}
	return strings.Join(reportParts, " ")
}
