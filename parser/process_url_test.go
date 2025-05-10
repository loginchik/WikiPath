package parser

import "testing"

func TestRemoveLangReference(t *testing.T) {
	got := RemoveLangReference("https://en.wikipedia.org/wiki/Lavatory_Madeleine")
	want := "https://wikipedia.org/wiki/Lavatory_Madeleine"

	if got != want {
		t.Errorf("Expected: %s, but got: %s", want, got)
	}
}

func TestFormatURL(t *testing.T) {
	got := FormatURL("/test-url")
	want := "https://wikipedia.org/test-url"

	if got != want {
		t.Errorf("Expected: %s, but got: %s", want, got)
	}
}

func TestValidateURLNoPrefix(t *testing.T) {
	_, err := ValidateURL("/data/wiki")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestValidateURLSpecificType(t *testing.T) {
	_, err := ValidateURL("/wiki/type:value")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestValidateURL(t *testing.T) {
	url, err := ValidateURL("/wiki/some-page")
	if err != nil {
		t.Error("Expected success, got error:", err)
	}

	if url != "/wiki/some-page" {
		t.Errorf("Expected </wiki/some-page>, got: %s", url)
	}
}
