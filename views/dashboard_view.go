package views

import (
	"github.com/gizak/termui"
	"github.com/rzh/montu/widgets"
)

type ViewEvent struct {
	mongosWidget   widgets.RPSWidget
	replica1Widget widgets.RPSWidget

	mongosWidget_2   widgets.RPSWidget
	replica1Widget_2 widgets.RPSWidget

	dash *widgets.DashWidget
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

var _eventView *ViewEvent

func (v *ViewEvent) SwitchToView() View {
	ResetView()

	statusRow := termui.NewRow(
		termui.NewCol(12, 0, v.dash.Widget()),
	)

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(1, 0, v.mongosWidget.Widget()),
			termui.NewCol(1, 0, v.replica1Widget.Widget()),
		),
		termui.NewRow(
			termui.NewCol(1, 0, v.mongosWidget_2.Widget()),
			termui.NewCol(1, 0, v.replica1Widget_2.Widget()),
		),
		statusRow,
	)

	// calculate layout
	termui.Body.Align()

	return v
}

func (v *ViewEvent) Render() {
	termui.Render(termui.Body)
}

func (v *ViewEvent) Init() {
	v.mongosWidget = widgets.NewClusterRPSWidget("mongos", "ms")
	v.replica1Widget = widgets.NewClusterRPSWidget("Rpl1", "p1")

	v.mongosWidget_2 = widgets.NewClusterRPSWidget("MongoS", "p4")
	v.replica1Widget_2 = widgets.NewClusterRPSWidget("Rpl1", "p7")

	v.dash = widgets.Dash()
}

func EventView() *ViewEvent {
	return _eventView
}

func init() {
	_eventView = &ViewEvent{}
	_eventView.Init()
}
