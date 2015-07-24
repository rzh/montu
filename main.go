package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gizak/termui"
	"github.com/rzh/montu/views"
)

var currentView views.View

func main() {
	// show flashscreen

	views.MainView().Init()
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	currentView = views.MainView().SwitchToView()

	go func() {
		for {
			currentView.Render()

			time.Sleep(1 * time.Second)
		}
	}()

	for {
		e := <-termui.EventCh()

		if e.Type == termui.EventKey {
			if e.Key == termui.KeyEsc {
				break
			} else if e.Key == termui.KeyArrowLeft {
				currentView = views.EventView().SwitchToView()
				currentView.Render()
			} else if e.Key == termui.KeyArrowRight {
				currentView = views.MainView().SwitchToView()
				currentView.Render()
			} else if e.Key == termui.KeyArrowUp {
				currentView = views.RPSView().SwitchToView()
				currentView.Render()
			} else {
				// FIXME: quit with any key
				break
			}
		}
	}
}

func init() {
	fmt.Printf("\033[H\033[2J")
	flash_screen, _ := ioutil.ReadFile("montu.txt")
	fmt.Printf("%s", string(flash_screen[0:len(flash_screen)-2]))

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	fmt.Printf("\033[H\033[2J")
}
