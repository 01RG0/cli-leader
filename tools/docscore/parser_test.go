package main

import (
	"strings"
	"testing"
)

func TestParseMarkdownReader(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []struct {
			heading string
			level   int
			hasText string
		}
	}{
		{
			name: "Basic H2 and H3 split",
			input: `# Document Title
Preamble line.

## Overview
This is the overview.

### Architecture Details
Here are the architecture details.
Line 2 of architecture.

## Roadmap
Upcoming milestones.
`,
			expected: []struct {
				heading string
				level   int
				hasText string
			}{
				{"Document Title", 1, "Preamble line."},
				{"Overview", 2, "This is the overview."},
				{"Architecture Details", 3, "Here are the architecture details.\nLine 2 of architecture."},
				{"Roadmap", 2, "Upcoming milestones."},
			},
		},
		{
			name: "Ignore headings inside code blocks",
			input: `## Code Examples
Here is some bash:
` + "```bash" + `
## This is a bash comment, not an H2!
### Neither is this!
echo "hello"
` + "```" + `
End of code section.

## Next Section
Real heading here.
`,
			expected: []struct {
				heading string
				level   int
				hasText string
			}{
				{"Code Examples", 2, "## This is a bash comment, not an H2!"},
				{"Next Section", 2, "Real heading here."},
			},
		},
		{
			name: "Handle H4 within H3 section",
			input: `## Main Feature
Feature intro.

### Sub-feature A
Sub-feature intro.

#### Sub-sub-detail
This H4 detail should stay inside Sub-feature A section.
`,
			expected: []struct {
				heading string
				level   int
				hasText string
			}{
				{"Main Feature", 2, "Feature intro."},
				{"Sub-feature A", 3, "This H4 detail should stay inside Sub-feature A section."},
			},
		},
		{
			name: "Empty sections skipped",
			input: `## First Section
## Second Section (immediately follows)
Actual text in second section.
`,
			expected: []struct {
				heading string
				level   int
				hasText string
			}{
				{"Second Section (immediately follows)", 2, "Actual text in second section."},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sections, err := ParseMarkdownReader("test.md", strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(sections) != len(tt.expected) {
				t.Fatalf("expected %d sections, got %d", len(tt.expected), len(sections))
			}

			for i, exp := range tt.expected {
				act := sections[i]
				if act.Heading != exp.heading {
					t.Errorf("[%d] expected heading %q, got %q", i, exp.heading, act.Heading)
				}
				if act.Level != exp.level {
					t.Errorf("[%d] expected level %d, got %d", i, exp.level, act.Level)
				}
				if !strings.Contains(act.Text, exp.hasText) {
					t.Errorf("[%d] expected text to contain %q, but got %q", i, exp.hasText, act.Text)
				}
			}
		})
	}
}
