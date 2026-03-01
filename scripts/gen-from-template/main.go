package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	OutputDir string   `yaml:"output_dir"`
	Outputs   []Output `yaml:"outputs"`
}

type Output struct {
	MD          string   `yaml:"md"`
	Description string   `yaml:"description"`
	Preamble    string   `yaml:"preamble"`
	Sources     []string `yaml:"sources"`
	TestSources []string `yaml:"test_sources"`
	Structure   bool     `yaml:"structure"`
}

func main() {
	rootDir := findProjectRoot()

	configPath := filepath.Join(rootDir, "scripts", "gen-skills.config.yaml")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		fatal("reading config: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(configData, &cfg); err != nil {
		fatal("parsing config: %v", err)
	}

	outDir := filepath.Join(rootDir, cfg.OutputDir)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fatal("creating output dir: %v", err)
	}

	for _, out := range cfg.Outputs {
		content, err := buildMarkdown(rootDir, out)
		if err != nil {
			fatal("building %s: %v", out.MD, err)
		}

		dest := filepath.Join(outDir, out.MD)
		if err := os.WriteFile(dest, []byte(content), 0o644); err != nil {
			fatal("writing %s: %v", dest, err)
		}
		fmt.Printf("  generated: %s\n", filepath.Join(cfg.OutputDir, out.MD))
	}

	fmt.Printf("\nDone. %d files generated in %s/\n", len(cfg.Outputs), cfg.OutputDir)
}

func buildMarkdown(rootDir string, out Output) (string, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n\n", strings.TrimSuffix(out.MD, ".md")))
	sb.WriteString(fmt.Sprintf("> %s\n\n", out.Description))
	if out.Preamble != "" {
		sb.WriteString(strings.TrimSpace(out.Preamble) + "\n\n")
	}

	if out.Structure {
		tree := buildDirTree(rootDir)
		sb.WriteString("```\n")
		sb.WriteString(tree)
		sb.WriteString("\n```\n")
		return sb.String(), nil
	}

	for i, src := range out.Sources {
		srcPath := filepath.Join(rootDir, src)
		data, err := os.ReadFile(srcPath)
		if err != nil {
			return "", fmt.Errorf("reading %s: %w", src, err)
		}

		if i > 0 {
			sb.WriteString("\n---\n\n")
		}

		sb.WriteString(fmt.Sprintf("## %s\n\n", src))
		lang := codeBlockLang(src)
		sb.WriteString(fmt.Sprintf("```%s\n", lang))
		sb.WriteString(strings.TrimRight(string(data), "\n"))
		sb.WriteString("\n```\n")
	}

	if len(out.TestSources) > 0 {
		sb.WriteString("\n---\n\n## Unit tests\n\n")
		sb.WriteString("When creating a new component, generate tests following this pattern:\n\n")
		for i, src := range out.TestSources {
			srcPath := filepath.Join(rootDir, src)
			data, err := os.ReadFile(srcPath)
			if err != nil {
				return "", fmt.Errorf("reading test %s: %w", src, err)
			}
			if i > 0 {
				sb.WriteString("\n---\n\n")
			}
			sb.WriteString(fmt.Sprintf("### %s\n\n", src))
			lang := codeBlockLang(src)
			sb.WriteString(fmt.Sprintf("```%s\n", lang))
			sb.WriteString(strings.TrimRight(string(data), "\n"))
			sb.WriteString("\n```\n")
		}
	}

	return sb.String(), nil
}

func codeBlockLang(path string) string {
	switch filepath.Ext(path) {
	case ".sql":
		return "sql"
	case ".yaml", ".yml":
		return "yaml"
	case ".json":
		return "json"
	default:
		return "go"
	}
}

func buildDirTree(rootDir string) string {
	exclude := map[string]bool{
		".git": true, "node_modules": true, ".generated": true,
		"gen-from-template": true, "skills": true, ".einar": true,
	}
	var sb strings.Builder
	sb.WriteString(".\n")
	walkTree(&sb, rootDir, "", true, exclude)
	return strings.TrimRight(sb.String(), "\n")
}

func walkTree(sb *strings.Builder, dir, prefix string, isLast bool, exclude map[string]bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	var valid []os.DirEntry
	for _, e := range entries {
		if exclude[e.Name()] {
			continue
		}
		valid = append(valid, e)
	}

	for i, e := range valid {
		last := i == len(valid)-1
		connector := "├── "
		if last {
			connector = "└── "
		}
		extend := "│   "
		if last {
			extend = "    "
		}

		sb.WriteString(prefix + connector + e.Name() + "\n")

		if e.IsDir() {
			walkTree(sb, filepath.Join(dir, e.Name()), prefix+extend, last, exclude)
		}
	}
}

func findProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		fatal("getting cwd: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			fatal("could not find project root (no go.mod found)")
		}
		dir = parent
	}
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "gen-from-template: "+format+"\n", args...)
	os.Exit(1)
}
