package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/ShunsakuIsaji/dashboard_cuttle/handler"
	"github.com/ShunsakuIsaji/dashboard_cuttle/internal/gcs"
	"github.com/ShunsakuIsaji/dashboard_cuttle/internal/model"
)

func main() {
	EnvLoad()

	GCS_BUCKET_NAME := os.Getenv("GCS_BUCKET_NAME")
	GCS_FILE_NAME := os.Getenv("GCS_FILE_NAME")

	// 1) テンプレートを起動時に読む
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("templates", "header.gohtml"),
		filepath.Join("templates", "footer.gohtml"),
		filepath.Join("templates", "index.gohtml"),
	))

	ctx := context.Background()
	records, err := gcs.DownloadFromGCS(ctx, GCS_BUCKET_NAME, GCS_FILE_NAME)
	if err != nil {
		log.Fatalf("failed to download data from GCS: %v", err)
	}
	log.Printf("fetched data size: %d", len(records))

	itemMetas, err := model.LoadItemMetas("item_meta.yaml")
	if err != nil {
		log.Fatalf("failed to load item metas: %v", err)
	}

	app := &handler.App{
		PriceRecords: records,
		ItemMetas:    itemMetas,
		Template:     tmpl,
	}

	// 2) ルータを作る
	r := chi.NewRouter()

	// 3) 最低限あると便利なミドルウェア
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// 4) ルート登録
	r.Get("/index", app.HandleIndex())

	// ヘルスチェック
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	addr := ":8080"
	log.Printf("server started: http://localhost%s", addr)

	// 5) サーバ起動
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

func EnvLoad() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using environment variables")
	}
}
