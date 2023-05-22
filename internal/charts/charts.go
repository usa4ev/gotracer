package charts

import (
	"io"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/usa4ev/gotracer/internal/model"
)

func DrawChart(w io.Writer, items map[string][]model.Entry, title opts.Title) {
	page := components.NewPage()
	page.AddCharts(
		barBasic(items, title),
	)

	page.Render(w)
}

const itemCnt = 24

var hours = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}

func barBasic(data map[string][]model.Entry, title opts.Title) *charts.Bar {
	bar := charts.NewBar()

	title.Left = "15%"

	bar.SetGlobalOptions(
		charts.WithTitleOpts(title),
		charts.WithToolboxOpts(opts.Toolbox{Show: true}),
		charts.WithLegendOpts(opts.Legend{Show: true, Left: "50%"}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "hours",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "rate", // AxisLabel: &opts.AxisLabel{VerticalAlign: "Bottom"},
		}),
	)

	bar = bar.SetXAxis(hours)

	for k, v := range data {
		bar = bar.AddSeries(k, makeBarItems(v))
	}

	return bar
}

func makeBarItems(items []model.Entry) []opts.BarData {
	res := make([]opts.BarData, itemCnt)

	for _,item := range items{
		res[item.Time.Hour()].Value = item.Count
	}

	return res
}

func NewTitle() opts.Title {
	return opts.Title{}
}
