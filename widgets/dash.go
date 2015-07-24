package widgets

import (
	"fmt"
	"time"

	"github.com/gizak/termui"
	"github.com/rzh/montu/gen"
)

type DashWidget struct {
	title  string
	addr   string
	server *gen.ServerMonitor
	p1     *gen.ServerMonitor

	bc *termui.Par
}

const CLR_0 = "\x1b[30;1m"
const CLR_R = "\x1b[31;1m"
const CLR_G = "\x1b[32;1m"
const CLR_Y = "\x1b[33;1m"
const CLR_B = "\x1b[34;1m"
const CLR_M = "\x1b[35;1m"
const CLR_C = "\x1b[36;1m"
const CLR_W = "\x1b[37;1m"
const CLR_N = "\x1b[0m"

var _dash *DashWidget

func (w *DashWidget) Init() {
	w.server = gen.MonitorServer("ms")
	w.p1 = gen.MonitorServer("p1")

	par3 := termui.NewPar("")
	par3.Height = 17
	par3.Width = 137
	par3.Y = 9
	par3.Border.Label = "> Stats Dashboard: <"

	par3.Text = ""

	w.bc = par3

	go func() {
		for {
			// TODO: refresh data here via channel or something
			stats := w.server.GetStats()
			stats_p1 := w.p1.GetStats()

			/*
				w.bc.Text += fmt.Sprintf("  CPU        : %5d%%", stats.P_CPU)              ; w.bc.Text += fmt.Sprintf("  |   CPU        : %5d%%\n", stats_p1.P_CPU)
				w.bc.Text += fmt.Sprintf("  MEM        : %5d%%", stats.P_MEM)              ; w.bc.Text += fmt.Sprintf("  |   MEM        : %5d%%\n", stats_p1.P_MEM)
				w.bc.Text += fmt.Sprintf("  DSK(dbs)   : %5d%%", stats.P_DISK_dbs)         ; w.bc.Text += fmt.Sprintf("  |   DSK(dbs)   : %5d%%\n", stats_p1.P_DISK_dbs)
				w.bc.Text += fmt.Sprintf("  DSK(jnl)   : %5d%%", stats.P_DISK_journal)     ; w.bc.Text += fmt.Sprintf("  |   DSK(jnl)   : %5d%%\n", stats_p1.P_DISK_journal)
				w.bc.Text += fmt.Sprintf("  OPS        : %5d ", stats.OPS)                 ; w.bc.Text += fmt.Sprintf("  |   OPS        : %5d\n", stats_p1.OPS)
				w.bc.Text += fmt.Sprintf("   - insert  : %5d ", stats.OPS_insert)          ; w.bc.Text += fmt.Sprintf("  |    - insert  : %5d\n", stats_p1.OPS_insert)
				w.bc.Text += fmt.Sprintf("   - query   : %5d ", stats.OPS_query)           ; w.bc.Text += fmt.Sprintf("  |    - query   : %5d\n", stats_p1.OPS_query)
				w.bc.Text += fmt.Sprintf("   - update  : %5d ", stats.OPS_update)          ; w.bc.Text += fmt.Sprintf("  |    - update  : %5d\n", stats_p1.OPS_update)
				w.bc.Text += fmt.Sprintf("   - getmore : %5d ", stats.OPS_getmore)         ; w.bc.Text += fmt.Sprintf("  |    - getmore : %5d\n", stats_p1.OPS_getmore)
				w.bc.Text += fmt.Sprintf("   - delete  : %5d ", stats.OPS_delete)          ; w.bc.Text += fmt.Sprintf("  |    - delete  : %5d\n", stats_p1.OPS_delete)
				w.bc.Text += fmt.Sprintf("   - command : %5d ", stats.OPS_command)         ; w.bc.Text += fmt.Sprintf("  |    - command : %5d\n", stats_p1.OPS_command)
			*/
			w.bc.Text = "\nmongos:                |   Shard_1_primary:\n-----------------------|-----------------------\n"

			w.bc.Text += fmt.Sprintf("  CPU        : %5d%%", stats.P_CPU)
			w.bc.Text += fmt.Sprintf("  |   CPU        : %5d%%\n", stats_p1.P_CPU)
			w.bc.Text += fmt.Sprintf("  MEM        : %5d%%", stats.P_MEM)
			w.bc.Text += fmt.Sprintf("  |   MEM        : %5d%%\n", stats_p1.P_MEM)
			w.bc.Text += fmt.Sprintf("  DSK(dbs)   : %5d%%", stats.P_DISK_dbs)
			w.bc.Text += fmt.Sprintf("  |   DSK(dbs)   : %5d%%\n", stats_p1.P_DISK_dbs)
			w.bc.Text += fmt.Sprintf("  DSK(jnl)   : %5d%%", stats.P_DISK_journal)
			w.bc.Text += fmt.Sprintf("  |   DSK(jnl)   : %5d%%\n", stats_p1.P_DISK_journal)
			w.bc.Text += fmt.Sprintf("  OPS        : %5d ", stats.OPS)
			w.bc.Text += fmt.Sprintf("  |   OPS        : %5d\n", stats_p1.OPS)
			w.bc.Text += fmt.Sprintf("   - insert  : %5d ", stats.OPS_insert)
			w.bc.Text += fmt.Sprintf("  |    - insert  : %5d\n", stats_p1.OPS_insert)
			w.bc.Text += fmt.Sprintf("   - query   : %5d ", stats.OPS_query)
			w.bc.Text += fmt.Sprintf("  |    - query   : %5d\n", stats_p1.OPS_query)
			w.bc.Text += fmt.Sprintf("   - update  : %5d ", stats.OPS_update)
			w.bc.Text += fmt.Sprintf("  |    - update  : %5d\n", stats_p1.OPS_update)
			w.bc.Text += fmt.Sprintf("   - getmore : %5d ", stats.OPS_getmore)
			w.bc.Text += fmt.Sprintf("  |    - getmore : %5d\n", stats_p1.OPS_getmore)
			w.bc.Text += fmt.Sprintf("   - delete  : %5d ", stats.OPS_delete)
			w.bc.Text += fmt.Sprintf("  |    - delete  : %5d\n", stats_p1.OPS_delete)
			w.bc.Text += fmt.Sprintf("   - command : %5d ", stats.OPS_command)
			w.bc.Text += fmt.Sprintf("  |    - command : %5d\n", stats_p1.OPS_command)

			time.Sleep(1 * time.Second)
		}
	}()
}

func (w *DashWidget) Widget() termui.GridBufferer {
	return w.bc
}

func Dash() *DashWidget {
	return _dash
}

func init() {
	_dash = &DashWidget{}
	_dash.Init()
}
