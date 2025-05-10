package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// WikiPrefix being present in URL means that the URL points to another Wikipedia page
const WikiPrefix string = "/wiki/"

// ValidateWikiURl checks if URL addresses Wiki page
func ValidateWikiURl(url string) (string, error) {
	url = RemoveLangReference(url)
	if !strings.HasPrefix(url, fmt.Sprintf("https://wikipedia.org%s", WikiPrefix)) {
		return "", errors.New("Invalid URL")
	}
	return url, nil
}

// ValidateURL checks URL to point on Wikipedia page and not to any other content
func ValidateURL(url string) (string, error) {
	// Check if URL points to Wikipedia
	if !strings.HasPrefix(url, WikiPrefix) {
		return "", errors.New("URL has no prefix")
	}

	// Check if URL doesn't point to file or map
	separation := strings.SplitN(url[len(WikiPrefix):], ":", 2)
	if len(separation) == 2 {
		return "", errors.New(fmt.Sprintf("URL points to specific type: %s", separation[0]))
	}

	return url, nil
}

// FormatURL joins provided href to Wikipedia host
func FormatURL(href string) string {
	return fmt.Sprintf("https://wikipedia.org%s", href)
}

// RemoveLangReference uses regular expression to remove language-specific domain from Wikipedia URL
func RemoveLangReference(url string) string {
	re := regexp.MustCompile(`^https?://[a-z]{1,3}\.wikipedia\.org`)
	return re.ReplaceAllString(url, "https://wikipedia.org")
}
