package guidelines

import (
	"go/format"
	"go/scanner"
	"go/token"
	"runtime"
	"strings"
	"testing"

	"github.com/JetBrains/go-modern-guidelines/internal/goversion"
)

// TestGuidelineExamplesAreFormattedGo checks that every example snippet parses as
// Go and is stored in gofmt form, so explain output is a usable style reference.
// Snippets are fragments, so each blank-line-separated block is formatted in
// whichever context parses: package-level declarations or function body statements.
func TestGuidelineExamplesAreFormattedGo(t *testing.T) {
	toolchainVersion, ok := toolchainLanguageVersion()
	if !ok {
		t.Skipf("cannot determine toolchain Go version from %q", runtime.Version())
	}

	skipped := 0
	for _, guideline := range modernGoGuidelines {
		if goversion.Compare(toolchainVersion, guideline.sinceVersion) < 0 {
			skipped += len(guideline.examples)
			continue
		}
		for i, example := range guideline.examples {
			for _, side := range []struct {
				name    string
				snippet string
			}{{"before", example.before}, {"after", example.after}} {
				formatted, err := formatGoSnippet(side.snippet)
				if err != nil {
					t.Errorf("guideline %q example %d %s does not parse as Go: %v\n%s",
						guideline.id, i+1, side.name, err, side.snippet)
					continue
				}
				if formatted != side.snippet {
					t.Errorf("guideline %q example %d %s is not gofmt-formatted\nhave:\n%s\nwant:\n%s",
						guideline.id, i+1, side.name, side.snippet, formatted)
				}
			}
		}
	}
	if skipped > 0 {
		t.Logf("skipped %d example snippets newer than the Go %s toolchain", skipped, toolchainVersion)
	}
}

// formatGoSnippet returns the gofmt form of an example snippet.
//
// gofmt accepts a partial file, so a snippet that is a list of declarations or
// a list of statements formats in one pass. A snippet that mixes the two is
// neither, so those are formatted one top-level part at a time.
func formatGoSnippet(snippet string) (string, error) {
	formatted, err := format.Source([]byte(snippet))
	if err == nil {
		return string(formatted), nil
	}

	parts := splitTopLevelParts(snippet)
	if len(parts) < 2 {
		return "", err
	}
	formattedParts := make([]string, 0, len(parts))
	for _, part := range parts {
		formattedPart, partErr := format.Source([]byte(part))
		if partErr != nil {
			return "", partErr
		}
		formattedParts = append(formattedParts, string(formattedPart))
	}
	// gofmt collapses a run of blank lines between two parts into one, so
	// rejoining with a single blank line reports an over-spaced snippet.
	return strings.Join(formattedParts, "\n\n"), nil
}

// splitTopLevelParts splits a snippet at the blank lines that separate whole
// declarations or statements. A blank line nested inside brackets, a string
// literal, or a comment belongs to the part around it, so it never splits.
func splitTopLevelParts(snippet string) []string {
	nested := nestedLines(snippet)

	var parts []string
	var part []string
	for i, line := range strings.Split(snippet, "\n") {
		if strings.TrimSpace(line) == "" && !nested[i+1] {
			if len(part) > 0 {
				parts = append(parts, strings.Join(part, "\n"))
				part = nil
			}
			continue
		}
		part = append(part, line)
	}
	if len(part) > 0 {
		parts = append(parts, strings.Join(part, "\n"))
	}
	return parts
}

// nestedLines reports the one-based lines of snippet that sit inside brackets
// or inside a string literal or comment that spans more than one line.
func nestedLines(snippet string) map[int]bool {
	fileSet := token.NewFileSet()
	file := fileSet.AddFile("", fileSet.Base(), len(snippet))
	var lexer scanner.Scanner
	lexer.Init(file, []byte(snippet), func(token.Position, string) {}, scanner.ScanComments)

	nested := make(map[int]bool)
	depth := 0
	previousEnd := 1
	for {
		position, tok, literal := lexer.Scan()
		if tok == token.EOF {
			return nested
		}

		start := fileSet.Position(position).Line
		if depth > 0 {
			for line := previousEnd; line <= start; line++ {
				nested[line] = true
			}
		}
		// Only string literals and comments span lines. An inserted semicolon
		// also carries a newline literal, so counting every token would mark
		// the line after each statement as nested.
		end := start
		if tok == token.STRING || tok == token.COMMENT {
			end += strings.Count(literal, "\n")
		}
		for line := start + 1; line <= end; line++ {
			nested[line] = true
		}
		previousEnd = end

		switch tok {
		case token.LPAREN, token.LBRACK, token.LBRACE:
			depth++
		case token.RPAREN, token.RBRACK, token.RBRACE:
			if depth > 0 {
				depth--
			}
		}
	}
}

