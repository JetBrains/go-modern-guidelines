package guidelines

import (
	"fmt"
	"go/format"
	"runtime"
	"strings"
	"testing"

	"github.com/JetBrains/go-modern-guidelines/internal/goversion"
)

func TestListForResolvedVersion(t *testing.T) {
	output := ListText("1.26")
	if !strings.HasPrefix(output, "new_expression: Use new(value) for pointer fields or arguments instead of generic/type-specific pointer helper functions or temporary variables used only for &value; convert untyped constants to the target type first, because new(30) is a *int.") {
		t.Fatalf("list output is not newest-first:\n%s", output)
	}
	for _, want := range []string{
		"errors_as_type: Use errors.AsType[T](err) when checking whether an error matches a specific type.",
		"sync_waitgroup_go: Use wg.Go when spawning goroutines tracked by a sync.WaitGroup.",
		"cmp_or: Use cmp.Or to pick the first non-zero value from a fallback chain, but only when every argument is cheap and free of side effects, because all of them are evaluated.",
		"loopvar_capture: Do not add redundant loop-variable copies before closures or taking addresses; Go 1.22 gives each iteration its own variables.",
		"maps_keys_values_iter: Use maps.Keys or maps.Values directly as iterators instead of manually looping over a map.",
		"testing_t_context: Use t.Context() when a test function needs a context tied to the test lifetime; it is already canceled by the time t.Cleanup functions run, so do not pass it to cleanup work.",
		"json_omitzero: Use omitzero on JSON-tagged bool, numeric, struct, and time fields whose zero value should be omitted; keep omitempty for empty strings, slices, and maps.",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("list output missing %q\noutput:\n%s", want, output)
		}
	}
	cmpOrIndex := strings.Index(output, "cmp_or:")
	timeSinceIndex := strings.Index(output, "time_since:")
	if cmpOrIndex == -1 || timeSinceIndex == -1 || cmpOrIndex > timeSinceIndex {
		t.Fatalf("list output placed older guidelines before newer ones:\n%s", output)
	}
}

func TestListForGo127(t *testing.T) {
	output := ListText("1.27")
	if !strings.HasPrefix(output, "generic_methods: ") {
		t.Fatalf("list output is not newest-first for Go 1.27:\n%s", output)
	}
	for _, want := range []string{
		"generic_methods: Use generic methods instead of package-level generic helper functions when the operation naturally belongs to the type itself.",
		"promoted_field_literals: Set embedded struct fields directly with promoted field names in Go 1.27+ struct literals instead of constructing the embedded struct explicitly.",
		"strings_bytes_cut_last: Use strings.CutLast and bytes.CutLast instead of LastIndex plus manual slicing around the last separator.",
		"stdlib_uuid: Use the standard library uuid package instead of third-party libraries or custom UUID implementations when targeting Go 1.27+.",
		"url_clone: Use net/url URL.Clone and Values.Clone methods to copy URLs and URL values instead of manual copying.",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("list output missing %q\noutput:\n%s", want, output)
		}
	}
}

func TestGuidelineTableIntegrity(t *testing.T) {
	if len(modernGoGuidelines) == 0 {
		t.Fatal("modernGoGuidelines is empty")
	}

	latestVersion := LatestKnownVersion()
	seenIDs := make(map[string]bool, len(modernGoGuidelines))
	for i, guideline := range modernGoGuidelines {
		if i > 0 && goversion.Compare(modernGoGuidelines[i-1].sinceVersion, guideline.sinceVersion) < 0 {
			t.Fatalf(
				"modernGoGuidelines must be ordered newest-first: guideline %q with Go %s appears before %q with Go %s",
				modernGoGuidelines[i-1].id,
				modernGoGuidelines[i-1].sinceVersion,
				guideline.id,
				guideline.sinceVersion,
			)
		}

		if strings.TrimSpace(guideline.id) == "" {
			t.Fatalf("guideline %d has no id", i+1)
		}
		if guideline.id != strings.TrimSpace(guideline.id) {
			t.Fatalf("guideline %q id has leading or trailing whitespace", guideline.id)
		}
		if seenIDs[guideline.id] {
			t.Fatalf("guideline id %q is duplicated", guideline.id)
		}
		seenIDs[guideline.id] = true

		if strings.TrimSpace(guideline.sinceVersion) == "" {
			t.Fatalf("guideline %q has no since version", guideline.id)
		}
		if guideline.sinceVersion != strings.TrimSpace(guideline.sinceVersion) {
			t.Fatalf("guideline %q since version has leading or trailing whitespace", guideline.id)
		}
		if !goversion.IsMajorMinor(guideline.sinceVersion) {
			t.Fatalf("guideline %q since version %q must use major.minor format", guideline.id, guideline.sinceVersion)
		}
		if goversion.Compare(latestVersion, guideline.sinceVersion) < 0 {
			t.Fatalf("latest known guideline version %q is lower than guideline %q since version %q", latestVersion, guideline.id, guideline.sinceVersion)
		}

		if strings.TrimSpace(guideline.guideline) == "" {
			t.Fatalf("guideline %q has no guideline", guideline.id)
		}
		if strings.TrimSpace(guideline.details) == "" {
			t.Fatalf("guideline %q has no details", guideline.id)
		}
		if len(guideline.examples) == 0 {
			t.Fatalf("guideline %q has no examples", guideline.id)
		}
		for i, example := range guideline.examples {
			if strings.TrimSpace(example.before) == "" {
				t.Fatalf("guideline %q example %d has empty before snippet", guideline.id, i+1)
			}
			if strings.TrimSpace(example.after) == "" {
				t.Fatalf("guideline %q example %d has empty after snippet", guideline.id, i+1)
			}
		}
	}

	if len(guidelineByID) != len(modernGoGuidelines) {
		t.Fatalf("guidelineByID has %d entries, want %d", len(guidelineByID), len(modernGoGuidelines))
	}
	for id := range seenIDs {
		if _, ok := guidelineByID[id]; !ok {
			t.Fatalf("guidelineByID missing guideline %q", id)
		}
	}
}

