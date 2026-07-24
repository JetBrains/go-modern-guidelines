package guidelines

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/JetBrains/go-modern-guidelines/internal/goversion"
)

type modernGoExample struct {
	before string
	after  string
}

type modernGoGuideline struct {
	id           string
	sinceVersion string
	guideline    string
	details      string
	examples     []modernGoExample
}

type modernGoExampleData struct {
	Before []string `json:"before"`
	After  []string `json:"after"`
}

type modernGoGuidelineData struct {
	ID           string                `json:"id"`
	SinceVersion string                `json:"since_version"`
	Guideline    string                `json:"guideline"`
	Details      string                `json:"details"`
	Examples     []modernGoExampleData `json:"examples"`
}

//go:embed guidelines.json
var modernGoGuidelinesJSON []byte

var modernGoGuidelines = mustLoadModernGoGuidelines()

func mustLoadModernGoGuidelines() []modernGoGuideline {
	var guidelineData []modernGoGuidelineData
	if err := json.Unmarshal(modernGoGuidelinesJSON, &guidelineData); err != nil {
		panic(fmt.Sprintf("parse embedded guidelines.json: %v", err))
	}

	guidelines := make([]modernGoGuideline, 0, len(guidelineData))
	for _, guideline := range guidelineData {
		examples := make([]modernGoExample, 0, len(guideline.Examples))
		for _, example := range guideline.Examples {
			examples = append(examples, modernGoExample{
				before: strings.Join(example.Before, "\n"),
				after:  strings.Join(example.After, "\n"),
			})
		}
		guidelines = append(guidelines, modernGoGuideline{
			id:           guideline.ID,
			sinceVersion: guideline.SinceVersion,
			guideline:    guideline.Guideline,
			details:      guideline.Details,
			examples:     examples,
		})
	}
	return guidelines
}

// ListText returns short guideline text for guidelines supported by targetGoVersion.
func ListText(targetGoVersion string) string {
	return toGuidelinesText(supportedGuidelines(targetGoVersion))
}

func toGuidelinesText(guidelines []modernGoGuideline) string {
	parts := make([]string, 0, len(guidelines))
	for _, guideline := range guidelines {
		parts = append(parts, guideline.id+": "+guideline.guideline)
	}
	return strings.Join(parts, "\n")
}

// ExplainText returns detailed guidance for the requested guideline IDs.
func ExplainText(guidelineIDs []string) (string, error) {
	resolvedGuidelines, err := resolveGuidelinesByID(guidelineIDs)
	if err != nil {
		return "", err
	}
	return toGuidelineDetailsText(resolvedGuidelines), nil
}

func toGuidelineDetailsText(guidelines []modernGoGuideline) string {
	var b strings.Builder
	for guidelineIndex, guideline := range guidelines {
		if guidelineIndex > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(guideline.id)
		b.WriteString(":\n")
		fmt.Fprintf(&b, "  Since: Go %s\n\n", guideline.sinceVersion)
		b.WriteString("  Summary:\n")
		b.WriteString(indentLines(guideline.guideline, "    "))
		b.WriteString("\n\n  Details:\n")
		b.WriteString(indentLines(guideline.details, "    "))
		writeExamples(&b, guideline.examples)
	}
	return b.String()
}

func writeExamples(b *strings.Builder, examples []modernGoExample) {
	if len(examples) == 0 {
		return
	}

	b.WriteString("\n\n  Examples:")
	for exampleIndex, example := range examples {
		if len(examples) > 1 {
			fmt.Fprintf(b, "\n\n  Example %d:", exampleIndex+1)
		}
		b.WriteString("\n\n  Before:\n")
		b.WriteString(indentLines(example.before, "    "))
		b.WriteString("\n\n  After:\n")
		b.WriteString(indentLines(example.after, "    "))
	}
}

func indentLines(s, indent string) string {
	var b strings.Builder
	first := true
	for line := range strings.SplitSeq(s, "\n") {
		if !first {
			b.WriteByte('\n')
		}
		first = false
		b.WriteString(indent)
		b.WriteString(line)
	}
	return b.String()
}

func resolveGuidelinesByID(guidelineIDs []string) ([]modernGoGuideline, error) {
	requestedIDs := normalizeGuidelineIDs(guidelineIDs)
	if len(requestedIDs) == 0 {
		return nil, fmt.Errorf("explain command requires at least one guideline id. Run list to list available ids")
	}

	var unknownIDs []string
	for _, id := range requestedIDs {
		if _, ok := guidelineByID[id]; !ok {
			unknownIDs = append(unknownIDs, id)
		}
	}
	if len(unknownIDs) > 0 {
		return nil, fmt.Errorf(
			"unknown Go modern code guideline ids: %s. Available ids: %s",
			strings.Join(unknownIDs, ", "),
			strings.Join(availableGuidelineIDs(), ", "),
		)
	}

	guidelines := make([]modernGoGuideline, 0, len(requestedIDs))
	for _, id := range requestedIDs {
		guidelines = append(guidelines, guidelineByID[id])
	}
	return guidelines, nil
}

func supportedGuidelines(targetGoVersion string) []modernGoGuideline {
	var guidelines []modernGoGuideline
	for _, guideline := range modernGoGuidelines {
		if goversion.Compare(targetGoVersion, guideline.sinceVersion) >= 0 {
			guidelines = append(guidelines, guideline)
		}
	}
	return guidelines
}

// LatestKnownVersion returns the newest Go version covered by the guideline table.
func LatestKnownVersion() string {
	latestVersion := "0.0"
	for _, guideline := range modernGoGuidelines {
		if goversion.Compare(guideline.sinceVersion, latestVersion) > 0 {
			latestVersion = guideline.sinceVersion
		}
	}
	return latestVersion
}

func normalizeGuidelineIDs(values []string) []string {
	seen := map[string]bool{}
	var result []string
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	return result
}

func availableGuidelineIDs() []string {
	ids := make([]string, 0, len(modernGoGuidelines))
	for _, guideline := range modernGoGuidelines {
		ids = append(ids, guideline.id)
	}
	return ids
}

var guidelineByID = func() map[string]modernGoGuideline {
	guidelines := make(map[string]modernGoGuideline, len(modernGoGuidelines))
	for _, guideline := range modernGoGuidelines {
		guidelines[guideline.id] = guideline
	}
	return guidelines
}()
