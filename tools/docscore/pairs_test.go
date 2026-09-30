package main

import (
	"strings"
	"testing"
)

func TestExtractKeywordsAndJaccard(t *testing.T) {
	tests := []struct {
		name     string
		textA    string
		textB    string
		minSim   float64
		maxSim   float64
	}{
		{
			name:   "Identical content",
			textA:  "Reverse proxy gateway with streaming SSE",
			textB:  "Reverse proxy gateway with streaming SSE",
			minSim: 0.99,
			maxSim: 1.01,
		},
		{
			name:   "Completely different content",
			textA:  "Authentication token secret keys",
			textB:  "Database migration schema indexes",
			minSim: 0.0,
			maxSim: 0.0,
		},
		{
			name:   "Partial overlap with stopwords",
			textA:  "The architecture of the gateway",
			textB:  "Design of the gateway server",
			minSim: 0.20,
			maxSim: 0.60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setA := ExtractKeywords(tt.textA)
			setB := ExtractKeywords(tt.textB)
			sim := JaccardSimilarity(setA, setB)

			if sim < tt.minSim || sim > tt.maxSim {
				t.Errorf("expected similarity in range [%f, %f], got %f", tt.minSim, tt.maxSim, sim)
			}
		})
	}
}

func TestFindCandidatePairs(t *testing.T) {
	sections := []Section{
		{
			File:    "README.md",
			Heading: "Architecture",
			Text:    "The architecture consists of a gateway and worker CLIs.",
		},
		{
			File:    "README.md",
			Heading: "Installation",
			Text:    "Run go build to install cli-leader.",
		},
		{
			File:    "docs/ARCHITECTURE.md",
			Heading: "Architecture Deep Dive",
			Text:    "Detailed overview of the gateway and worker CLIs architecture.",
		},
		{
			File:    "docs/ROADMAP.md",
			Heading: "Milestones",
			Text:    "Phases and schedule for upcoming releases.",
		},
	}

	pairs := FindCandidatePairs(sections, 10)

	// Should not pair README Architecture with README Installation (same file)
	for _, p := range pairs {
		if p.SectionA.File == p.SectionB.File {
			t.Errorf("found pair from same file: %s and %s", p.SectionA.File, p.SectionB.File)
		}
	}

	// Should prioritize README Architecture <=> docs/ARCHITECTURE.md Deep Dive
	if len(pairs) == 0 {
		t.Fatalf("expected candidate pairs, got 0")
	}

	topPair := pairs[0]
	if !((topPair.SectionA.Heading == "Architecture" && topPair.SectionB.Heading == "Architecture Deep Dive") ||
		(topPair.SectionB.Heading == "Architecture" && topPair.SectionA.Heading == "Architecture Deep Dive")) {
		t.Errorf("expected top pair to be Architecture <=> Architecture Deep Dive, got %s <=> %s",
			topPair.SectionA.Heading, topPair.SectionB.Heading)
	}

	// Test FormatState output
	formatted := topPair.FormatState()
	if !strings.Contains(formatted, "--- Section A (File:") || !strings.Contains(formatted, "--- Section B (File:") {
		t.Errorf("FormatState does not contain proper section labels:\n%s", formatted)
	}

	// Test capping limit
	capped := FindCandidatePairs(sections, 1)
	if len(capped) != 1 {
		t.Errorf("expected exactly 1 pair when capped, got %d", len(capped))
	}
}
