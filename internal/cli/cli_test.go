package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

var errWriteFailed = errors.New("write failed")

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errWriteFailed
}

func TestExplainRequiresIDs(t *testing.T) {
	var stdout bytes.Buffer

	err := Run([]string{"explain"}, &stdout)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "requires at least one guideline id") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInvalidSubcommandFlagReturnsSingleErrorWithUsage(t *testing.T) {
	var stdout bytes.Buffer

	err := Run([]string{"list", "--df"}, &stdout)
	if err == nil {
		t.Fatal("expected error")
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
	output := err.Error()
	if count := strings.Count(output, "flag provided but not defined: -df"); count != 1 {
		t.Fatalf("parse error appeared %d times, want once:\n%s", count, output)
	}
	if !strings.Contains(output, "Usage of list:") {
		t.Fatalf("parse error should include usage:\n%s", output)
	}
}

func TestVersion(t *testing.T) {
	oldVersion := version
	version = "v0.1.0"
	t.Cleanup(func() {
		version = oldVersion
	})

	var stdout bytes.Buffer

	err := Run([]string{"--version"}, &stdout)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stdout.String() != "v0.1.0\n" {
		t.Fatalf("version output = %q, want %q", stdout.String(), "v0.1.0\n")
	}
}

func TestListRejectsMultiplePositionalPaths(t *testing.T) {
	var stdout bytes.Buffer

	err := Run([]string{"list", "one.go", "two.go"}, &stdout)
	if err == nil {
		t.Fatal("expected error")
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
	if !strings.Contains(err.Error(), "at most one positional file path") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListRejectsFilePathFlagWithPositionalPath(t *testing.T) {
	var stdout bytes.Buffer

	err := Run([]string{"list", "--file-path", "one.go", "two.go"}, &stdout)
	if err == nil {
		t.Fatal("expected error")
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
	if !strings.Contains(err.Error(), "accepts only one Go version source") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListRejectsGoVersionWithFilePath(t *testing.T) {
	var stdout bytes.Buffer

	err := Run([]string{"list", "--go-version=1.26", "--file-path", "one.go"}, &stdout)
	if err == nil {
		t.Fatal("expected error")
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}
	if !strings.Contains(err.Error(), "accepts only one Go version source") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListRejectsGoVersionWithPositionalPath(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "single path", args: []string{"list", "--go-version=1.26", "one.go"}},
		{name: "multiple paths", args: []string{"list", "--go-version=1.26", "one.go", "two.go"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer

			err := Run(tt.args, &stdout)
			if err == nil {
				t.Fatal("expected error")
			}
			if stdout.Len() != 0 {
				t.Fatalf("unexpected stdout: %q", stdout.String())
			}
			if !strings.Contains(err.Error(), "accepts only one Go version source") {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestRunReturnsOutputWriteErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "top-level help", args: []string{"--help"}},
		{name: "version output", args: []string{"--version"}},
		{name: "subcommand help", args: []string{"list", "-h"}},
		{name: "list output", args: []string{"list", "--go-version=1.26"}},
		{name: "explain output", args: []string{"explain", "atomic_types"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Run(tt.args, failingWriter{})
			if !errors.Is(err, errWriteFailed) {
				t.Fatalf("Run() error = %v, want %v", err, errWriteFailed)
			}
		})
	}
}

func TestSubcommandHelpWritesUsageToStdout(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		want             string
		unexpectedOutput string
	}{
		{
			name:             "list",
			args:             []string{"list", "-h"},
			want:             "Usage of list:",
			unexpectedOutput: "atomic_types:",
		},
		{
			name:             "explain",
			args:             []string{"explain", "-h"},
			want:             "Usage of explain:",
			unexpectedOutput: "requires at least one guideline id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer

			err := Run(tt.args, &stdout)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			output := stdout.String()
			if !strings.Contains(output, tt.want) {
				t.Fatalf("help output missing usage:\n%s", output)
			}
			if tt.name == "list" && !strings.Contains(output, "list [--go-version <version> | --file-path <path> | <path>]") {
				t.Fatalf("list help missing positional path form:\n%s", output)
			}
			if strings.Contains(output, tt.unexpectedOutput) {
				t.Fatalf("help output included command output:\n%s", output)
			}
		})
	}
}

func TestUnknownGuidelineID(t *testing.T) {
	var stdout bytes.Buffer

	err := Run([]string{"explain", "--guideline-id=missing"}, &stdout)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "Unknown") && !strings.Contains(err.Error(), "unknown") {
		t.Fatalf("unexpected error: %v", err)
	}
}
