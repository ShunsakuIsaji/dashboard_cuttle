package handler

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/ShunsakuIsaji/dashboard_cuttle/internal/model"
)

func (app *App) HandleMilk() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		latest := app.PriceRecords.GetLatest()

		MilkCategriesNames := app.ItemMetas.GetItemNamesByCategory("milk")
		chartData := app.PriceRecords.BuildChartData(MilkCategriesNames)

		payload := struct {
			Labels   []string                 `json:"labels"`
			Datasets []map[string]interface{} `json:"datasets"`
		}{
			Labels:   chartData.Labels,
			Datasets: make([]map[string]interface{}, 0, len(chartData.DataSets)),
		}

		for _, ds := range chartData.DataSets {
			id := ds.Label
			display := id
			if meta, ok := app.ItemMetas.GetItemMeta(id); ok {
				display = meta.Label
			}
			values := make([]interface{}, 0, len(ds.Data))
			for _, vp := range ds.Data {
				if vp == nil {
					values = append(values, nil)
				} else {
					values = append(values, *vp)
				}
			}
			payload.Datasets = append(payload.Datasets, map[string]interface{}{
				"label":   display,
				"data":    values,
				"yAxisID": id,
			})
		}

		b, err := json.Marshal(payload)
		if err != nil {
			slog.Error("failed to marshal chart payload", "error", err)
		}

		data := struct {
			Title                 string
			Latest_milk           model.PriceData
			Latest_wagyu_cow_a5   model.PriceData
			Latest_wagyu_female   model.PriceData
			Latest_compound_price model.PriceData
			ChartJSON             template.JS
		}{
			Title:                 "牛肉価格ダッシュボード",
			Latest_milk:           latest.GetDataByItems([]string{"milk_price"})[0],
			Latest_wagyu_cow_a5:   latest.GetDataByItems([]string{"wagyu_cow_a5"})[0],
			Latest_wagyu_female:   latest.GetDataByItems([]string{"wagyu_female"})[0],
			Latest_compound_price: latest.GetDataByItems([]string{"compound_price"})[0],
			ChartJSON:             template.JS(string(b)),
		}

		if err := app.Template.ExecuteTemplate(w, "milk", data); err != nil {
			slog.Error("failed to execute template", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
