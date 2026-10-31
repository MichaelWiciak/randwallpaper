package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestRunVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer

	if err := run([]string{"-version"}, &stdout, &stderr); err != nil {
		t.Fatalf("run(-version) failed: %v", err)
	}
	if got := strings.TrimSpace(stdout.String()); got != resolveVersion() {
		t.Fatalf("version output = %q, want %q", got, resolveVersion())
	}
	if stderr.Len() != 0 {
		t.Fatalf("unexpected stderr: %q", stderr.String())
	}
}

func TestRunInvalidArguments(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{name: "zero width", args: []string{"-w", "0"}, wantErr: "width and height must be positive"},
		{name: "negative width", args: []string{"-width", "-1"}, wantErr: "width and height must be positive"},
		{name: "zero height", args: []string{"-h", "0"}, wantErr: "width and height must be positive"},
		{name: "negative height", args: []string{"-height", "-1"}, wantErr: "width and height must be positive"},
		{name: "zero count", args: []string{"-count", "0"}, wantErr: "count must be positive"},
		{name: "negative count", args: []string{"-count", "-1"}, wantErr: "count must be positive"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer

			err := run(tc.args, &stdout, &stderr)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("error = %q, want it to contain %q", err, tc.wantErr)
			}
			if stdout.Len() != 0 {
				t.Fatalf("unexpected stdout: %q", stdout.String())
			}
		})
	}
}

func TestRunSingleOutput(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "wallpaper")
	var stdout, stderr bytes.Buffer

	err := run([]string{"-w", "8", "-h", "8", "-seed", "42", "-out", output}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if _, err := os.Stat(output + ".png"); err != nil {
		t.Fatalf("expected single output file: %v", err)
	}
}

func TestRunSeedReproducibility(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()
	output1 := filepath.Join(dir1, "seeded")
	output2 := filepath.Join(dir2, "seeded")
	var stdout, stderr bytes.Buffer

	for _, output := range []string{output1, output2} {
		if err := run([]string{"-w", "16", "-h", "16", "-seed", "42", "-out", output}, &stdout, &stderr); err != nil {
			t.Fatalf("run failed: %v", err)
		}
	}

	first, err := os.ReadFile(output1 + ".png")
	if err != nil {
		t.Fatalf("cannot read first image: %v", err)
	}
	second, err := os.ReadFile(output2 + ".png")
	if err != nil {
		t.Fatalf("cannot read second image: %v", err)
	}
	if !bytes.Equal(first, second) {
		t.Fatal("the same seed produced different images")
	}
}

func TestRunRandomSeed(t *testing.T) {
	output := filepath.Join(t.TempDir(), "random")
	var stdout, stderr bytes.Buffer

	if err := run([]string{"-w", "8", "-h", "8", "-out", output}, &stdout, &stderr); err != nil {
		t.Fatalf("run without -seed failed: %v", err)
	}
	if _, err := os.Stat(output + ".png"); err != nil {
		t.Fatalf("expected output file: %v", err)
	}
}

func TestRunBatchOutput(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "batch")
	var stdout, stderr bytes.Buffer

	err := run([]string{"-w", "8", "-h", "8", "-seed", "42", "-count", "3", "-out", output}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}

	for i := 1; i <= 3; i++ {
		filename := filepath.Join(dir, "batch", fmt.Sprintf("batch_%d.png", i))
		if _, err := os.Stat(filename); err != nil {
			t.Fatalf("missing batch output %d: %v", i, err)
		}
	}
	if _, err := os.Stat(output + ".png"); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("batch mode should not create a single output file")
	}
}

func TestRunUnknownFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer

	err := run([]string{"-nope"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected an error for an unknown flag")
	}
	if !strings.Contains(stderr.String(), "flag provided but not defined") {
		t.Fatalf("stderr = %q, want flag error", stderr.String())
	}
}

func normalizeHelp(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i, line := range lines {
		lines[i] = strings.Join(strings.Fields(line), " ")
	}
	return strings.Join(lines, "\n")
}

func TestRunHelpMatchesReadme(t *testing.T) {
	var stdout, stderr bytes.Buffer

	if err := run([]string{"-help"}, &stdout, &stderr); !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("run(-help) error = %v, want flag.ErrHelp", err)
	}

	data, err := os.ReadFile(filepath.Join("..", "..", "README.md"))
	if err != nil {
		t.Fatalf("cannot read README.md: %v", err)
	}
	readme := strings.ReplaceAll(string(data), "\r\n", "\n")
	block := regexp.MustCompile("(?s)There is a `-help` flag[^\\n]*\n\n```\n(.*?)\n```").FindStringSubmatch(readme)
	if len(block) != 2 {
		t.Fatal("README help block not found")
	}

	got := normalizeHelp(stderr.String())
	want := normalizeHelp(block[1])
	if got != want {
		t.Fatalf("README help block differs from -help output\nREADME:\n%s\nCLI:\n%s", want, got)
	}
}

func TestRunHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer

	err := run([]string{"-help"}, &stdout, &stderr)
	if !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("run(-help) error = %v, want flag.ErrHelp", err)
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout: %q", stdout.String())
	}

	help := stderr.String()
	if !strings.Contains(help, "Usage: randwallpaper [options]") {
		t.Fatalf("help output missing usage line: %q", help)
	}
	for _, name := range []string{"-count", "-h", "-height", "-out", "-seed", "-version", "-w", "-width"} {
		if !strings.Contains(help, name) {
			t.Fatalf("help output missing %s: %q", name, help)
		}
	}
}
