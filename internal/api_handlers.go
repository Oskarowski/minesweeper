package internal

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"minesweeper/internal/db"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	rawData, err := h.Queries.GetGamesInfo(r.Context())

	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching DB information: %v", err), http.StatusInternalServerError)
		return
	}

	parsedData := make([]opts.PieData, 0)
	var colors = []string{"#28a745", "#dc3545", "#ffc107"}

	parsedData = append(parsedData, opts.PieData{Name: "Wins", Value: rawData.WonGames, ItemStyle: &opts.ItemStyle{Color: colors[0]}})
	parsedData = append(parsedData, opts.PieData{Name: "Losses", Value: rawData.LostGames, ItemStyle: &opts.ItemStyle{Color: colors[1]}})
	parsedData = append(parsedData, opts.PieData{Name: "Incomplete", Value: rawData.NotFinishedGames, ItemStyle: &opts.ItemStyle{Color: colors[2]}})

	pie := charts.NewPie()

	pie.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Wins vs Losses vs Incomplete",
		Subtitle: fmt.Sprintf("Total games: %v", rawData.TotalGames),
	}))

	pie.AddSeries("Game Status", parsedData).
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

	rawData, err := h.Queries.GetGamesPlayedPerGridSize(r.Context())

	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching grid size data: %v", err), http.StatusInternalServerError)
		return
	}

	gridSizes := make([]string, len(rawData))
	parsedBarData := make([]opts.BarData, len(rawData))

	// find the most popular grid size
	maxGamesPlayed := int64(0)
	for _, dbData := range rawData {
		if dbData.GamesPlayed > maxGamesPlayed {
			maxGamesPlayed = dbData.GamesPlayed
		}
	}

	for i, dbData := range rawData {
		gridSizes[i] = fmt.Sprintf("%vx%v", dbData.GridSize, dbData.GridSize)

		if dbData.GamesPlayed == maxGamesPlayed {
			parsedBarData[i] = opts.BarData{Value: dbData.GamesPlayed, ItemStyle: &opts.ItemStyle{Color: "#ffa500"}}
		} else {
			parsedBarData[i] = opts.BarData{Value: dbData.GamesPlayed}
		}
	}

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

	bar.SetXAxis(gridSizes).
		AddSeries("Games Played", parsedBarData).
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

func (h *ApiHandler) MinesAmountBarChart(w http.ResponseWriter, r *http.Request) {
	rawDbData, err := h.Queries.GetMinesPopularity(r.Context())

	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching mines popularity data: %v", err), http.StatusInternalServerError)
		return
	}

	// find the most popular amount of mines
	maxAmount := int64(0)
	for _, dbData := range rawDbData {
		if dbData.MinesCount > maxAmount {
			maxAmount = dbData.MinesCount
		}
	}

	minesPopularity := make([]string, len(rawDbData))
	parsedBarData := make([]opts.BarData, len(rawDbData))

	for i, dbData := range rawDbData {
		minesPopularity[i] = strconv.FormatInt(dbData.MinesAmount, 10)

		if dbData.MinesCount == maxAmount {
			parsedBarData[i] = opts.BarData{Value: dbData.MinesCount, ItemStyle: &opts.ItemStyle{Color: "#ffa500"}}
		} else {
			parsedBarData[i] = opts.BarData{Value: dbData.MinesCount}
		}
	}

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Amount of Mines Popularity",
			Subtitle: "Games with particular amount of mines",
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:  "slider",
			Start: 10,
			End:   75,
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(false),
		}),
	)

	bar.SetXAxis(minesPopularity).
		AddSeries("Mines Amount", parsedBarData).
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