func TestExplainFormatting(t *testing.T) {
	output, err := ExplainText([]string{"atomic_types"})
	if err != nil {
		t.Fatalf("explain guideline: %v", err)
	}

	want := `atomic_types:
  Since: Go 1.19

  Summary:
    Use typed atomics such as atomic.Bool, atomic.Int64, and atomic.Pointer[T] instead of untyped atomic functions in new code; converting existing fields changes struct layout, and atomic.Value and atomic.Pointer[T] differ in nil and type handling.

  Details:
    Typed atomic wrapper values keep the storage and the atomic operations together. They make the value type visible, reduce accidental non-atomic access, and avoid old pointer-alignment pitfalls. Prefer them in new code: replacing a raw int64 or unsafe.Pointer field in an existing type changes its layout and size, and atomic.Value differs from atomic.Pointer[T] in nil and type handling, so migrate a type deliberately rather than call by call.

  Examples:

  Example 1:

  Before:
    var enabled int32
    atomic.StoreInt32(&enabled, 1)
    if atomic.LoadInt32(&enabled) != 0 {
    	run()
    }

  After:
    var enabled atomic.Bool
    enabled.Store(true)
    if enabled.Load() {
    	run()
    }

  Example 2:

  Before:
    var ptr unsafe.Pointer
    atomic.StorePointer(&ptr, unsafe.Pointer(cfg))

  After:
    var ptr atomic.Pointer[Config]
    ptr.Store(cfg)`
	if output != want {
		t.Fatalf("explain output mismatch\nwant:\n%s\ngot:\n%s", want, output)
	}
}

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
func formatGoSnippet(snippet string) (string, error) {
	blocks := strings.Split(snippet, "\n\n")
	formatted := make([]string, 0, len(blocks))
	for _, block := range blocks {
		if out, err := formatGoDeclarations(block); err == nil {
			formatted = append(formatted, out)
			continue
		}
		out, err := formatGoStatements(block)
		if err != nil {
			return "", err
		}
		formatted = append(formatted, out)
	}
	return strings.Join(formatted, "\n\n"), nil
}

func formatGoDeclarations(block string) (string, error) {
	out, err := format.Source([]byte("package p\n" + block + "\n"))
	if err != nil {
		return "", err
	}
	lines := strings.Split(strings.TrimRight(string(out), "\n"), "\n")
	lines = lines[1:]
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	return strings.Join(lines, "\n"), nil
}

func formatGoStatements(block string) (string, error) {
	out, err := format.Source([]byte("package p\nfunc _() {\n" + block + "\n}\n"))
	if err != nil {
		return "", err
	}
	lines := strings.Split(strings.TrimRight(string(out), "\n"), "\n")
	start := -1
	for i, line := range lines {
		if line == "func _() {" {
			start = i + 1
			break
		}
	}
	if start < 0 || lines[len(lines)-1] != "}" {
		return "", fmt.Errorf("unexpected gofmt wrapper output:\n%s", out)
	}
	body := lines[start : len(lines)-1]
	for i, line := range body {
		body[i] = strings.TrimPrefix(line, "\t")
	}
	return strings.Join(body, "\n"), nil
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
