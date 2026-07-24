package cli

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"runtime/debug"
	"strings"

	"github.com/JetBrains/go-modern-guidelines/internal/goversion"
	"github.com/JetBrains/go-modern-guidelines/internal/guidelines"
)

const listVersionSourceConflictError = "list accepts only one Go version source: --go-version, --file-path, or one positional file path. Run list -h for usage"

var version = detectedVersion()

// Run executes the CLI with args, writes command output to stdout, and returns user-facing errors.
func Run(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return printUsage(stdout)
	}

	switch args[0] {
	case "-h", "--help", "help":
		return printUsage(stdout)
	case "--version", "version":
		return printVersion(stdout)
	case "list":
		return runList(args[1:], stdout)
	case "explain":
		return runExplain(args[1:], stdout)
	default:
		return fmt.Errorf("unknown command %q\n\nRun with --help for usage.", args[0])
	}
}

func detectedVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" || info.Main.Version == "(devel)" {
		return "dev"
	}
	return info.Main.Version
}

func runList(args []string, stdout io.Writer) error {
	var filePath string
	var goVersion string
	fs := newFlagSet("list", "list [--go-version <version> | --file-path <path> | <path>]")
	fs.StringVar(&filePath, "file-path", "", "Optional absolute or project-relative path to a Go file, go.mod, or go.work file.")
	fs.StringVar(&goVersion, "go-version", "", "Optional Go version override, for example 1.24, go1.24.3, or devel.")
	parsed, err := parseFlagSet(fs, args, stdout)
	if err != nil {
		return err
	}
	if !parsed {
		return nil
	}
	positionalArgs := fs.Args()
	if goVersion != "" {
		switch {
		case filePath != "", len(positionalArgs) > 0:
			return errors.New(listVersionSourceConflictError)
		}
	}

	switch {
	case len(positionalArgs) == 0:
	case filePath != "":
		return errors.New(listVersionSourceConflictError)
	case len(positionalArgs) == 1:
		filePath = positionalArgs[0]
	default:
		return errors.New("list accepts at most one positional file path. Run list -h for usage")
	}

	targetVersion, err := goversion.Resolve(filePath, goVersion, guidelines.LatestKnownVersion())
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(stdout, guidelines.ListText(targetVersion))
	return err
}

func runExplain(args []string, stdout io.Writer) error {
	var guidelineIDs stringListFlag
	fs := newFlagSet("explain", "explain [--guideline-id <id>]... <id>...")
	fs.Var(&guidelineIDs, "guideline-id", "Guideline id from list output. May be repeated or comma-separated.")
	fs.Var(&guidelineIDs, "guideline-ids", "Alias for --guideline-id.")
	parsed, err := parseFlagSet(fs, args, stdout)
	if err != nil {
		return err
	}
	if !parsed {
		return nil
	}
	for _, arg := range fs.Args() {
		if err := guidelineIDs.Set(arg); err != nil {
			return err
		}
	}

	output, err := guidelines.ExplainText(guidelineIDs)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(stdout, output)
	return err
}

func newFlagSet(name, usage string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s:\n  %s\n\nFlags:\n", name, usage)
		fs.PrintDefaults()
	}
	return fs
}

func parseFlagSet(fs *flag.FlagSet, args []string, stdout io.Writer) (bool, error) {
	var output bytes.Buffer
	fs.SetOutput(&output)

	err := fs.Parse(args)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, flag.ErrHelp) {
		_, err := fmt.Fprint(stdout, output.String())
		return false, err
	}
	if output.Len() > 0 {
		return false, errors.New(strings.TrimRight(output.String(), "\n"))
	}
	return false, err
}

func printUsage(w io.Writer) error {
	_, err := fmt.Fprint(w, `go-modern-guidelines provides modern Go coding guidelines for AI agents, so they can write up-to-date Go code despite their knowledge cutoff.

Commands:
  list [--go-version <version> | --file-path <path> | <path>]
      Return short modern Go coding guidelines supported by the resolved Go version,
      ordered newest first.

  explain [--guideline-id <id>]... <id>...
      Return detailed guidance and before/after examples for specific guideline ids.
`)
	return err
}

func printVersion(w io.Writer) error {
	_, err := fmt.Fprintln(w, version)
	return err
}

type stringListFlag []string

func (f *stringListFlag) String() string {
	return strings.Join(*f, ",")
}

func (f *stringListFlag) Set(value string) error {
	for part := range strings.SplitSeq(value, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			*f = append(*f, part)
		}
	}
	return nil
}
