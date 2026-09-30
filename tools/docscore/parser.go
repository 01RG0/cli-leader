package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Section represents a parsed Markdown section delimited by H2 or H3 headings.
type Section struct {
	File    string `json:"file"`
	Heading string `json:"heading"`
	Level   int    `json:"level"`
	Text    string `json:"text"`
	Line    int    `json:"line"`
}

// ID returns a unique identifier for the section.
func (s Section) ID() string {
	return fmt.Sprintf("%s#%s", s.File, s.Heading)
}

// ParseMarkdownReader parses markdown text from an io.Reader into sections split by H2/H3.
func ParseMarkdownReader(filePath string, r io.Reader) ([]Section, error) {
	scanner := bufio.NewScanner(r)
	// Allow scanning lines up to 1MB
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var sections []Section
	var currentHeading = "(Preamble)"
	var currentLevel = 1
	var currentLines []string
	var startLine = 1
	var lineNum = 0
	var inCodeBlock = false
	var codeFence = ""

	flushSection := func() {
		content := strings.TrimSpace(strings.Join(currentLines, "\n"))
		if content != "" {
			sections = append(sections, Section{
				File:    filePath,
				Heading: currentHeading,
				Level:   currentLevel,
				Text:    content,
				Line:    startLine,
			})
		}
		currentLines = nil
	}

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Check for code fence start/end
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			fence := trimmed[:3]
			if !inCodeBlock {
				inCodeBlock = true
				codeFence = fence
			} else if strings.HasPrefix(trimmed, codeFence) {
				inCodeBlock = false
				codeFence = ""
			}
			currentLines = append(currentLines, line)
			continue
		}

		if inCodeBlock {
			currentLines = append(currentLines, line)
			continue
		}

		// Detect H2 and H3 headings
		if strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "### ") {
			flushSection()
			if strings.HasPrefix(line, "## ") {
				currentHeading = strings.TrimSpace(strings.TrimPrefix(line, "## "))
				currentLevel = 2
			} else {
				currentHeading = strings.TrimSpace(strings.TrimPrefix(line, "### "))
				currentLevel = 3
			}
			startLine = lineNum
			continue
		}

		// If first section is preamble and we encounter an H1 title, use it as heading
		if currentHeading == "(Preamble)" && strings.HasPrefix(line, "# ") && len(currentLines) == 0 {
			currentHeading = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			currentLevel = 1
			startLine = lineNum
			continue
		}

		currentLines = append(currentLines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading markdown %s: %w", filePath, err)
	}

	flushSection()
	return sections, nil
}

// ParseMarkdownFile parses a file on disk into sections.
func ParseMarkdownFile(filePath string) ([]Section, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Normalize relative path
	cleanPath := filepath.Clean(filePath)
	return ParseMarkdownReader(cleanPath, f)
}

// FindMarkdownFiles searches for markdown files matching a glob or default files list.
func FindMarkdownFiles(pattern string) ([]string, error) {
	if pattern != "" {
		// If specific glob is provided
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		var clean []string
		for _, m := range matches {
			if strings.HasSuffix(m, ".md") {
				clean = append(clean, filepath.Clean(m))
			}
		}
		return clean, nil
	}

	// Default targets: README.md, CONTRIBUTING.md, docs/**
	var files []string
	candidates := []string{"README.md", "CONTRIBUTING.md"}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			files = append(files, c)
		}
	}

	err := filepath.Walk("docs", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip if docs folder doesn't exist yet
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			// Skip _review directory to prevent circular inspection
			if !strings.Contains(path, "_review") {
				files = append(files, filepath.Clean(path))
			}
		}
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return files, nil
}
