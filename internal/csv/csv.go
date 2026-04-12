package csv

import (
	"encoding/csv"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/ShunsakuIsaji/dashboard_cuttle/internal/model"
)

func ReadCattlePrices(filePath string) ([]model.CattlePrice, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var prices []model.CattlePrice

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for i, record := range records {
		// ヘッダー行はスキップ
		if i == 0 {
			continue
		}

		// priceのparse
		price, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			// parseに失敗したら、その行はスキップする
			slog.Info("failed to parse price, skipping row", "row", i, "error", err)
			continue
		}

		// Dateのparse
		date, err := time.Parse("2006-01-02", record[0]+"-01")
		if err != nil {
			slog.Info("failed to parse date, skipping row", "row", i, "error", err)
			continue
		}

		prices = append(prices, model.CattlePrice{
			Date:     date,
			Price:    price,
			Unit:     record[2],
			Category: record[3],
		})
	}

	return prices, nil
}
