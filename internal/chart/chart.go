package chart

import (
	"fmt"
	"io"
	"time"

	"github.com/google/go-github/v39/github"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func WriteStarsChart(stars []*github.Stargazer, w io.Writer) error {
	var series = chart.TimeSeries{
		Style: chart.Style{
			Show: true,
			StrokeColor: drawing.Color{
				R: 129,
				G: 199,
				B: 239,
				A: 255,
			},
			StrokeWidth: 2,
		},
	}
	for i, star := range stars {
		series.XValues = append(series.XValues, star.StarredAt.Time)
		series.YValues = append(series.YValues, float64(i))
	}
	if len(series.XValues) < 2 {
		series.XValues = append(series.XValues, time.Now())
		series.YValues = append(series.YValues, 1)
	}

	var graph = chart.Chart{
		XAxis: chart.XAxis{
			Name:      "Time",
			NameStyle: chart.StyleShow(),
			Style: chart.Style{
				Show:        true,
				StrokeWidth: 2,
				StrokeColor: drawing.Color{
					R: 85,
					G: 85,
					B: 85,
					A: 255,
				},
			},
		},
		YAxis: chart.YAxis{
			Name:      "Stargazers",
			NameStyle: chart.StyleShow(),
			Style: chart.Style{
				Show:        true,
				StrokeWidth: 2,
				StrokeColor: drawing.Color{
					R: 85,
					G: 85,
					B: 85,
					A: 255,
				},
			},
			ValueFormatter: IntValueFormatter,
		},
		Series: []chart.Series{series},
	}
	return graph.Render(chart.SVG, w)
}

func IntValueFormatter(v interface{}) string {
	return fmt.Sprintf("%.0f", v)
}
