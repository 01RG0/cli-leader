package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Finding represents an individual scored issue or question result.
type Finding struct {
	File        string  `json:"file"`
	Heading     string  `json:"heading"`
	Question    string  `json:"question"`
	Probability float64 `json:"probability"`
	Snippet     string  `json:"snippet"`
	RequestID   string  `json:"request_id"`
	PairFile    string  `json:"pair_file,omitempty"`
	PairHeading string  `json:"pair_heading,omitempty"`
}

// FileSummary summarizes scores per document.
type FileSummary struct {
	File              string  `json:"file"`
	SectionCount      int     `json:"section_count"`
	HighRiskFindings  int     `json:"high_risk_findings"` // Count with P >= 0.7
	TotalFlaggedCount int     `json:"total_flagged_count"` // Count with P >= 0.5
	MaxProbability    float64 `json:"max_probability"`
	AvgProbability    float64 `json:"avg_probability"`
}

// FullReport contains all findings and metadata.
type FullReport struct {
	GeneratedAt     string                       `json:"generated_at"`
	Model           string                       `json:"model"`
	TotalSections   int                          `json:"total_sections"`
	TotalPairs      int                          `json:"total_pairs"`
	TotalRequests   int                          `json:"total_requests"`
	CacheHits       int                          `json:"cache_hits"`
	TotalTokens     UsageInfo                    `json:"total_tokens"`
	FileSummaries   []FileSummary                `json:"file_summaries"`
	SectionFindings []Finding                    `json:"section_findings"`
	PairFindings    []Finding                    `json:"pair_findings"`
	SectionRoles    map[string]map[string]string `json:"section_roles"`
}

// TruncateSnippet produces a clean single-line snippet up to 200 runes.
func TruncateSnippet(text string, maxRunes int) string {
	cleaned := strings.Join(strings.Fields(text), " ")
	runes := []rune(cleaned)
	if len(runes) <= maxRunes {
		return string(runes)
	}
	return string(runes[:maxRunes-3]) + "..."
}

// GenerateMarkdownReport converts a FullReport into GitHub-flavored Markdown.
func GenerateMarkdownReport(report FullReport) string {
	var sb strings.Builder

	sb.WriteString("# Drex Documentation Quality & Verification Report\n\n")
	sb.WriteString(fmt.Sprintf("> Generated: `%s` | Model: `%s` | Total Requests: `%d` (Cache hits: `%d`)\n",
		report.GeneratedAt, report.Model, report.TotalRequests, report.CacheHits))
	sb.WriteString(fmt.Sprintf("> Total Input Tokens: `%d` | Total Output Tokens: `%d`\n\n",
		report.TotalTokens.InputTokens, report.TotalTokens.OutputTokens))

	// Summary Table Per File
	sb.WriteString("## 1. Summary by File\n\n")
	sb.WriteString("| File | Sections | Flagged (P ≥ 0.5) | High Risk (P ≥ 0.7) | Max Prob | Avg Prob |\n")
	sb.WriteString("| :--- | :---: | :---: | :---: | :---: | :---: |\n")
	for _, fs := range report.FileSummaries {
		sb.WriteString(fmt.Sprintf("| `%s` | %d | %d | %d | `%.3f` | `%.3f` |\n",
			fs.File, fs.SectionCount, fs.TotalFlaggedCount, fs.HighRiskFindings, fs.MaxProbability, fs.AvgProbability))
	}
	sb.WriteString("\n---\n\n")

	// Section Findings Table
	sb.WriteString("## 2. Section Findings (Sorted by Probability)\n\n")
	if len(report.SectionFindings) == 0 {
		sb.WriteString("No section findings detected.\n\n")
	} else {
		sb.WriteString("| Probability | Question | File & Heading | Snippet (First 200 chars) | Request ID |\n")
		sb.WriteString("| :---: | :--- | :--- | :--- | :--- |\n")
		for _, f := range report.SectionFindings {
			sb.WriteString(fmt.Sprintf("| **`%.3f`** | `%s` | **`%s`**<br>`%s` | %s | `%s` |\n",
				f.Probability, f.Question, f.File, f.Heading, f.Snippet, f.RequestID))
		}
		sb.WriteString("\n")
	}
	sb.WriteString("---\n\n")

	// Cross-File Pair Findings Table
	sb.WriteString("## 3. Cross-File Contradictions & Duplications\n\n")
	if len(report.PairFindings) == 0 {
		sb.WriteString("No cross-file contradiction or duplication pairs scored P ≥ 0.30.\n\n")
	} else {
		sb.WriteString("| Probability | Question | Section A | Section B | Request ID |\n")
		sb.WriteString("| :---: | :--- | :--- | :--- | :--- |\n")
		for _, f := range report.PairFindings {
			sb.WriteString(fmt.Sprintf("| **`%.3f`** | `%s` | **`%s`**<br>`%s` | **`%s`**<br>`%s` | `%s` |\n",
				f.Probability, f.Question, f.File, f.Heading, f.PairFile, f.PairHeading, f.RequestID))
		}
		sb.WriteString("\n")
	}
	sb.WriteString("---\n\n")

	// Section Role Classifications
	sb.WriteString("## 4. Section Roles Distribution\n\n")
	sb.WriteString("| File | Heading | Role |\n")
	sb.WriteString("| :--- | :--- | :--- |\n")
	for file, headings := range report.SectionRoles {
		for heading, role := range headings {
			sb.WriteString(fmt.Sprintf("| `%s` | %s | `%s` |\n", file, heading, role))
		}
	}
	sb.WriteString("\n")

	return sb.String()
}

// SaveReports writes both JSON and Markdown reports to disk.
func SaveReports(outputDir string, report FullReport) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 1. Write JSON report
	jsonPath := filepath.Join(outputDir, "DREX_REPORT.json")
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON report: %w", err)
	}
	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", jsonPath, err)
	}

	// 2. Write Markdown report
	mdPath := filepath.Join(outputDir, "DREX_REPORT.md")
	mdData := GenerateMarkdownReport(report)
	if err := os.WriteFile(mdPath, []byte(mdData), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", mdPath, err)
	}

	return nil
}

// BuildSummaries generates summary statistics per file from findings and sections.
func BuildSummaries(sections []Section, findings []Finding) []FileSummary {
	fileSecCounts := make(map[string]int)
	for _, s := range sections {
		fileSecCounts[s.File]++
	}

	fileProbMap := make(map[string][]float64)
	fileFlagged := make(map[string]int)
	fileHighRisk := make(map[string]int)

	for _, f := range findings {
		fileProbMap[f.File] = append(fileProbMap[f.File], f.Probability)
		if f.Probability >= 0.5 {
			fileFlagged[f.File]++
		}
		if f.Probability >= 0.7 {
			fileHighRisk[f.File]++
		}
	}

	var summaries []FileSummary
	for file, count := range fileSecCounts {
		probs := fileProbMap[file]
		maxProb := 0.0
		sumProb := 0.0
		for _, p := range probs {
			if p > maxProb {
				maxProb = p
			}
			sumProb += p
		}
		avgProb := 0.0
		if len(probs) > 0 {
			avgProb = sumProb / float64(len(probs))
		}

		summaries = append(summaries, FileSummary{
			File:              file,
			SectionCount:      count,
			HighRiskFindings:  fileHighRisk[file],
			TotalFlaggedCount: fileFlagged[file],
			MaxProbability:    maxProb,
			AvgProbability:    avgProb,
		})
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].File < summaries[j].File
	})

	return summaries
}
