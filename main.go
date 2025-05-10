package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/loginchik/WikiPath/parser"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

const logo string = `
 __          ___ _    _ _____      _   _     
 \ \        / (_) |  (_)  __ \    | | | |    
  \ \  /\  / / _| | ___| |__) |_ _| |_| |__  
   \ \/  \/ / | | |/ / |  ___/ _' | __| '_ \ 
    \  /\  /  | |   <| | |  | (_| | |_| | | |
     \/  \/   |_|_|\_\_|_|   \__,_|\__|_| |_|
`

// getUserParams reads flags and validates the input data.
// The return value are the address of a string variables initialUrl, targetUrl,
// and int variable maxDepth
func getUserParams() (*string, *string, *int, error) {
	initialUrl := flag.String("start", "", "Wikipedia page to start with")
	targetUrl := flag.String("target", "", "Wikipedia page to look for")
	maxDepth := flag.Int("depth", 5, "Maximum search depth")
	flag.Parse()

	if *initialUrl == "" || *targetUrl == "" {
		flag.Usage()
		return initialUrl, targetUrl, maxDepth, errors.New("invalid params")
	}

	var err error
	*initialUrl, err = parser.ValidateWikiURl(strings.TrimSpace(*initialUrl))
	if err != nil {
		log.Fatalf("Invalid starting URL: <%s>. %s", *initialUrl, err)
	}
	*targetUrl, err = parser.ValidateWikiURl(strings.TrimSpace(*targetUrl))
	if err != nil {
		log.Fatalf("Invalid URL to look for: <%s>. %s", *targetUrl, err)
	}

	return initialUrl, targetUrl, maxDepth, nil
}

func main() {
	log.SetLevel(log.InfoLevel)

	// Get user specified params
	initialUrl, targetUrl, maxDepth, err := getUserParams()
	if err != nil {
		os.Exit(1)
	}

	// Max number of requests that can be processed in parallel
	maxConcurrency := 5

	// Print application info
	fmt.Println(logo)
	fmt.Println("Navigate through Wikipedia with minimum number of clicks")
	fmt.Println("Start URL: ", *initialUrl)
	fmt.Println("Target URL: ", *targetUrl)

	// Context to prevent additional requests when path is found
	ctx, cancel := context.WithCancel(context.Background())

	// Process request
	var report *parser.PageReport
	report, err = parser.WideSearch(initialUrl, targetUrl, &maxConcurrency, maxDepth, ctx, cancel)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(parser.PrintReport(report))
}
