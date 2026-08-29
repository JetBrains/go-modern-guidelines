package schema

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/JetBrains/go-modern-guidelines/internal/goversion"
)

// Example contains the before and after snippets for a guideline example.
type Example struct {
	Before []string `json:"before"`
	After  []string `json:"after"`
}

// Guideline is the JSON representation of a modern Go guideline.
type Guideline struct {
	ID           string    `json:"id"`
	SinceVersion string    `json:"since_version"`
	Modernizer   *bool     `json:"modernizer"`
	Category     string    `json:"category"`
	Impact       string    `json:"impact"`
	Guideline    string    `json:"guideline"`
	Details      string    `json:"details"`
	Examples     []Example `json:"examples"`
}

// Parse decodes and validates a guideline collection.
func Parse(data []byte) ([]Guideline, error) {
	var guidelines []Guideline
	if err := json.Unmarshal(data, &guidelines); err != nil {
		return nil, fmt.Errorf("parse guidelines JSON: %w", err)
	}
	if len(guidelines) == 0 {
		return nil, fmt.Errorf("guidelines are empty")
	}

	seenIDs := make(map[string]bool, len(guidelines))
	for i, guideline := range guidelines {
		if !validID(guideline.ID) {
			return nil, fmt.Errorf("invalid guideline id %q: use lowercase letters, digits, and underscores", guideline.ID)
		}
		if seenIDs[guideline.ID] {
			return nil, fmt.Errorf("guideline id %q is duplicated", guideline.ID)
		}
		seenIDs[guideline.ID] = true

		if !goversion.IsMajorMinor(guideline.SinceVersion) {
			return nil, fmt.Errorf("invalid Go version %q for guideline %q: use major.minor", guideline.SinceVersion, guideline.ID)
		}
		if i > 0 && goversion.Compare(guidelines[i-1].SinceVersion, guideline.SinceVersion) < 0 {
			return nil, fmt.Errorf(
				"guidelines must be ordered newest-first: %q with Go %s appears before %q with Go %s",
				guidelines[i-1].ID,
				guidelines[i-1].SinceVersion,
				guideline.ID,
				guideline.SinceVersion,
			)
		}
		if guideline.Modernizer == nil {
			return nil, fmt.Errorf("guideline %q has no modernizer flag", guideline.ID)
		}
		if guideline.Category == "" {
			return nil, fmt.Errorf("guideline %q has no category", guideline.ID)
		}
		if !validImpact(guideline.Impact) {
			return nil, fmt.Errorf("guideline %q has invalid impact %q", guideline.ID, guideline.Impact)
		}
		if strings.TrimSpace(guideline.Guideline) == "" {
			return nil, fmt.Errorf("guideline %q has no guideline text", guideline.ID)
		}
		if strings.TrimSpace(guideline.Details) == "" {
			return nil, fmt.Errorf("guideline %q has no details", guideline.ID)
		}
		if len(guideline.Examples) == 0 {
			return nil, fmt.Errorf("guideline %q has no examples", guideline.ID)
		}
		for exampleIndex, example := range guideline.Examples {
			if strings.TrimSpace(strings.Join(example.Before, "\n")) == "" {
				return nil, fmt.Errorf("guideline %q example %d has empty before snippet", guideline.ID, exampleIndex+1)
			}
			if strings.TrimSpace(strings.Join(example.After, "\n")) == "" {
				return nil, fmt.Errorf("guideline %q example %d has empty after snippet", guideline.ID, exampleIndex+1)
			}
		}
	}
	return guidelines, nil
}

func validID(id string) bool {
	if id == "" {
		return false
	}
	for _, r := range id {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '_' {
			return false
		}
	}
	return true
}

func validImpact(impact string) bool {
	switch impact {
	case "Critical", "High", "Medium", "Low":
		return true
	default:
		return false
	}
}
