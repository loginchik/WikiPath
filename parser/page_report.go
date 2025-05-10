package parser

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type PageReport struct {
	Parent   *PageReport
	Url      string
	Children *[]string
	Depth    int
}

// Process requests URLs that are collected in page.
// If targetUrl is found among the collected URLs, then the process is stopped
// and final report is compiled. Otherwise, the queue is updated and parsing continues
func (p *PageReport) Process(
	ctx context.Context,
	sem chan struct{},
	wg *sync.WaitGroup,
	mu *sync.Mutex,
	processedUrls map[string]bool,
	reports *[]*PageReport,
	queue *[]*PageReport,
	finalReport **PageReport,
	targetUrl *string,
	pageCount *int,
	cancel context.CancelFunc,
) {
	// Ensure channel is unblocked when page processing is finished
	defer wg.Done()
	defer func() { <-sem }()

	select {
	case <-ctx.Done():
		// Break processing, when context is cancelled
		return
	default:
		// Continue processing
	}

	links := ProcessPage(p.Url)
	p.Children = links
	*reports = append(*reports, p)

	for _, link := range *links {
		link = RemoveLangReference(link)
		// Save final report, if target Url is found
		if link == *targetUrl {
			*finalReport = &PageReport{
				Url:    link,
				Parent: p,
				Depth:  p.Depth + 1,
			}
			cancel()
			return
		}
		// Add new URl to queue to request it next time
		mu.Lock()
		if !processedUrls[link] {
			child := &PageReport{
				Url:    link,
				Parent: p,
				Depth:  p.Depth + 1,
			}
			*queue = append(*queue, child)
		}
		mu.Unlock()
	}
	*pageCount++

	log.Debugf("[d %d] Processed %s, got %d links", p.Depth, p.Url, len(*links))
	fmt.Printf("\r[d %d] Processed %d links", p.Depth, *pageCount)
	time.Sleep(600 * time.Millisecond)
}
