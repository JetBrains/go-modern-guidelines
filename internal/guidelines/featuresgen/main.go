package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/JetBrains/go-modern-guidelines/internal/guidelines/schema"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: featuresgen <guidelines.json> <FEATURES.md>")
	}

	data, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("read guidelines: %w", err)
	}

	output, err := render(data)
	if err != nil {
		return err
	}
	if err := os.WriteFile(args[1], output, 0o644); err != nil {
		return fmt.Errorf("write features: %w", err)
	}
	return nil
}

func render(data []byte) ([]byte, error) {
	guidelines, err := schema.Parse(data)
	if err != nil {
		return nil, err
	}

	var b strings.Builder
	b.WriteString("<!-- Generated from internal/guidelines/guidelines.json; DO NOT EDIT. -->\n\n")
	b.WriteString("# Modern Go Guidelines Explained\n\n")
	b.WriteString("This file provides a more detailed description of the features supported in the Go Modern Guidelines.\n\n")
	b.WriteString("**Modernizer legend:**\n\n")
	b.WriteString("- [x] — supported by a modernize analyzer\n")
	b.WriteString("- [ ] — no modernizer available\n\n")
	b.WriteString("**Impact legend:**\n\n")
	b.WriteString("- Critical — found in almost every project, dozens of occurrences\n")
	b.WriteString("- High — found often, 5–20 occurrences per project\n")
	b.WriteString("- Medium — found regularly, 1–5 occurrences per project\n")
	b.WriteString("- Low — found rarely or in specific code\n\n")
	b.WriteString("## Guidelines\n\n")
	b.WriteString("| Category | Guideline | Modernizer | Go | Impact |\n")
	b.WriteString("|----------|-----------|------------|----|--------|\n")
	for _, guideline := range guidelines {
		fmt.Fprintf(
			&b,
			"| %s | [`%s`](#%s) | %s | %s | %s |\n",
			guideline.Category,
			guideline.ID,
			guideline.ID,
			modernizerMark(*guideline.Modernizer),
			guideline.SinceVersion,
			guideline.Impact,
		)
	}

	for _, guideline := range guidelines {
		b.WriteString("\n---\n\n")
		fmt.Fprintf(&b, "<a id=\"%s\"></a>\n\n", guideline.ID)
		fmt.Fprintf(&b, "## `%s`\n\n", guideline.ID)
		fmt.Fprintf(
			&b,
			"**Category: %s · Go %s+ · Impact: %s · Modernizer: %s**\n\n",
			guideline.Category,
			guideline.SinceVersion,
			guideline.Impact,
			modernizerLabel(*guideline.Modernizer),
		)
		b.WriteString(guideline.Guideline)
		b.WriteString("\n\n")
		b.WriteString(guideline.Details)
		b.WriteByte('\n')

		for i, example := range guideline.Examples {
			b.WriteString("\n### Example")
			if len(guideline.Examples) > 1 {
				fmt.Fprintf(&b, " %d", i+1)
			}
			b.WriteString("\n\n**Before:**\n\n")
			writeCodeBlock(&b, example.Before)
			b.WriteString("\n**After:**\n\n")
			writeCodeBlock(&b, example.After)
		}
	}

	return []byte(b.String()), nil
}

func modernizerMark(supported bool) string {
	if supported {
		return "[x]"
	}
	return "[ ]"
}

func modernizerLabel(supported bool) string {
	if supported {
		return "yes"
	}
	return "no"
}

func writeCodeBlock(b *strings.Builder, lines []string) {
	b.WriteString("```go\n")
	for _, line := range lines {
		b.WriteString(line)
		b.WriteByte('\n')
	}
	b.WriteString("```\n")
}
