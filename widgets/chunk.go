package widgets

import (
	"fmt"
	"time"

	"github.com/gizak/termui"
	"github.com/rzh/montu/gen"
)

type ChunkWidget struct {
	title string
	shard int
	bc    *termui.Gauge
}

func max3(a [3]int) int {
	var m int = 0
	for i := 0; i < 3; i++ {
		if a[i] > m {
			m = a[i]
		}

	}
	return m
}

func total3(a [3]int) int {
	var m int = 0
	for i := 0; i < 3; i++ {
		m = a[i] + m
	}
	return m
}
func (w *ChunkWidget) getChunkCount(s gen.MongosStats) int {
	return (s.ChunkDistribution[w.shard] * 100) / total3(s.ChunkDistribution)
}

func (w *ChunkWidget) Init() {
	g := termui.NewGauge()
	g.Percent = 10
	g.Width = 50
	g.Height = 3
	g.Y = 11
	g.Border.Label = "> " + w.title + " chunks <"
	g.Label = "test"
	g.LabelAlign = termui.AlignRight
	g.BarColor = termui.ColorGreen
	// g.PercentColor = termui.ColorYellow

	w.bc = g

	go func() {
		for {
			// TODO: refresh data here via channel or something
			s := gen.MongosMonitor()

			if max3(s.GetStats().ChunkDistribution) == 0 {
				w.bc.Percent = 0
			} else {
				w.bc.Percent = w.getChunkCount(s.GetStats())
			}
			w.bc.Label = fmt.Sprintf("chunks: %d", s.GetStats().ChunkDistribution[w.shard])
			w.bc.Border.Label = "> " + w.title + " - " + w.bc.Label + " <"

			time.Sleep(1 * time.Second)
		}
	}()
}

func (w *ChunkWidget) Widget() termui.GridBufferer {
	return w.bc
}

func NewChunkWidget(title string, shard int) ChunkWidget {
	v := ChunkWidget{title: title, shard: shard}
	v.Init()

	return v
}
