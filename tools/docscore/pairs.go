package main

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// SectionPair represents a candidate pair of sections from different files.
type SectionPair struct {
	SectionA Section `json:"section_a"`
	SectionB Section `json:"section_b"`
	Score    float64 `json:"score"`
}

// PairID returns a stable identifier for the pair.
func (p SectionPair) PairID() string {
	return fmt.Sprintf("%s <=> %s", p.SectionA.ID(), p.SectionB.ID())
}

// FormatState formats both sections clearly labeled as Section A and Section B.
func (p SectionPair) FormatState() string {
	return fmt.Sprintf("--- Section A (File: %s, Heading: %s) ---\n%s\n\n--- Section B (File: %s, Heading: %s) ---\n%s",
		p.SectionA.File, p.SectionA.Heading, p.SectionA.Text,
		p.SectionB.File, p.SectionB.Heading, p.SectionB.Text)
}

// Common stopwords to exclude from keyword extraction
var commonStopwords = map[string]bool{
	"a": true, "about": true, "above": true, "after": true, "again": true, "against": true,
	"all": true, "am": true, "an": true, "and": true, "any": true, "are": true, "as": true,
	"at": true, "be": true, "because": true, "been": true, "before": true, "being": true,
	"below": true, "between": true, "both": true, "but": true, "by": true, "can": true,
	"did": true, "do": true, "does": true, "doing": true, "down": true, "during": true,
	"each": true, "few": true, "for": true, "from": true, "further": true, "had": true,
	"has": true, "have": true, "having": true, "he": true, "her": true, "here": true,
	"hers": true, "herself": true, "him": true, "himself": true, "his": true, "how": true,
	"if": true, "in": true, "into": true, "is": true, "it": true, "its": true, "itself": true,
	"just": true, "me": true, "more": true, "most": true, "my": true, "myself": true,
	"no": true, "nor": true, "not": true, "now": true, "of": true, "off": true, "on": true,
	"once": true, "only": true, "or": true, "other": true, "our": true, "ours": true,
	"ourselves": true, "out": true, "over": true, "own": true, "same": true, "she": true,
	"should": true, "so": true, "some": true, "such": true, "than": true, "that": true,
	"the": true, "their": true, "theirs": true, "them": true, "themselves": true, "then": true,
	"there": true, "these": true, "they": true, "this": true, "those": true, "through": true,
	"to": true, "too": true, "under": true, "until": true, "up": true, "very": true, "was": true,
	"we": true, "were": true, "what": true, "when": true, "where": true, "which": true,
	"while": true, "who": true, "whom": true, "why": true, "will": true, "with": true,
	"you": true, "your": true, "yours": true, "yourself": true, "yourselves": true,
}

// ExtractKeywords extracts normalized, non-stopword tokens from text.
func ExtractKeywords(text string) map[string]bool {
	tokens := make(map[string]bool)
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	words := strings.FieldsFunc(strings.ToLower(text), f)
	for _, w := range words {
		if len(w) >= 3 && !commonStopwords[w] {
			tokens[w] = true
		}
	}
	return tokens
}

// JaccardSimilarity computes Jaccard index between two token sets.
func JaccardSimilarity(setA, setB map[string]bool) float64 {
	if len(setA) == 0 || len(setB) == 0 {
		return 0.0
	}
	intersection := 0
	for k := range setA {
		if setB[k] {
			intersection++
		}
	}
	if intersection == 0 {
		return 0.0
	}
	union := len(setA)
	for k := range setB {
		if !setA[k] {
			union++
		}
	}
	return float64(intersection) / float64(union)
}

// ComputePairScore computes a relevance score between two cross-file sections.
func ComputePairScore(a, b Section) float64 {
	// Headings
	hA := ExtractKeywords(a.Heading)
	hB := ExtractKeywords(b.Heading)
	hSim := JaccardSimilarity(hA, hB)

	// Body text
	tA := ExtractKeywords(a.Text)
	tB := ExtractKeywords(b.Text)
	tSim := JaccardSimilarity(tA, tB)

	// Heading match is weighted 4x higher than general text match
	return (4.0 * hSim) + (1.0 * tSim)
}

// FindCandidatePairs selects and ranks candidate cross-file pairs using the similarity heuristic.
func FindCandidatePairs(sections []Section, maxPairs int) []SectionPair {
	if maxPairs <= 0 {
		maxPairs = 300
	}

	var candidates []SectionPair

	for i := 0; i < len(sections); i++ {
		for j := i + 1; j < len(sections); j++ {
			a := sections[i]
			b := sections[j]

			// Must be cross-file
			if a.File == b.File {
				continue
			}

			score := ComputePairScore(a, b)
			// Only include pairs with some lexical overlap
			if score > 0.02 {
				candidates = append(candidates, SectionPair{
					SectionA: a,
					SectionB: b,
					Score:    score,
				})
			}
		}
	}

	// Sort descending by score
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})

	if len(candidates) > maxPairs {
		candidates = candidates[:maxPairs]
	}

	return candidates
}
