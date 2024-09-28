package internal

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"math/rand"
	"minesweeper/internal/db"
	"net/http"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/gorilla/sessions"

	chartrender "github.com/go-echarts/go-echarts/v2/render"
)

type ApiHandler struct {
	Templates *template.Template
	Store     *sessions.CookieStore
	Queries   *db.Queries
}

func NewApiHandler(templates *template.Template, store *sessions.CookieStore, queries *db.Queries) *ApiHandler {
	return &ApiHandler{templates, store, queries}
}

// renderToHtml renders a chart as a template.HTML value.
//
// The argument should be a go-echarts chart that implements the Renderer interface.
// The rendered chart is returned as a template.HTML value.
//
// If the chart fails to render, an error is returned.
func renderToHtml(c interface{}) (template.HTML, error) {
	r, ok := c.(chartrender.Renderer)
	if !ok {
		return "", fmt.Errorf("provided chart does not implement the Renderer interface")
	}

	var buf bytes.Buffer

	err := r.Render(&buf)
	if err != nil {
		return "", fmt.Errorf("failed to render chart: %v", err)

	}

	htmlContent := buf.String()

	// Remove the <head>, <title>, <script> tags from the rendered chart
	// TODO use net/html package with RemoveChild to remove the tags
	htmlContent = strings.ReplaceAll(htmlContent, "<head>", "")
	htmlContent = strings.ReplaceAll(htmlContent, "</head>", "")
	htmlContent = strings.ReplaceAll(htmlContent, "<title>Awesome go-echarts</title>", "")
	htmlContent = strings.ReplaceAll(htmlContent, "<script src=\"https://go-echarts.github.io/go-echarts-assets/assets/echarts.min.js\"></script>", "")

	return template.HTML(htmlContent), nil
}

func (h *ApiHandler) PieWinsLossesIncompleteChart(w http.ResponseWriter, r *http.Request) {
	var (
		itemCntPie = 3
		options    = []string{"Wins", "Losses", "Incomplete"}
		colors     = []string{"#28a745", "#dc3545", "#ffc107"}
	)

	// TODO remove this placeholder data
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

	htmlPieSnipper, err := renderToHtml(pie)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering chart: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(htmlPieSnipper))
}

func (h *ApiHandler) GridSizeBar(w http.ResponseWriter, r *http.Request) {
	// TODO remove this placeholder data
	gridSizes := []string{"3x3", "4x4", "5x5", "6x6", "7x7"}
	gamesPlayed := []int{15, 40, 25, 10, 5}

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Grid Size Popularity",
			Subtitle: "Games played per grid size",
		}), charts.WithDataZoomOpts(opts.DataZoom{
			Type:  "slider",
			Start: 0,
			End:   100,
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(false),
		}),
	)

	items := make([]opts.BarData, 0)
	for _, v := range gamesPlayed {
		items = append(items, opts.BarData{Value: v})
	}

	bar.SetXAxis(gridSizes).
		AddSeries("Games Played", items).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{Show: opts.Bool(true)}),
		)

	htmlBarSnippet, err := renderToHtml(bar)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering chart: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(htmlBarSnippet))
}
func (h *ApiHandler) PlayedGamesInMonthBarChart(w http.ResponseWriter, r *http.Request) {
	pickedDate := r.URL.Query().Get("picked-date-range")

	if pickedDate == "" {
		now := time.Now()
		pickedDate = now.Format("2006-01")
	}

	parsedDate, err := time.Parse("2006-01", pickedDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid date format: %v", err), http.StatusBadRequest)
		return
	}

	startOfMonth := parsedDate
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	// Create sql.NullTime for start and end of the month
	startTime := sql.NullTime{
		Time:  startOfMonth,
		Valid: true,
	}

	endTime := sql.NullTime{
		Time:  endOfMonth,
		Valid: true,
	}

	gamesPerDay, err := h.Queries.GetGamesByMonthYearGroupedByDay(r.Context(), db.GetGamesByMonthYearGroupedByDayParams{
		CreatedAt:   startTime,
		CreatedAt_2: endTime,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching games: %v", err), http.StatusInternalServerError)
		return
	}

	days := []string{}
	gamesPlayed := []opts.BarData{}

	for _, gameDay := range gamesPerDay {
		days = append(days, fmt.Sprintf("%v", gameDay.Day))
		gamesPlayed = append(gamesPlayed, opts.BarData{Value: gameDay.GamesPlayed})
	}

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    fmt.Sprintf("Games Per Day in %v", parsedDate.Format("January 2006")),
			Subtitle: "Amount of games played per day",
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:  "slider",
			Start: 0,
			End:   100,
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(false),
		}),
	)

	bar.SetXAxis(days).
		AddSeries("Games Played", gamesPlayed).
		SetSeriesOptions(charts.WithLabelOpts(opts.Label{Show: opts.Bool(true)}))

	htmlBarSnippet, err := renderToHtml(bar)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering chart: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(htmlBarSnippet))
}
