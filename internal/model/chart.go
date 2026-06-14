package model

import (
	"sort"
)

type ChartData struct {
	Labels   []string       `json:"labels"`
	DataSets []ChartDataSet `json:"datasets"`
}

type ChartDataSet struct {
	Label string     `json:"label"`
	Data  []*float64 `json:"data"`
}

func (records PriceRecords) BuildLabels() []string {
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

func (records PriceRecords) BuildDatasets(items []string, labels []string) []ChartDataSet {
	// item -> date -> value
	index := make(map[string]map[string]float64)
	for _, c := range items {
		index[c] = map[string]float64{}
	}
	for _, r := range records {
		if _, ok := index[r.Item]; !ok {
			continue
		}
		date := r.Date.Format("2006-01")
		index[r.Item][date] = r.Price
	}

	datasets := make([]ChartDataSet, 0, len(items))
	for _, c := range items {
		values := make([]*float64, 0, len(labels))
		for _, label := range labels {
			if value, ok := index[c][label]; ok {
				// capture value into new variable so pointer is stable
				v := value
				values = append(values, &v)
			} else {
				values = append(values, nil)
			}
		}
		datasets = append(datasets, ChartDataSet{
			Label: c,
			Data:  values,
		})
	}
	return datasets
}

func (records PriceRecords) BuildChartData(items []string) ChartData {
	labels := records.BuildLabels()
	datasets := records.BuildDatasets(items, labels)
	return ChartData{
		Labels:   labels,
		DataSets: datasets,
	}
}

func ConvertChartDataLabels(chartData ChartData, itemMetas ItemMetas) ChartData {
	for i, dataset := range chartData.DataSets {
		if itemMeta, ok := itemMetas.GetItemMeta(dataset.Label); ok {
			chartData.DataSets[i].Label = itemMeta.Label
		}
	}
	return chartData
}
