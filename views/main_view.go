package views

import (
	"fmt"

	"github.com/gizak/termui"
	"github.com/rzh/montu/widgets"
)

type View interface {
	SwitchToView() View
	Render()
}

type ViewMain struct {
	// list of widgets
	mongosWidget     widgets.RPSWidget
	rpsWidgets       [10]widgets.RPSWidget
	replica1Widget   widgets.RPSWidget
	replica1Widget_2 widgets.RPSWidget
	replica1Widget_3 widgets.RPSWidget
	serverW_mongos   *widgets.ServerWidget
	shards           [3][3]*widgets.ServerWidget
	statusBar        *widgets.StatusBarWidget
}

var _mainView *ViewMain

func (v *ViewMain) SwitchToView() View {

	ResetView()

	statusRow := termui.NewRow(
		termui.NewCol(12, 0, v.statusBar.Widget()),
	)

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(3, 0, v.serverW_mongos.Widget()),
			termui.NewCol(3, 0, v.shards[0][0].Widget()),
			termui.NewCol(3, 0, v.shards[1][0].Widget()),
			termui.NewCol(3, 0, v.shards[2][0].Widget()),
		),
		termui.NewRow(
			termui.NewCol(3, 0, nil),
			termui.NewCol(3, 0, v.shards[0][1].Widget()),
			termui.NewCol(3, 0, v.shards[1][1].Widget()),
			termui.NewCol(3, 0, v.shards[2][1].Widget()),
		),
		termui.NewRow(
			termui.NewCol(3, 0, nil),
			termui.NewCol(3, 0, v.shards[0][2].Widget()),
			termui.NewCol(3, 0, v.shards[1][2].Widget()),
			termui.NewCol(3, 0, v.shards[2][2].Widget()),
		),
		/*
			termui.NewRow(
				termui.NewCol(1, 0, v.rpsWidgets[0].Widget()),
				termui.NewCol(1, 1, v.rpsWidgets[1].Widget()),
				termui.NewCol(1, 0, v.rpsWidgets[2].Widget()),
				termui.NewCol(1, 0, v.rpsWidgets[3].Widget()),
				termui.NewCol(1, 0, v.rpsWidgets[4].Widget()),
				termui.NewCol(1, 0, v.rpsWidgets[5].Widget()),
				termui.NewCol(1, 0, v.rpsWidgets[6].Widget()),
				termui.NewCol(1, 0, v.rpsWidgets[7].Widget()),
				termui.NewCol(1, 0, v.rpsWidgets[8].Widget()),
				termui.NewCol(1, 0, v.rpsWidgets[9].Widget()),
			),
		*/
		statusRow,
	)

	// calculate layout
	termui.Body.Align()

	return v
}

func (v *ViewMain) Init() {
	v.rpsWidgets[0] = widgets.NewClusterRPSWidget("MongoS", "ms")

	for i := 1; i < len(v.rpsWidgets); i++ {
		v.rpsWidgets[i] = widgets.NewClusterRPSWidget("MongoD_"+fmt.Sprint(i), fmt.Sprintf("p%d", i))
	}

	v.serverW_mongos = widgets.NewServerWidget("MongoS", "ms")

	for i := 0; i < 3; i++ {
		v.shards[i][0] = widgets.NewServerWidget("Repl_"+fmt.Sprint(i)+"_PRI", fmt.Sprintf("p%d", i*3+1))
		v.shards[i][1] = widgets.NewServerWidget("Repl_"+fmt.Sprint(i)+"_SEC", fmt.Sprintf("p%d", i*3+2))
		v.shards[i][2] = widgets.NewServerWidget("Repl_"+fmt.Sprint(i)+"_SEC", fmt.Sprintf("p%d", i*3+3))
	}

	v.statusBar = widgets.StatusBar()
}

func (v *ViewMain) Render() {
	termui.Render(termui.Body)
}

func MainView() *ViewMain {
	return _mainView
}

func ResetView() {
	termui.Close()
	err := termui.Init()
	if err != nil {
		panic(err)
	}
}

func init() {
	_mainView = &ViewMain{}
}
