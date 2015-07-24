package views

import (
	"fmt"

	"github.com/gizak/termui"
	"github.com/rzh/montu/widgets"
)

type ViewRPS struct {
	mongosWidget widgets.RPSWidget

	replicas [9]widgets.RPSWidget

	chunks [3]widgets.ChunkWidget
}

var _rpsView *ViewRPS

func (v *ViewRPS) SwitchToView() View {
	ResetView()

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(3, 0, v.mongosWidget.Widget()),
			termui.NewCol(3, 0, v.chunks[0].Widget(), v.chunks[1].Widget(), v.chunks[2].Widget()),
			termui.NewCol(3, 0, widgets.Migration().Widget()),
		),
		termui.NewRow(
			termui.NewCol(3, 0, v.replicas[0].Widget()),
			termui.NewCol(3, 0, v.replicas[1].Widget()),
			termui.NewCol(3, 0, v.replicas[2].Widget()),
		),
		termui.NewRow(
			termui.NewCol(3, 0, v.replicas[3].Widget()),
			termui.NewCol(3, 0, v.replicas[4].Widget()),
			termui.NewCol(3, 0, v.replicas[5].Widget()),
		),
		termui.NewRow(
			termui.NewCol(3, 0, v.replicas[6].Widget()),
			termui.NewCol(3, 0, v.replicas[7].Widget()),
			termui.NewCol(3, 0, v.replicas[8].Widget()),
		),
	)

	// calculate layout
	termui.Body.Align()

	return v
}

func (v *ViewRPS) Render() {
	termui.Render(termui.Body)
}

func (v *ViewRPS) Init() {
	v.mongosWidget = widgets.NewClusterRPSWidget("mongos", "ms")

	for i := 0; i < len(v.replicas); i++ {
		v.replicas[i] = widgets.NewClusterRPSWidget("Repl_"+fmt.Sprint(int(i/3+1))+"_"+fmt.Sprint(i-int(i/3)*3), fmt.Sprintf("p%d", i+1))
	}

	for i := 0; i < len(v.chunks); i++ {
		v.chunks[i] = widgets.NewChunkWidget("sh_"+fmt.Sprint(i), i)
	}

}

func RPSView() *ViewRPS {
	return _rpsView
}

func init() {
	_rpsView = &ViewRPS{}
	_rpsView.Init()
}
