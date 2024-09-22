package internal

import (
	"html/template"
	"math/rand"
	"minesweeper/internal/db"
	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/gorilla/sessions"
)

type ApiHandler struct {
	Templates *template.Template
	Store     *sessions.CookieStore
	Queries   *db.Queries
}

func NewApiHandler(templates *template.Template, store *sessions.CookieStore, queries *db.Queries) *ApiHandler {
	return &ApiHandler{templates, store, queries}
}

func (h *ApiHandler) PieWinsLossesIncompleteChart(w http.ResponseWriter, r *http.Request) {
	var (
		itemCntPie = 3
		options    = []string{"Wins", "Losses", "Incomplete"}
		colors     = []string{"#28a745", "#dc3545", "#ffc107"}
	)

	items := make([]opts.PieData, 0)
	for i := 0; i < itemCntPie; i++ {
		items = append(items, opts.PieData{Name: options[i], Value: rand.Intn(100), ItemStyle: &opts.ItemStyle{Color: colors[i]}})
	}

	pie := charts.NewPie()

	pie.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Wins vs Losses vs Incomplete",
		Subtitle: "Minesweeper Global Statistics",
	}))

	pie.AddSeries("Game Status", items).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Formatter: "{b}: {d}%", // Label formatter to show percentage
			}),
		)

	pie.Render(w)
}
