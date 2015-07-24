package widgets

import (
	"fmt"
	"strings"
	"time"

	"github.com/gizak/termui"
	"github.com/rzh/montu/gen"
)

type StatusBarWidget struct {
	msg  []string
	addr string
	bc   *termui.Par
	G    *gen.ServerMonitor
}

var _statusBar *StatusBarWidget

func StatusBar() *StatusBarWidget {
	return _statusBar
}

func (w *StatusBarWidget) Init() {
	w.addr = "p4" // FIXME
	w.G = gen.MonitorServer(w.addr)

	par3 := termui.NewPar("")
	par3.Height = LINE + 2
	par3.Width = 137
	par3.Y = 9
	par3.Border.Label = "> Events <"

	par3.Text = ""

	w.bc = par3

	go func() {
		for {
			time.Sleep(1 * time.Second)

			// request data
			stat := w.G.GetStats()
			m := fmt.Sprintf("msg: %s CPU: %d DSK_dbs: %d DSK_j: %d", stat.LastLine, stat.P_CPU, stat.P_DISK_dbs, stat.P_DISK_journal)

			if len(w.msg) == LINE {
				w.msg = append(w.msg[1:LINE], m)
			} else {
				w.msg = append(w.msg, m)
			}

			text := strings.Join(w.msg, "\n ")
			w.bc.Text = " " + text
		}
	}()
}

func (w *StatusBarWidget) Widget() termui.GridBufferer {
	return w.bc
}

func init() {
	_statusBar = &StatusBarWidget{}
	_statusBar.Init()
}
