package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// LabelRecord represents a ground-truth calibration entry.
type LabelRecord struct {
	File     string
	Heading  string
	Question string
	Label    int // 0 or 1
}

// CalibrationMetric tracks performance metrics at a specific threshold.
type CalibrationMetric struct {
	Threshold float64
	TP        int
	FP        int
	TN        int
	FN        int
	Total     int
	Agreement float64 // (TP + TN) / Total
	Precision float64 // TP / (TP + FP)
	Recall    float64 // TP / (TP + FN)
}

// LoadCalibrationLabels reads a CSV file containing ground-truth labels.
func LoadCalibrationLabels(csvPath string) ([]LabelRecord, error) {
	f, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read labels CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, errorsNew("labels CSV is empty or missing headers")
	}

	header := records[0]
	colIdx := make(map[string]int)
	for i, col := range header {
		colIdx[strings.ToLower(strings.TrimSpace(col))] = i
	}

	reqCols := []string{"file", "heading", "question", "label"}
	for _, col := range reqCols {
		if _, exists := colIdx[col]; !exists {
			return nil, fmt.Errorf("missing required column %q in labels CSV", col)
		}
	}

	var labels []LabelRecord
	for lineIdx, row := range records[1:] {
		if len(row) <= colIdx["label"] {
			continue
		}
		lblStr := strings.TrimSpace(row[colIdx["label"]])
		lblVal, err := strconv.Atoi(lblStr)
		if err != nil || (lblVal != 0 && lblVal != 1) {
			return nil, fmt.Errorf("invalid label %q on line %d (must be 0 or 1)", lblStr, lineIdx+2)
		}

		labels = append(labels, LabelRecord{
			File:     strings.TrimSpace(row[colIdx["file"]]),
			Heading:  strings.TrimSpace(row[colIdx["heading"]]),
			Question: strings.TrimSpace(row[colIdx["question"]]),
			Label:    lblVal,
		})
	}

	return labels, nil
}

func errorsNew(msg string) error {
	return fmt.Errorf("%s", msg)
}

// EvaluateCalibration evaluates agreement between model predictions and ground-truth labels.
func EvaluateCalibration(w io.Writer, labels []LabelRecord, findings []Finding) {
	// Index findings by Key = File#Heading#Question
	findMap := make(map[string]float64)
	for _, f := range findings {
		key := fmt.Sprintf("%s#%s#%s", strings.TrimSpace(f.File), strings.TrimSpace(f.Heading), strings.TrimSpace(f.Question))
		findMap[key] = f.Probability
	}

	// Group labels by Question
	questionLabels := make(map[string][]LabelRecord)
	for _, l := range labels {
		questionLabels[l.Question] = append(questionLabels[l.Question], l)
	}

	thresholds := []float64{0.5, 0.7, 0.9}

	fmt.Fprintln(w, "\n================================================================================")
	fmt.Fprintln(w, "🎯 CALIBRATION AGREEMENT REPORT (Ground Truth vs. Drex Probabilities)")
	fmt.Fprintln(w, "================================================================================")

	for q, qLabels := range questionLabels {
		fmt.Fprintf(w, "\nQuestion: [%s] (Total Labeled Samples: %d)\n", q, len(qLabels))
		fmt.Fprintf(w, "%-11s | %-5s | %-5s | %-5s | %-5s | %-10s | %-10s | %-10s\n",
			"Threshold", "TP", "FP", "TN", "FN", "Agreement", "Precision", "Recall")
		fmt.Fprintln(w, "--------------------------------------------------------------------------------")

		for _, thresh := range thresholds {
			var tp, fp, tn, fn int

			for _, l := range qLabels {
				key := fmt.Sprintf("%s#%s#%s", l.File, l.Heading, l.Question)
				prob, found := findMap[key]
				if !found {
					continue
				}

				pred := 0
				if prob >= thresh {
					pred = 1
				}

				if pred == 1 && l.Label == 1 {
					tp++
				} else if pred == 1 && l.Label == 0 {
					fp++
				} else if pred == 0 && l.Label == 0 {
					tn++
				} else if pred == 0 && l.Label == 1 {
					fn++
				}
			}

			total := tp + fp + tn + fn
			agreement := 0.0
			precision := 0.0
			recall := 0.0

			if total > 0 {
				agreement = float64(tp+tn) / float64(total)
			}
			if tp+fp > 0 {
				precision = float64(tp) / float64(tp+fp)
			}
			if tp+fn > 0 {
				recall = float64(tp) / float64(tp+fn)
			}

			fmt.Fprintf(w, "%-11.1f | %-5d | %-5d | %-5d | %-5d | %-9.1f%% | %-9.1f%% | %-9.1f%%\n",
				thresh, tp, fp, tn, fn, agreement*100, precision*100, recall*100)
		}
	}
	fmt.Fprintln(w, "================================================================================")
}
