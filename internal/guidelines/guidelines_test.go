package guidelines

import (
	"strings"
	"testing"
)

func TestListForResolvedVersion(t *testing.T) {
	output := ListText("1.26")
	if !strings.HasPrefix(output, "new_expression: Use `new(value)` for pointer fields or arguments instead of generic/type-specific pointer helper functions or temporary variables used only for `&value`.") {
		t.Fatalf("list output is not newest-first:\n%s", output)
	}
	for _, want := range []string{
		"errors_as_type: Use `errors.AsType[T](err)` when checking whether an error matches a specific type.",
		"sync_waitgroup_go: Use `wg.Go` when spawning goroutines tracked by a `sync.WaitGroup`.",
		"cmp_or: Use `cmp.Or` to pick the first non-zero value from a fallback chain.",
		"loopvar_capture: Do not add redundant loop-variable copies before closures or taking addresses; Go 1.22 gives each iteration its own variables.",
		"maps_keys_values_iter: Use `maps.Keys` or `maps.Values` directly as iterators instead of manually looping over a map.",
		"testing_t_context: Use `t.Context()` when a test function needs a context tied to the test lifetime.",
		"json_omitzero: Use `omitzero` on JSON-tagged bool, numeric, struct, and time fields whose zero value should be omitted; keep `omitempty` for empty strings, slices, and maps.",
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
		"strings_bytes_cut_last: Use `strings.CutLast` and `bytes.CutLast` instead of `LastIndex` plus manual slicing around the last separator.",
		"stdlib_uuid: Use the standard library `uuid` package instead of third-party libraries or custom UUID implementations when targeting Go 1.27+.",
		"url_clone: Use the `URL.Clone` and `Values.Clone` methods from `net/url` to copy URLs and `URL` values instead of manual copying.",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("list output missing %q\noutput:\n%s", want, output)
		}
	}
}

func TestGuidelineIndexIntegrity(t *testing.T) {
	if len(guidelineByID) != len(modernGoGuidelines) {
		t.Fatalf("guidelineByID has %d entries, want %d", len(guidelineByID), len(modernGoGuidelines))
	}
	for _, guideline := range modernGoGuidelines {
		if indexed, ok := guidelineByID[guideline.id]; !ok || indexed.id != guideline.id {
			t.Fatalf("guidelineByID missing guideline %q", guideline.id)
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
    ` + "Use typed atomics such as `atomic.Bool`, `atomic.Int64`, and `atomic.Pointer[T]` instead of untyped atomic functions." + `

  Details:
    Typed atomic wrapper values keep the storage and the atomic operations together. They make the value type visible, reduce accidental non-atomic access, and avoid old pointer-alignment pitfalls.

  Examples:

  Example 1:

  Before:
    var enabled int32
    atomic.StoreInt32(&enabled, 1)
    if atomic.LoadInt32(&enabled) != 0 {
` + "    \trun()" + `
    }

  After:
    var enabled atomic.Bool
    enabled.Store(true)
    if enabled.Load() {
` + "    \trun()" + `
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
