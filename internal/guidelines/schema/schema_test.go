package schema

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseValidGuidelines(t *testing.T) {
	data, err := json.Marshal([]Guideline{validGuideline()})
	if err != nil {
		t.Fatal(err)
	}

	guidelines, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(guidelines) != 1 || guidelines[0].ID != "test_guideline" {
		t.Fatalf("Parse() = %#v", guidelines)
	}
}

func TestParseRejectsInvalidGuideline(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Guideline)
		want   string
	}{
		{"id", func(g *Guideline) { g.ID = "Unsafe-ID" }, "invalid guideline id"},
		{"version", func(g *Guideline) { g.SinceVersion = "1.25.1" }, "invalid Go version"},
		{"modernizer", func(g *Guideline) { g.Modernizer = nil }, "has no modernizer flag"},
		{"category", func(g *Guideline) { g.Category = "" }, "has no category"},
		{"impact", func(g *Guideline) { g.Impact = "Huge" }, "invalid impact"},
		{"guideline", func(g *Guideline) { g.Guideline = " " }, "has no guideline text"},
		{"details", func(g *Guideline) { g.Details = " " }, "has no details"},
		{"examples", func(g *Guideline) { g.Examples = nil }, "has no examples"},
		{"before", func(g *Guideline) { g.Examples[0].Before = nil }, "empty before snippet"},
		{"after", func(g *Guideline) { g.Examples[0].After = nil }, "empty after snippet"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			guideline := validGuideline()
			test.mutate(&guideline)
			data, err := json.Marshal([]Guideline{guideline})
			if err != nil {
				t.Fatal(err)
			}

			_, err = Parse(data)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Parse() error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestParseRejectsDuplicateIDs(t *testing.T) {
	guideline := validGuideline()
	data, err := json.Marshal([]Guideline{guideline, guideline})
	if err != nil {
		t.Fatal(err)
	}

	_, err = Parse(data)
	if err == nil || !strings.Contains(err.Error(), "duplicated") {
		t.Fatalf("Parse() error = %v, want duplicate id error", err)
	}
}

func TestParseRejectsOldestFirstOrder(t *testing.T) {
	old := validGuideline()
	old.SinceVersion = "1.24"
	newer := validGuideline()
	newer.ID = "newer_guideline"
	data, err := json.Marshal([]Guideline{old, newer})
	if err != nil {
		t.Fatal(err)
	}

	_, err = Parse(data)
	if err == nil || !strings.Contains(err.Error(), "ordered newest-first") {
		t.Fatalf("Parse() error = %v, want ordering error", err)
	}
}

func TestParseRejectsInvalidJSON(t *testing.T) {
	_, err := Parse([]byte("{"))
	if err == nil || !strings.Contains(err.Error(), "parse guidelines JSON") {
		t.Fatalf("Parse() error = %v, want JSON error", err)
	}
}

func validGuideline() Guideline {
	modernizer := true
	return Guideline{
		ID:           "test_guideline",
		SinceVersion: "1.25",
		Modernizer:   &modernizer,
		Category:     "Testing",
		Impact:       "Low",
		Guideline:    "Use `new API`.",
		Details:      "Details.",
		Examples: []Example{{
			Before: []string{"old()"},
			After:  []string{"new()"},
		}},
	}
}
