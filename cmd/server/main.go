package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ShunsakuIsaji/dashboard_cuttle/handler"
)

func main() {
	// 1) テンプレートを起動時に読む
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("template", "header.html"),
		filepath.Join("template", "footer.html"),
		filepath.Join("template", "index.html"),
	))

	// 2) ルータを作る
	r := chi.NewRouter()

	// 3) 最低限あると便利なミドルウェア
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// 4) ルート登録
	r.Get("/index", handler.HandleIndex(tmpl))

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
