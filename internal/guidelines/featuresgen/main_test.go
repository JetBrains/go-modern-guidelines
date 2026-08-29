package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TestFeaturesMarkdownInSync fails if FEATURES.md has drifted from the
// generator output. Run `make generate-features` to regenerate it.
func TestFeaturesMarkdownInSync(t *testing.T) {
	data, err := os.ReadFile("../guidelines.json")
	if err != nil {
		t.Fatal(err)
	}
	want, err := render(data)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile("../../../FEATURES.md")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatal("FEATURES.md is out of date; run `make generate-features`")
	}
}

func TestWriteCodeBlock(t *testing.T) {
	var b strings.Builder
	writeCodeBlock(&b, []string{`fmt.Println("hello")`})

	want := "```go\nfmt.Println(\"hello\")\n```\n"
	if got := b.String(); got != want {
		t.Fatalf("writeCodeBlock() = %q, want %q", got, want)
	}
}

func TestRenderPreservesValidatedIDInLink(t *testing.T) {
	output, err := render([]byte(`[{"id":"safe_id","since_version":"1.25","modernizer":true,"category":"Testing","impact":"High","guideline":"Use ` + "`" + `errors.AsType[T](err)` + "`" + `.","details":"Keep ` + "`" + `[T]` + "`" + ` as code.","examples":[{"before":["old()"],"after":["new()"]}]}]`))
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{
		"[`safe_id`](#safe_id)",
		"| Testing | [`safe_id`](#safe_id) | [x] | 1.25 | High |",
		"**Category: Testing · Go 1.25+ · Impact: High · Modernizer: yes**",
		"Use `errors.AsType[T](err)`.",
		"Keep `[T]` as code.",
	} {
		if !bytes.Contains(output, []byte(want)) {
			t.Fatalf("render() output does not contain %q", want)
		}
	}
}
