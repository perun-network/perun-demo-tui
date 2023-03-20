package view

import (
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"perun.network/perun-demo-tui/client"
	"polycry.pt/poly-go/sync"
)

type View struct {
	id              uuid.UUID
	Client          client.DemoClient
	Pages           *tview.Pages
	partyAndBalance *tview.TextView
	onStateUpdate   func(string)
	updateLock      sync.Mutex
}

var Left = newView()
var Right = newView()

func newView() *View {
	return &View{
		id:              uuid.New(),
		partyAndBalance: tview.NewTextView().SetTextAlign(tview.AlignCenter).SetDynamicColors(true).SetChangedFunc(func() { App.TUI.Draw() }),
	}
}

func (v *View) UpdateState(s string) {
	v.updateLock.Lock()
	defer v.updateLock.Unlock()
	if v.onStateUpdate != nil {
		v.onStateUpdate(s)
	}
}

func (v *View) UpdateBalance(s string) {
	v.updateLock.Lock()
	defer v.updateLock.Unlock()
	v.partyAndBalance.SetText("[red]" + v.Client.DisplayName() + "[white]: " + v.Client.DisplayAddress() + "\nOn-Chain Balance: " + s)

}

func (v *View) GetID() uuid.UUID {
	return v.id
}

func (v *View) SetClient(c client.DemoClient) {
	if v.Client == c {
		return
	}
	if v.Client != nil {
		v.Client.Deregister(v)
	}
	v.Client = c
	v.Client.Register(v)
}
