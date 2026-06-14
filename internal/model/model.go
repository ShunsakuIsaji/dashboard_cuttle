package model

import (
	"sort"
	"time"
)

type PriceData struct {
	Date     time.Time `json:"date"`
	Price    float64   `json:"price"`
	Unit     string    `json:"unit"`
	Category string    `json:"category"`
	Item     string    `json:"item"`
}

type PriceRecords []PriceData

func (r PriceData) ParseDateToStr() string {
	return r.Date.Format("2006-01-02")
}

func (records PriceRecords) GetDataByCategory(category string) PriceRecords {
	var result PriceRecords
	for _, r := range records {
		if r.Category == category {
			result = append(result, r)
		}
	}
	return result
}

func (records PriceRecords) GetDataByDate(since, until time.Time) PriceRecords {
	var result PriceRecords
	for _, r := range records {
		if (r.Date.Equal(since) || r.Date.After(since)) && (r.Date.Equal(until) || r.Date.Before(until)) {
			result = append(result, r)
		}
	}
	return result
}

func (records PriceRecords) SortDataByDate() {
	// 日付の古い順にソート
	sort.Slice(records, func(i, j int) bool {
		return records[i].Date.Before(records[j].Date)
	})
}

func (records PriceRecords) GetLatest() PriceRecords {
	latestMap := make(map[string]PriceData)

	for _, r := range records {
		cur, ok := latestMap[r.Item]
		if !ok || cur.Date.Before(r.Date) {
			latestMap[r.Item] = PriceData{
				Date:     r.Date,
				Price:    r.Price,
				Unit:     r.Unit,
				Category: r.Category,
				Item:     r.Item,
			}
		}
	}

	result := make(PriceRecords, 0, len(latestMap))
	for _, v := range latestMap {
		result = append(result, v)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Item < result[j].Item
	})

	return result
}

func (records PriceRecords) GetDataByItems(items []string) PriceRecords {
	itemSet := make(map[string]struct{})
	for _, item := range items {
		itemSet[item] = struct{}{}
	}

	var result PriceRecords
	for _, r := range records {
		if _, ok := itemSet[r.Item]; ok {
			result = append(result, r)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if records[i].Item < records[j].Item {
			return true
		} else if records[i].Item == records[j].Item {
			return records[i].Date.Before(records[j].Date)
		}
		return false
	})

	return result
}
