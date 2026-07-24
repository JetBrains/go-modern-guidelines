package goversion

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizeGoVersion(t *testing.T) {
	tests := map[string]string{
		"1.24":                  "1.24",
		"go1.24.3":              "1.24",
		"Go1.25":                "1.25",
		"go version go1.26.0 x": "1.26",
	}

	for input, want := range tests {
		got, err := normalizeGoVersion(input, "1.27")
		if err != nil {
			t.Fatalf("normalizeGoVersion(%q): %v", input, err)
		}
		if got != want {
			t.Fatalf("normalizeGoVersion(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestNormalizeGoVersionDevelUsesConfiguredVersion(t *testing.T) {
	got, err := normalizeGoVersion("devel", "1.99")
	if err != nil {
		t.Fatalf("normalizeGoVersion(%q): %v", "devel", err)
	}

	if got != "1.99" {
		t.Fatalf("normalizeGoVersion(%q) = %q, want configured devel version %q", "devel", got, "1.99")
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		left  string
		right string
		want  int
	}{
		{left: "1.24", right: "1.24", want: 0},
		{left: "1.25", right: "1.24", want: 1},
		{left: "1.9", right: "1.10", want: -1},
		{left: "2.0", right: "1.99", want: 1},
	}

	for _, tt := range tests {
		got := Compare(tt.left, tt.right)
		if got < 0 && tt.want < 0 || got == 0 && tt.want == 0 || got > 0 && tt.want > 0 {
			continue
		}
		t.Fatalf("Compare(%q, %q) = %d, want sign %d", tt.left, tt.right, got, tt.want)
	}
}

func TestIsMajorMinor(t *testing.T) {
	tests := map[string]bool{
		"1.24":   true,
		"1":      false,
		"1.24.3": false,
		"1.x":    false,
		"go1.24": false,
	}

	for version, want := range tests {
		if got := IsMajorMinor(version); got != want {
			t.Fatalf("IsMajorMinor(%q) = %t, want %t", version, got, want)
		}
	}
}

func TestParseGoDirectiveFromGoMod(t *testing.T) {
	path := filepath.Join(t.TempDir(), "go.mod")
	writeFile(t, path, "module example.test/m\n\ngo 1.24.3\n")

	got, ok, err := parseGoDirective(path, "1.27")
	if err != nil {
		t.Fatalf("parseGoDirective: %v", err)
	}
	if !ok {
		t.Fatal("parseGoDirective did not find go directive")
	}
	if got != "1.24" {
		t.Fatalf("parseGoDirective = %q, want %q", got, "1.24")
	}
}

func TestParseGoDirectiveFromGoWork(t *testing.T) {
	path := filepath.Join(t.TempDir(), "go.work")
	writeFile(t, path, "go 1.25\n\nuse ./m\n")

	got, ok, err := parseGoDirective(path, "1.27")
	if err != nil {
		t.Fatalf("parseGoDirective: %v", err)
	}
	if !ok {
		t.Fatal("parseGoDirective did not find go directive")
	}
	if got != "1.25" {
		t.Fatalf("parseGoDirective = %q, want %q", got, "1.25")
	}
}

func TestParseGoDirectiveRejectsInvalidGoModVersion(t *testing.T) {
	path := filepath.Join(t.TempDir(), "go.mod")
	writeFile(t, path, "module example.test/m\n\ngo 1.24x\n")

	_, _, err := parseGoDirective(path, "1.27")
	if err == nil {
		t.Fatal("parseGoDirective succeeded for invalid go directive")
	}
	if !strings.Contains(err.Error(), "invalid go version") {
		t.Fatalf("parseGoDirective error = %q, want invalid go version", err)
	}
}

func TestResolveFromNearestGoMod(t *testing.T) {
	dir := t.TempDir()
	moduleDir := filepath.Join(dir, "module")
	pkgDir := filepath.Join(moduleDir, "pkg")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(moduleDir, "go.mod"), "module example.test/m\n\ngo 1.24\n")
	goFile := filepath.Join(pkgDir, "x.go")
	writeFile(t, goFile, "package pkg\n")

	got, err := Resolve(goFile, "", "1.27")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if got != "1.24" {
		t.Fatalf("Resolve = %q, want %q", got, "1.24")
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
