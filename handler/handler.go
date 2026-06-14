package handler

import (
	"html/template"

	"github.com/ShunsakuIsaji/dashboard_cuttle/internal/model"
)

type App struct {
	PriceRecords model.PriceRecords
	ItemMetas    model.ItemMetas
	Template     *template.Template
}
