package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type RequestResult struct {
	ID       int    `json:"id"`
	URL      string `json:"url"`
	Status   string `json:"status"`
	Error    string `json:"error,omitempty"`
	Duration string `json:"duration"`
}

type CycleStats struct {
	CycleNumber    int         `json:"cycle_number"`
	TotalRequests  int         `json:"total_requests"`
	SuccessCount   int         `json:"success_count"`
	ErrorCount     int         `json:"error_count"`
	AvgDurationMS  string      `json:"avg_duration_ms"`
	AvgDurationSec string      `json:"avg_duration_sec"`
	AvgDurationMin string      `json:"avg_duration_min"`
	MinDurationMS  string      `json:"min_duration_ms"`
	MinDurationSec string      `json:"min_duration_sec"`
	MinDurationMin string      `json:"min_duration_min"`
	MaxDurationMS  string      `json:"max_duration_ms"`
	MaxDurationSec string      `json:"max_duration_sec"`
	MaxDurationMin string      `json:"max_duration_min"`
	ResponseCodes  map[int]int `json:"response_codes"`
}

type Report struct {
	TotalCycles    int          `json:"total_cycles"`
	TotalRequests  int          `json:"total_requests"`
	SuccessCount   int          `json:"success_count"`
	ErrorCount     int          `json:"error_count"`
	AvgDurationMS  string       `json:"avg_duration_ms"`
	AvgDurationSec string       `json:"avg_duration_sec"`
	AvgDurationMin string       `json:"avg_duration_min"`
	MinDurationMS  string       `json:"min_duration_ms"`
	MinDurationSec string       `json:"min_duration_sec"`
	MinDurationMin string       `json:"min_duration_min"`
	MaxDurationMS  string       `json:"max_duration_ms"`
	MaxDurationSec string       `json:"max_duration_sec"`
	MaxDurationMin string       `json:"max_duration_min"`
	ResponseCodes  map[int]int  `json:"response_codes"`
	CycleDetails   []CycleStats `json:"cycle_details"`
}

func checkStatus(url string, id int, wg *sync.WaitGroup, results chan<- RequestResult) {
	defer wg.Done()

	start := time.Now()
	resp, err := http.Get(url)
	duration := time.Since(start)

	result := RequestResult{
		ID:       id,
		URL:      url,
		Duration: duration.String(),
	}

	if err != nil {
		result.Error = err.Error()
		results <- result
		return
	}
	defer resp.Body.Close()

	result.Status = resp.Status
	results <- result
}

func generateReport(report Report, filename string) error {
	// Converte o relatório para JSON
	reportJSON, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao gerar relatório JSON: %v", err)
	}

	// Salva o relatório em um arquivo
	err = ioutil.WriteFile(filename, reportJSON, 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar relatório em arquivo: %v", err)
	}

	fmt.Printf("Relatório detalhado gerado com sucesso: %s\n", filename)
	return nil
}

func main() {
	url := "http://localhost:8025/" // URL do site
	numGoroutines := 200            // Número de requisições simultâneas por ciclo
	maxCycles := 3000               // Número máximo de ciclos (para evitar loop infinito)

	var totalDurations []time.Duration
	report := Report{
		ResponseCodes: make(map[int]int),
		CycleDetails:  []CycleStats{},
	}

	for cycle := 1; cycle <= maxCycles; cycle++ {
		var wg sync.WaitGroup
		results := make(chan RequestResult, numGoroutines)

		fmt.Printf("Iniciando ciclo %d...\n", cycle)

		// Inicia as goroutines
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go checkStatus(url, i, &wg, results)
		}

		wg.Wait()
		close(results)

		// Processa os resultados do ciclo
		cycleStats := CycleStats{
			CycleNumber:   cycle,
			ResponseCodes: make(map[int]int),
		}
		var cycleDurations []time.Duration

		for result := range results {
			report.TotalRequests++
			cycleStats.TotalRequests++

			if result.Error == "" {
				report.SuccessCount++
				cycleStats.SuccessCount++

				// Extrai o código HTTP da resposta
				var statusCode int
				fmt.Sscanf(result.Status, "%d", &statusCode)
				report.ResponseCodes[statusCode]++
				cycleStats.ResponseCodes[statusCode]++

				// Calcula duração
				duration, _ := time.ParseDuration(result.Duration)
				totalDurations = append(totalDurations, duration)
				cycleDurations = append(cycleDurations, duration)
			} else {
				report.ErrorCount++
				cycleStats.ErrorCount++
			}
		}

		// Calcula estatísticas do ciclo
		if len(cycleDurations) > 0 {
			avg := calculateAverageDuration(cycleDurations)
			min := calculateMinDuration(cycleDurations)
			max := calculateMaxDuration(cycleDurations)

			cycleStats.AvgDurationMS = avg.String()
			cycleStats.AvgDurationSec = formatSeconds(avg)
			cycleStats.AvgDurationMin = formatMinutes(avg)

			cycleStats.MinDurationMS = min.String()
			cycleStats.MinDurationSec = formatSeconds(min)
			cycleStats.MinDurationMin = formatMinutes(min)

			cycleStats.MaxDurationMS = max.String()
			cycleStats.MaxDurationSec = formatSeconds(max)
			cycleStats.MaxDurationMin = formatMinutes(max)
		}

		report.CycleDetails = append(report.CycleDetails, cycleStats)
	}

	// Calcula estatísticas globais
	if len(totalDurations) > 0 {
		avg := calculateAverageDuration(totalDurations)
		min := calculateMinDuration(totalDurations)
		max := calculateMaxDuration(totalDurations)

		report.AvgDurationMS = avg.String()
		report.AvgDurationSec = formatSeconds(avg)
		report.AvgDurationMin = formatMinutes(avg)

		report.MinDurationMS = min.String()
		report.MinDurationSec = formatSeconds(min)
		report.MinDurationMin = formatMinutes(min)

		report.MaxDurationMS = max.String()
		report.MaxDurationSec = formatSeconds(max)
		report.MaxDurationMin = formatMinutes(max)
	}

	// Gera o relatório final
	filename := "api_test_detailed_report.json"
	err := generateReport(report, filename)
	if err != nil {
		fmt.Println("Erro ao gerar relatório:", err)
	}
}

// Funções auxiliares para calcular estatísticas de duração
func calculateAverageDuration(durations []time.Duration) time.Duration {
	var total time.Duration
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}

func calculateMinDuration(durations []time.Duration) time.Duration {
	min := durations[0]
	for _, d := range durations {
		if d < min {
			min = d
		}
	}
	return min
}

func calculateMaxDuration(durations []time.Duration) time.Duration {
	max := durations[0]
	for _, d := range durations {
		if d > max {
			max = d
		}
	}
	return max
}

// Funções para formatar durações
func formatSeconds(duration time.Duration) string {
	seconds := duration.Seconds()
	return fmt.Sprintf("%.3f segundos", seconds)
}

func formatMinutes(duration time.Duration) string {
	minutes := duration.Minutes()
	return fmt.Sprintf("%.3f minutos", minutes)
}
