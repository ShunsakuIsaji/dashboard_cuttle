package handler

import (
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"sort"
	"time"

	"github.com/ShunsakuIsaji/dashboard_cuttle/internal/model"

	"github.com/ShunsakuIsaji/dashboard_cuttle/internal/csv"
)

type LatestPrice struct {
	Category string
	Date     time.Time
	Value    float64
	Unit     string
}

type ChartData struct {
	Labels   []string       `json:"labels"`
	DataSets []ChartDataSet `json:"datasets"`
}

type ChartDataSet struct {
	Label string    `json:"label"`
	Data  []float64 `json:"data"`
}

func HandleIndex(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		record, err := csv.ReadCattlePrices(filepath.Join("data", "cattle_prices_full_vs.csv"))
		if err != nil {
			slog.Error("failed to read cattle prices", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// 最新の価格を取得
		latestPrices := getLatest(record)

		if err := tmpl.ExecuteTemplate(w, "index.html", latestPrices); err != nil {
			slog.Error("failed to execute template", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func getLatest(records []model.CattlePrice) []LatestPrice {
	latestMap := make(map[string]LatestPrice)

	for _, r := range records {
		cur, ok := latestMap[r.Category]
		if !ok || cur.Date.Before(r.Date) {
			latestMap[r.Category] = LatestPrice{
				Category: r.Category,
				Date:     r.Date,
				Value:    r.Price,
				Unit:     r.Unit,
			}
		}
	}

	result := make([]LatestPrice, 0, len(latestMap))
	for _, v := range latestMap {
		result = append(result, v)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Category < result[j].Category
	})

	return result
}

func BuildLabels(records []model.CattlePrice) []string {
	m := map[string]struct{}{}
	for _, r := range records {
		m[r.Date.Format("2006-01")] = struct{}{}
	}

	labels := make([]string, 0, len(m))
	for k := range m {
		labels = append(labels, k)
	}
	sort.Strings(labels)
	return labels
}

func BuildDatasets(records []model.CattlePrice, categories []string, labels []string) []ChartDataSet {
	// category -> date -> value
	index := make(map[string]map[string]float64)
	for _, c := range categories {
		index[c] = map[string]float64{}
	}
	for _, r := range records {
		date := r.Date.Format("2006-01")
		index[r.Category][date] = r.Price
	}

	datasets := make([]ChartDataSet, 0, len(categories))
	for _, c := range categories {
		values := make([]float64, 0, len(labels))
		for _, label := range labels {
			values = append(values, index[c][label])
		}
		datasets = append(datasets, ChartDataSet{
			Label: c,
			Data:  values,
		})
	}
	return datasets
}
