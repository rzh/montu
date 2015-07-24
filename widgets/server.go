// to show ops/cpu of server with sparking line

package widgets

import (
	"math/rand"
	"time"

	"github.com/gizak/termui"
	"github.com/rzh/montu/gen"
)

type ServerWidget struct {
	title string
	addr  string

	bc *termui.Sparklines

	spl  []termui.Sparkline
	data [][]int

	server *gen.ServerMonitor
}

func makeSparkline(title string, data []int, color termui.Attribute) termui.Sparkline {
	s := termui.NewSparkline()
	s.Data = data
	s.Title = title
	s.LineColor = color
	return s
}

const SPARKLINE_COUNT = 5

func (w *ServerWidget) Init() {
	// something
	w.server = gen.MonitorServer(w.addr)

	w.data = make([][]int, SPARKLINE_COUNT, SPARKLINE_COUNT)
	w.spl = make([]termui.Sparkline, SPARKLINE_COUNT, SPARKLINE_COUNT)

	for i := 0; i < SPARKLINE_COUNT; i++ {
		w.data[i] = make([]int, 60, 60)
	}

	w.spl[0] = makeSparkline(" Ops/s", w.data[0], termui.ColorGreen)
	w.spl[1] = makeSparkline(" % CPU", w.data[1], termui.ColorRed)
	w.spl[2] = makeSparkline(" % DISK_dbs", w.data[2], termui.ColorMagenta)
	w.spl[3] = makeSparkline(" % DISK_journal", w.data[3], termui.ColorCyan)

	// group
	spls1 := termui.NewSparklines()
	spls1.Height = 10
	spls1.Width = 20
	spls1.Y = 3
	spls1.Border.Label = "> " + w.title + " <"

	w.bc = spls1

	go func() {
		for {
			stat := w.server.GetStats()

			w.data[0] = append(w.data[0][1:len(w.data[0])], stat.OPS)               // OPS
			w.data[1] = append(w.data[1][1:len(w.data[1])], stat.P_CPU)             // CPU
			w.data[2] = append(w.data[2][1:len(w.data[2])], stat.P_DISK_dbs)        // DISK_dbs
			w.data[3] = append(w.data[3][1:len(w.data[3])], stat.P_DISK_journal)    // DISK_journal
			w.data[4] = append(w.data[4][1:len(w.data[4])], int(rand.Int63n(90))+5) // NET

			// set max/min to 100, FIXME
			w.data[0][0] = 100
			w.data[1][0] = 100
			w.data[2][0] = 100
			w.data[3][0] = 100
			w.data[4][0] = 100

			for i := 0; i < SPARKLINE_COUNT; i++ {
				w.spl[i].Data = w.data[i]
			}

			w.bc.Lines = w.spl
			time.Sleep(1 * time.Second)
		}
	}()
}

func (w *ServerWidget) Widget() termui.GridBufferer {
	return w.bc
}

func NewServerWidget(title string, addr string) *ServerWidget {
	w := ServerWidget{title: title, addr: addr}
	w.Init()

	return &w
}
