package gcs

import (
	"bytes"
	"context"
	"encoding/csv"
	"io"
	"log/slog"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/ShunsakuIsaji/dashboard_cuttle/internal/model"
)

func DownloadFromGCS(ctx context.Context, bucket, filename string) (model.PriceRecords, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	r, err := client.Bucket(bucket).Object(filename).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	result, err := ByteToCSV(data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ByteToCSV(data []byte) (model.PriceRecords, error) {
	var result model.PriceRecords

	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for i, record := range records {
		if i == 0 {
			// ヘッダー行は飛ばす
			continue
		}
		// date,category,item,value,unitの順に並んでいる
		// valueのパース
		value, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			slog.Info("failed to parse value, skipping row", "row", i, "error", err)
			continue
		}

		// dateのパース
		date, err := time.Parse("2006-01-02", record[0]+"-01")
		if err != nil {
			slog.Info("failed to parse date, skipping row", "row", i, "error", err)
		}

		result = append(result, model.PriceData{
			Date:     date,
			Price:    value,
			Unit:     record[4],
			Category: record[1],
			Item:     record[2],
		})
	}
	return result, nil
}