// toolchainLanguageVersion reports the major.minor Go version of the running
// toolchain, so examples for a newer Go than the test runs on are not checked.
func toolchainLanguageVersion() (string, bool) {
	version := runtime.Version()
	index := strings.Index(version, "go1.")
	if index < 0 {
		return "", false
	}
	parts := strings.SplitN(version[index+len("go"):], ".", 3)
	if len(parts) < 2 {
		return "", false
	}
	minor := parts[1]
	for i, r := range minor {
		if r < '0' || r > '9' {
			minor = minor[:i]
			break
		}
	}
	candidate := parts[0] + "." + minor
	if !goversion.IsMajorMinor(candidate) {
		return "", false
	}
	return candidate, true
}

// TestFormatGoSnippetKeepsFormattedSnippets covers snippets that are already in
// gofmt form. Blank lines and indentation inside a bracketed block, a string
// literal, or a comment belong to the code around them, so formatting must
// leave every one of these snippets byte for byte alone.
func TestFormatGoSnippetKeepsFormattedSnippets(t *testing.T) {
	snippets := map[string]string{
		"blank line inside a function body":         "func f() {\n\ta := 1\n\n\tb := 2\n\t_ = b\n}",
		"blank line between struct field groups":    "type T struct {\n\tA int\n\n\tB int\n}",
		"blank line between switch cases":           "switch x {\ncase 1:\n\tprintln(1)\n\ncase 2:\n\tprintln(2)\n}",
		"blank line between import groups":          "import (\n\t\"fmt\"\n\n\t\"golang.org/x/mod/modfile\"\n)",
		"blank line inside a raw string literal":    "const query = `\nSELECT 1\n\nSELECT 2\n`",
		"indented line inside a raw string literal": "query := `\n\tSELECT 1\n`",
		"declarations mixed with statements":        "type Set[T comparable] map[T]struct{}\n\nfunc Map[T comparable](s Set[T]) []T {\n\treturn nil\n}\n\nnames := Map(users)",
		"import group inside a mixed snippet":       "import (\n\t\"fmt\"\n\n\t\"example.com/x\"\n)\n\nvalue := x.New()",
	}

	for name, snippet := range snippets {
		t.Run(name, func(t *testing.T) {
			formatted, err := formatGoSnippet(snippet)
			if err != nil {
				t.Fatalf("formatGoSnippet returned an error for valid Go: %v\n%s", err, snippet)
			}
			if formatted != snippet {
				t.Errorf("formatGoSnippet rewrote an already formatted snippet\nhave:\n%s\nwant:\n%s", snippet, formatted)
			}
		})
	}
}

// TestFormatGoSnippetReformatsUnformattedSnippets covers snippets that gofmt
// would change, so the guideline test reports them instead of passing them.
func TestFormatGoSnippetReformatsUnformattedSnippets(t *testing.T) {
	cases := map[string]struct {
		snippet string
		want    string
	}{
		"extra spaces in a statement": {
			snippet: "x  :=  1",
			want:    "x := 1",
		},
		"space indented declaration": {
			snippet: "func f()  {\n  return\n}",
			want:    "func f() {\n\treturn\n}",
		},
		"repeated blank lines between statements": {
			snippet: "a := 1\n\n\nb := 2\n_ = b",
			want:    "a := 1\n\nb := 2\n_ = b",
		},
		"repeated blank lines between mixed parts": {
			snippet: "type T int\n\n\nvalue := T(1)",
			want:    "type T int\n\nvalue := T(1)",
		},
	}

	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			formatted, err := formatGoSnippet(testCase.snippet)
			if err != nil {
				t.Fatalf("formatGoSnippet: %v", err)
			}
			if formatted != testCase.want {
				t.Errorf("formatGoSnippet = %q, want %q", formatted, testCase.want)
			}
		})
	}
}

func TestFormatGoSnippetRejectsNonGo(t *testing.T) {
	if formatted, err := formatGoSnippet("this is not go"); err == nil {
		t.Errorf("formatGoSnippet accepted non-Go input and returned %q", formatted)
	}
}
