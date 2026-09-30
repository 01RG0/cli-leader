package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

func main() {
	filesGlob := flag.String("files", "", "Glob pattern for markdown files (default: README.md, CONTRIBUTING.md, docs/**/*.md)")
	concurrency := flag.Int("concurrency", 4, "Number of concurrent API requests")
	maxPairs := flag.Int("max-pairs", 300, "Maximum number of cross-file candidate pairs to evaluate")
	dryRun := flag.Bool("dry-run", false, "Simulate execution and print planned calls and token estimates without calling API")
	jsonOnly := flag.Bool("json-only", false, "Output only final JSON to stdout")
	outputDir := flag.String("output-dir", filepath.Join("docs", "_review"), "Directory for generated reports")
	cacheDir := flag.String("cache-dir", filepath.Join(".cache", "docscore"), "Directory for request caching")
	calibrateFile := flag.String("calibrate", "", "Path to labels.csv for calibration evaluation")
	flag.Parse()

	// 1. Security & Environment setup
	LoadEnv(".env")
	apiKey := os.Getenv("DREX_API_KEY")

	// 2. Discover markdown files
	files, err := FindMarkdownFiles(*filesGlob)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering markdown files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "No markdown files found matching pattern.\n")
		os.Exit(1)
	}

	// 3. Parse files into sections
	var allSections []Section
	for _, f := range files {
		sections, err := ParseMarkdownFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing file %s: %v\n", f, err)
			continue
		}
		allSections = append(allSections, sections...)
	}

	// 4. Find candidate cross-file pairs
	candidatePairs := FindCandidatePairs(allSections, *maxPairs)

	// 5. Calculate token estimates and planned calls
	totalSectionCalls := len(allSections)
	totalPairCalls := len(candidatePairs)
	totalPlannedCalls := totalSectionCalls + totalPairCalls

	var estimatedInputTokens int
	for _, s := range allSections {
		// ~4 characters per token + ~80 prompt/question overhead tokens
		estimatedInputTokens += (len(s.Text) / 4) + 80
	}
	for _, p := range candidatePairs {
		estimatedInputTokens += ((len(p.SectionA.Text) + len(p.SectionB.Text)) / 4) + 50
	}

	if !*jsonOnly {
		fmt.Println("================================================================================")
		fmt.Println("🔍 Drex Documentation Quality Scorer (tools/docscore)")
		fmt.Println("================================================================================")
		fmt.Printf("Files Discovered:         %d\n", len(files))
		fmt.Printf("Sections Parsed:          %d (H2/H3 blocks)\n", len(allSections))
		fmt.Printf("Candidate Pairs Selected: %d (Capped at %d)\n", len(candidatePairs), *maxPairs)
		fmt.Printf("Total Planned Calls:      %d (Sections: %d, Pairs: %d)\n", totalPlannedCalls, totalSectionCalls, totalPairCalls)
		fmt.Printf("Estimated Input Tokens:   ~%d tokens\n", estimatedInputTokens)
		fmt.Println("================================================================================")
	}

	// If dry-run mode, stop here and display estimates
	if *dryRun {
		if !*jsonOnly {
			fmt.Println("\n🚀 DRY RUN MODE ACTIVE: No external API requests were dispatched.")
			fmt.Printf("Planned calls: %d | Estimated tokens: ~%d\n\n", totalPlannedCalls, estimatedInputTokens)
			printOpenQuestions()
		}
		return
	}

	if apiKey == "" {
		fmt.Fprintf(os.Stderr, "Error: DREX_API_KEY is not set in environment or .env file.\n")
		os.Exit(1)
	}

	// Initialize API Client
	client, err := NewClient(apiKey, *cacheDir, *concurrency)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing client: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	var (
		findingsMu      sync.Mutex
		sectionFindings []Finding
		pairFindings    []Finding
		sectionRoles    = make(map[string]map[string]string)
		totalTokens     UsageInfo
		totalRequests   int
		cacheHits       int
	)

	// Progress tracker
	progressCount := 0
	totalToRun := totalPlannedCalls

	// 6. Execute Section Questions
	if !*jsonOnly {
		fmt.Printf("Evaluating %d sections (concurrency: %d)...\n", len(allSections), *concurrency)
	}

	var wg sync.WaitGroup
	secQuestions := SectionQuestions()

	for _, sec := range allSections {
		wg.Add(1)
		go func(s Section) {
			defer wg.Done()

			req := SystemOneRequest{
				Model:     DefaultModel,
				State:     s.Text,
				Questions: secQuestions,
			}

			resp, hit, evalErr := client.Evaluate(ctx, req)
			findingsMu.Lock()
			progressCount++
			if !*jsonOnly && (progressCount%10 == 0 || progressCount == totalToRun) {
				fmt.Printf("  [%d/%d] Completed calls (cache hits: %d)...\r", progressCount, totalToRun, cacheHits)
			}

			if evalErr != nil {
				fmt.Fprintf(os.Stderr, "\nWarning: Section evaluation failed (%s#%s): %s\n", s.File, s.Heading, client.RedactKey(evalErr.Error()))
				findingsMu.Unlock()
				return
			}

			totalRequests++
			if hit {
				cacheHits++
			}
			totalTokens.InputTokens += resp.Usage.InputTokens
			totalTokens.OutputTokens += resp.Usage.OutputTokens

			snippet := TruncateSnippet(s.Text, 200)

			// Process Noul Answers
			for qKey, answer := range resp.Answers {
				if answer.Type == "noul" {
					sectionFindings = append(sectionFindings, Finding{
						File:        s.File,
						Heading:     s.Heading,
						Question:    qKey,
						Probability: answer.Noul,
						Snippet:     snippet,
						RequestID:   resp.RequestID,
					})
				} else if answer.Type == "choice" {
					if _, ok := sectionRoles[s.File]; !ok {
						sectionRoles[s.File] = make(map[string]string)
					}
					sectionRoles[s.File][s.Heading] = answer.Choice
				}
			}
			findingsMu.Unlock()
		}(sec)
	}
	wg.Wait()

	// 7. Execute Pair Questions
	if !*jsonOnly {
		fmt.Printf("\nEvaluating %d cross-file candidate pairs (concurrency: %d)...\n", len(candidatePairs), *concurrency)
	}

	pairQuestions := PairQuestions()
	for _, p := range candidatePairs {
		wg.Add(1)
		go func(pair SectionPair) {
			defer wg.Done()

			req := SystemOneRequest{
				Model:     DefaultModel,
				State:     pair.FormatState(),
				Questions: pairQuestions,
			}

			resp, hit, evalErr := client.Evaluate(ctx, req)
			findingsMu.Lock()
			progressCount++
			if !*jsonOnly && (progressCount%10 == 0 || progressCount == totalToRun) {
				fmt.Printf("  [%d/%d] Completed calls (cache hits: %d)...\r", progressCount, totalToRun, cacheHits)
			}

			if evalErr != nil {
				fmt.Fprintf(os.Stderr, "\nWarning: Pair evaluation failed (%s <=> %s): %s\n", pair.SectionA.ID(), pair.SectionB.ID(), client.RedactKey(evalErr.Error()))
				findingsMu.Unlock()
				return
			}

			totalRequests++
			if hit {
				cacheHits++
			}
			totalTokens.InputTokens += resp.Usage.InputTokens
			totalTokens.OutputTokens += resp.Usage.OutputTokens

			pairSnippet := TruncateSnippet(fmt.Sprintf("[A: %s] %s | [B: %s] %s", pair.SectionA.Heading, pair.SectionA.Text, pair.SectionB.Heading, pair.SectionB.Text), 200)

			for qKey, answer := range resp.Answers {
				if answer.Type == "noul" && answer.Noul >= 0.20 {
					pairFindings = append(pairFindings, Finding{
						File:        pair.SectionA.File,
						Heading:     pair.SectionA.Heading,
						PairFile:    pair.SectionB.File,
						PairHeading: pair.SectionB.Heading,
						Question:    qKey,
						Probability: answer.Noul,
						Snippet:     pairSnippet,
						RequestID:   resp.RequestID,
					})
				}
			}
			findingsMu.Unlock()
		}(p)
	}
	wg.Wait()

	if !*jsonOnly {
		fmt.Printf("\nDone evaluating. Sorting findings...\n")
	}

	// Sort section findings by probability descending
	sort.Slice(sectionFindings, func(i, j int) bool {
		return sectionFindings[i].Probability > sectionFindings[j].Probability
	})

	// Sort pair findings by probability descending
	sort.Slice(pairFindings, func(i, j int) bool {
		return pairFindings[i].Probability > pairFindings[j].Probability
	})

	// Build summary statistics
	summaries := BuildSummaries(allSections, sectionFindings)

	report := FullReport{
		GeneratedAt:     time.Now().UTC().Format(time.RFC3339),
		Model:           DefaultModel,
		TotalSections:   len(allSections),
		TotalPairs:      len(candidatePairs),
		TotalRequests:   totalRequests,
		CacheHits:       cacheHits,
		TotalTokens:     totalTokens,
		FileSummaries:   summaries,
		SectionFindings: sectionFindings,
		PairFindings:    pairFindings,
		SectionRoles:    sectionRoles,
	}

	// 8. Save Reports
	if err := SaveReports(*outputDir, report); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving reports: %v\n", err)
	} else if !*jsonOnly {
		fmt.Printf("\n✅ Reports successfully generated at:\n - %s\n - %s\n",
			filepath.Join(*outputDir, "DREX_REPORT.json"),
			filepath.Join(*outputDir, "DREX_REPORT.md"))
	}

	// 9. Calibration if requested
	if *calibrateFile != "" {
		labels, err := LoadCalibrationLabels(*calibrateFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Calibration Error: %v\n", err)
		} else {
			allFindings := append([]Finding{}, sectionFindings...)
			allFindings = append(allFindings, pairFindings...)
			EvaluateCalibration(os.Stdout, labels, allFindings)
		}
	}

	if !*jsonOnly {
		fmt.Printf("\nTotal calls planned: %d | Total calls executed: %d (Cache hits: %d)\n", totalPlannedCalls, totalRequests, cacheHits)
		fmt.Printf("Estimated tokens used: %d input / %d output\n\n", totalTokens.InputTokens, totalTokens.OutputTokens)
		printOpenQuestions()
	}
}

func printOpenQuestions() {
	fmt.Println("================================================================================")
	fmt.Println("❓ OPEN QUESTIONS")
	fmt.Println("================================================================================")
	fmt.Println("1. [Question Schema]: The user request specified a 'multinomial' type for section_role.")
	fmt.Println("   However, Drex API documentation (OpenAPI 3.2.0 at https://drex.nace.ai/llms-full.txt)")
	fmt.Println("   strictly defines supported types as 'noul', 'choice', and 'score'.")
	fmt.Println("   Any request with type 'multinomial' is rejected with HTTP 422 ('must be \"noul\", \"choice\", or \"score\"').")
	fmt.Println("   Resolution applied: 'section_role' is implemented using the official 'choice' type with discrete criteria labels.")
	fmt.Println("2. [Pair Threshold]: What probability cutoff should constitute a hard block for CI/CD?")
	fmt.Println("   Drex documentation emphasizes that thresholds are business decisions. Using --calibrate with")
	fmt.Println("   labels.csv allows measuring empirical precision and recall at 0.5, 0.7, and 0.9.")
	fmt.Println("================================================================================")
}
