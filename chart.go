package main

import (
	"fmt"
	"io"
	"time"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func writeStarsChart(stars []stargazer, w io.Writer) error {
	series := chart.TimeSeries{
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
		series.XValues = append(series.XValues, star.StarredAt)
		series.YValues = append(series.YValues, float64(i))
	}
	if len(series.XValues) < 2 {
		series.XValues = append(series.XValues, time.Now())
		series.YValues = append(series.YValues, 1)
	}

	graph := chart.Chart{
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
			ValueFormatter: intValueFormatter,
		},
		Series: []chart.Series{series},
	}
	return graph.Render(chart.SVG, w)
}

func intValueFormatter(v any) string {
	return fmt.Sprintf("%.0f", v)
}
