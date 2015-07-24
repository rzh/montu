package widgets

import (
	"fmt"
	"strings"
	"time"

	"github.com/gizak/termui"
	"github.com/rzh/montu/gen"
)

const LINE = 8

type MigrationWidget struct {
	msg []string
	bc  *termui.Par
	G   *gen.ChunkMigrationGen
}

var _migration *MigrationWidget

func Migration() *MigrationWidget {
	return _migration
}

func (w *MigrationWidget) Init() {
	w.G = gen.ChunkMigration()

	par3 := termui.NewPar("")
	par3.Height = LINE + 2
	par3.Width = 137
	par3.Y = 9
	par3.Border.Label = "> Migrations <"

	par3.Text = ""

	w.bc = par3

	go func() {
		for {
			time.Sleep(1 * time.Second)

			// request data
			w.G.C <- [3]int64{1, 1, 1}
			a := <-w.G.C
			m := fmt.Sprintf(" %d: %d --> %d", a[0], a[1], a[2])

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

func (w *MigrationWidget) Widget() termui.GridBufferer {
	return w.bc
}

func init() {
	_migration = &MigrationWidget{}
	_migration.Init()
}
