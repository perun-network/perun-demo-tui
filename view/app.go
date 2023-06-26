package view

import (
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"perun.network/perun-demo-tui/asset"
	"perun.network/perun-demo-tui/client"
	"time"
)

const MaxClients = 8

type app struct {
	TUI               *tview.Application
	Mapping           asset.Register
	Clients           []client.DemoClient
	MicroPaymentDelay time.Duration
	left              *tview.Pages
	right             *tview.Pages
}

var App *app

func initApp(clients []client.DemoClient, mapping asset.Register) error {
	if App != nil {
		return errors.New("app already initialized")
	}
	if len(clients) > MaxClients {
		return fmt.Errorf("too many clients. Max: %d, Actual: %d", MaxClients, len(clients))
	}
	App = &app{
		TUI:               tview.NewApplication(),
		Clients:           clients,
		MicroPaymentDelay: 50 * time.Millisecond,
		Mapping:           mapping,
	}
	return nil
}

func setControls() {
	App.TUI.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlA:
			App.TUI.SetFocus(App.left)
		case tcell.KeyCtrlB:
			App.TUI.SetFocus(App.right)
		default:
			switch event.Rune() {
			case 'q':
				App.TUI.Stop()
			}
		}
		return event
	})
}

func setupLayout(title string) {
	header := tview.NewBox().SetBorder(true).SetTitle(" " + title + " ")
	flex := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(header, 3, 0, false).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(App.left, 0, 1, false).
				AddItem(tview.NewBox().SetBorder(true), 2, 0, false).
				AddItem(App.right, 0, 1, false),
			0,
			1,
			false,
		)
	App.TUI.SetRoot(flex, true).EnableMouse(true)
}

func RunDemo(title string, clients []client.DemoClient, mapping asset.Register) error {
	if err := initApp(clients, mapping); err != nil {
		return err
	}
	Left = newView()
	Right = newView()
	setControls()
	App.left = initColumn(Left, "Party A")
	App.right = initColumn(Right, "Party B")
	setupLayout(title)
	return App.TUI.Run()
}
