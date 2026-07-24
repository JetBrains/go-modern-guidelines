package goversion

import (
	"cmp"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/mod/modfile"
)

var goVersionInText = regexp.MustCompile(`(?i)(?:^|\s)(?:go)?(\d+\.\d+)`)

// Resolve returns the normalized Go major.minor version for the given version source.
func Resolve(filePath, goVersion, develVersion string) (string, error) {
	explicitVersion := strings.TrimSpace(goVersion)
	if explicitVersion != "" {
		version, err := normalizeGoVersion(explicitVersion, develVersion)
		if err != nil {
			return "", fmt.Errorf("cannot parse Go version %q. Use values like 1.24, go1.24.3, or devel", goVersion)
		}
		return version, nil
	}

	filePath = strings.TrimSpace(filePath)
	if filePath != "" {
		return resolveGoVersionFromPath(filePath, develVersion)
	}

	return resolveGoToolVersion("local Go toolchain", develVersion)
}

// Compare compares two Go version strings.
func Compare(left, right string) int {
	leftVersion, leftOK := parseMajorMinor(left)
	rightVersion, rightOK := parseMajorMinor(right)
	if !leftOK || !rightOK {
		return strings.Compare(left, right)
	}
	if c := cmp.Compare(leftVersion.major, rightVersion.major); c != 0 {
		return c
	}
	return cmp.Compare(leftVersion.minor, rightVersion.minor)
}

// IsMajorMinor reports whether version uses the Go major.minor form.
func IsMajorMinor(version string) bool {
	if strings.Count(version, ".") != 1 {
		return false
	}
	_, ok := parseMajorMinor(version)
	return ok
}

func resolveGoVersionFromPath(filePath, develVersion string) (string, error) {
	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(absolutePath)
	if err != nil {
		return "", fmt.Errorf("cannot read --file-path %q: %w", filePath, err)
	}

	if !info.IsDir() && isModuleVersionFile(absolutePath) {
		if version, ok, err := parseGoDirective(absolutePath, develVersion); err != nil {
			return "", err
		} else if ok {
			return version, nil
		}
	}

	searchDir := absolutePath
	if !info.IsDir() {
		searchDir = filepath.Dir(absolutePath)
	}
	moduleVersion, hasModuleVersion, err := resolveGoVersionFromModuleFiles(searchDir, develVersion)
	if err != nil {
		return "", err
	}

	if hasModuleVersion {
		return moduleVersion, nil
	}
	return resolveGoToolVersion("Go toolchain for "+filePath, develVersion)
}

func resolveGoVersionFromModuleFiles(startDir, develVersion string) (string, bool, error) {
	if goMod := findUp(startDir, "go.mod"); goMod != "" {
		version, ok, err := parseGoDirective(goMod, develVersion)
		return version, ok, err
	}
	if goWork := findUp(startDir, "go.work"); goWork != "" {
		version, ok, err := parseGoDirective(goWork, develVersion)
		return version, ok, err
	}
	return "", false, nil
}

func isModuleVersionFile(path string) bool {
	name := filepath.Base(path)
	return name == "go.mod" || name == "go.work"
}

func findUp(startDir, name string) string {
	dir := filepath.Clean(startDir)
	for {
		candidate := filepath.Join(dir, name)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func parseGoDirective(path, develVersion string) (string, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false, fmt.Errorf("cannot read %s: %w", path, err)
	}

	switch filepath.Base(path) {
	case "go.mod":
		goMod, err := modfile.Parse(path, data, nil)
		if err != nil {
			return "", false, fmt.Errorf("cannot parse %s: %w", path, err)
		}
		if goMod.Go == nil {
			return "", false, nil
		}
		return parseGoDirectiveVersion(path, goMod.Go.Version, develVersion)
	case "go.work":
		goWork, err := modfile.ParseWork(path, data, nil)
		if err != nil {
			return "", false, fmt.Errorf("cannot parse %s: %w", path, err)
		}
		if goWork.Go == nil {
			return "", false, nil
		}
		return parseGoDirectiveVersion(path, goWork.Go.Version, develVersion)
	default:
		return "", false, fmt.Errorf("cannot parse Go version directive from unsupported file %s", path)
	}
}

func parseGoDirectiveVersion(path, rawVersion, develVersion string) (string, bool, error) {
	version, err := normalizeGoVersion(rawVersion, develVersion)
	if err != nil {
		return "", false, fmt.Errorf("cannot parse Go version directive in %s: %q", path, rawVersion)
	}
	return version, true, nil
}

func resolveGoToolVersion(source, develVersion string) (string, error) {
	output, err := exec.Command("go", "env", "GOVERSION").Output()
	if err != nil {
		return "", fmt.Errorf("cannot determine the Go version from %s. Pass --go-version explicitly, for example --go-version=1.24", source)
	}
	version, err := normalizeGoVersion(string(output), develVersion)
	if err != nil {
		return "", fmt.Errorf("cannot determine the Go version from %s. Pass --go-version explicitly, for example --go-version=1.24", source)
	}
	return version, nil
}

func normalizeGoVersion(rawVersion, develVersion string) (string, error) {
	trimmed := strings.TrimSpace(rawVersion)
	if strings.EqualFold(trimmed, "devel") {
		return develVersion, nil
	}

	match := goVersionInText.FindStringSubmatch(trimmed)
	if len(match) < 2 {
		return "", fmt.Errorf("cannot parse Go version %q", rawVersion)
	}
	majorMinor, ok := parseMajorMinor(match[1])
	if !ok {
		return "", fmt.Errorf("cannot parse Go version %q", rawVersion)
	}
	return fmt.Sprintf("%d.%d", majorMinor.major, majorMinor.minor), nil
}

type majorMinorVersion struct {
	major int
	minor int
}

func parseMajorMinor(version string) (majorMinorVersion, bool) {
	majorPart, minorPart, ok := strings.Cut(version, ".")
	if !ok {
		return majorMinorVersion{}, false
	}
	minorPart, _, _ = strings.Cut(minorPart, ".")
	major, err := strconv.Atoi(majorPart)
	if err != nil {
		return majorMinorVersion{}, false
	}
	minor, err := strconv.Atoi(minorPart)
	if err != nil {
		return majorMinorVersion{}, false
	}
	return majorMinorVersion{major: major, minor: minor}, true
}
