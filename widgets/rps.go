package widgets

import (
	"math/rand"
	"time"

	"github.com/gizak/termui"
	"github.com/rzh/montu/gen"
)

type Widget interface {
	Widget() termui.GridBufferer
}

type RPSWidget struct {
	title  string
	addr   string
	server *gen.ServerMonitor

	bc *termui.BarChart
}

func (w *RPSWidget) Init() {
	w.server = gen.MonitorServer(w.addr)

	bc := termui.NewBarChart()
	data := []int{3, 2, 5, 3, 9, 5, 3, 2, 5, 8, 3, 2, 4, 5, 3, 2, 5, 7, 5, 3, 2, 6, 7, 4, 6, 3, 6, 7, 8, 3, 6, 4, 5, 3, 2, 4, 6, 4, 8, 5, 9, 4, 3, 6, 5, 3, 6}
	bclabels := []string{"INS", "QRY", "UPD", "DEL", "GTM", "", "CPU", "MEM", "DSK"}
	bc.Border.Label = "> " + w.title + " <"
	bc.Border.FgColor = termui.ColorCyan
	bc.Data = data
	bc.Width = 30
	bc.Height = 10
	bc.SetMax(100)
	bc.DataLabels = bclabels
	//	bc.BgColor = termui.ColorBlack
	bc.TextColor = termui.ColorYellow
	bc.BarColor = termui.ColorRed
	bc.NumColor = termui.ColorRed

	w.bc = bc

	go func() {
		for {
			// TODO: refresh data here via channel or something
			var data []int
			stats := w.server.GetStats()

			for i := 0; i < len(bclabels); i++ {
				if bclabels[i] != "" {
					switch bclabels[i] {

					case "CPU":
						c := stats.P_CPU
						if stats.OPS > 100 && c > 10 {
							c = int(float64(stats.OPS) * (float64(c) / 100.0))
						}
						data = append(data, c)
					case "MEM":
						c := stats.P_MEM
						if stats.OPS > 100 {
							c = int(float64(stats.OPS) * (float64(c) / 100.0))
						}
						data = append(data, c)
					case "INS":
						data = append(data, stats.OPS_insert)
					case "QRY":
						data = append(data, stats.OPS_query)
					case "UPD":
						data = append(data, stats.OPS_update)
					case "GTM":
						data = append(data, stats.OPS_getmore)
					case "DEL":
						data = append(data, stats.OPS_delete)
					case "DSK":
						c := (stats.P_DISK_dbs + stats.P_DISK_journal) / 2
						if stats.OPS > 100 {
							c = int(float64(stats.OPS) * (float64(c) / 100.0))
						}
						data = append(data, c)
					default:
						data = append(data, rand.Intn(100))
					}
				} else {
					data = append(data, 0)
				}
			}

			w.bc.Data = data
			time.Sleep(1 * time.Second)
		}
	}()
}

func (w *RPSWidget) Widget() termui.GridBufferer {
	return w.bc
}

func NewClusterRPSWidget(title string, addr string) RPSWidget {
	v := RPSWidget{title: title, addr: addr}
	v.Init()

	return v
}
