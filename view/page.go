package view

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"perun.network/perun-demo-tui/client"
	"strconv"
	"time"
)

const (
	PartySelectionPage = "PartySelectionPage"
	PartyMenuPage      = "PartyMenuPage"
	OpenChannelPage    = "OpenChannelPage"
	DisplayChannelPage = "DisplayChannelPage"
)

func initColumn(view *View, party string) *tview.Pages {
	pages := tview.NewPages()
	pages.SetFocusFunc(func() {
		_, frontPage := pages.GetFrontPage()
		App.TUI.SetFocus(frontPage)
	})
	pages.AddPage(PartySelectionPage, newPartiesPage(party, view), true, true)
	pages.AddPage(PartyMenuPage, newPartyMenuPage(view), true, false)
	displayChannelPage, channelUpdateHandler := newDisplayChannelPage(view)
	view.onStateUpdate = channelUpdateHandler
	pages.AddPage(DisplayChannelPage, displayChannelPage, true, false)
	pages.AddPage(OpenChannelPage, newOpenChannelPage(view), true, false)
	pages.SwitchToPage(PartySelectionPage)
	return pages
}

func setClientAndSwitchToPartyMenuPage(client client.DemoClient, view *View) func() {
	return func() {
		view.SetClient(client)
		log.Println("Switching to PartyMenuPage")
		view.Pages.SwitchToPage(PartyMenuPage)
	}
}

var digitRunes = []rune("0123456789")

func newPartiesPage(title string, view *View) tview.Primitive {
	content := tview.NewFlex().SetDirection(tview.FlexRow)
	content.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetDynamicColors(true).SetText("[red]"+title), 2, 0, false)
	list := tview.NewList()

	for i, c := range App.Clients {
		if i > MaxClients {
			break
		}

		list.AddItem(c.DisplayName(), "Address: "+c.DisplayAddress(), digitRunes[i+1], setClientAndSwitchToPartyMenuPage(c, view))
	}

	content.AddItem(list, 0, 1, true)
	list.SetSelectedFocusOnly(true)
	return content
}

func newPartyMenuPage(view *View) tview.Primitive {
	content := tview.NewFlex().SetDirection(tview.FlexRow)
	content.AddItem(view.partyAndBalance, 2, 0, false)
	header := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Menu")
	content.AddItem(header, 2, 0, false)

	list := tview.NewList().SetSelectedFocusOnly(true)
	list.AddItem("Open Channel", "Open a new Channel with another party", 'o', func() {
		view.Pages.SwitchToPage(OpenChannelPage)
	})
	list.AddItem("View Channels", "View open channel", 'v', func() {
		view.Pages.SwitchToPage(DisplayChannelPage)
	})
	content.AddItem(list, 0, 1, true)

	content.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'r':
			view.Pages.SwitchToPage(PartySelectionPage)
		}
		return event
	})
	return content
}

func newDisplayChannelPage(view *View) (tview.Primitive, func(string)) {
	content := tview.NewFlex().SetDirection(tview.FlexRow)
	content.AddItem(view.partyAndBalance, 2, 0, false)
	content.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("View Channel"), 2, 0, false)
	channelView := tview.NewTextView().SetDynamicColors(true).SetChangedFunc(func() { App.TUI.Draw() })
	sendForm := tview.NewForm()
	channelView.SetText("Currently no open channel for this client")
	content.AddItem(channelView, 0, 1, true)
	channelView.SetFocusFunc(func() {
		App.TUI.SetFocus(sendForm)

	})
	content.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'r':
			view.Pages.SwitchToPage(PartyMenuPage)
		}
		return event
	})
	content.AddItem(sendForm, 0, 1, false)
	setForm := func() {
		sendField := tview.NewInputField().SetLabel("Send Payment").SetFieldWidth(20).SetText("")
		microPaymentAmount := tview.NewInputField().SetLabel("Micro Payment Amount").SetFieldWidth(20).SetText("")
		microPaymentRepetitions := tview.NewInputField().SetLabel("Repetitions").SetFieldWidth(20).SetText("")
		*sendForm = *tview.NewForm().AddFormItem(sendField).
			AddButton("Send", func() {
				amount, err := strconv.ParseFloat(sendField.GetText(), 64)
				if err != nil {
					return
				}
				go view.Client.SendPaymentToPeer(amount)
			}).
			AddFormItem(microPaymentAmount).AddFormItem(microPaymentRepetitions).
			AddButton("Send Micro Payment", func() {
				amount, err := strconv.ParseFloat(microPaymentAmount.GetText(), 64)
				if err != nil {
					return
				}
				repetitions, err := strconv.ParseInt(microPaymentRepetitions.GetText(), 10, 64)
				if err != nil {
					return
				}
				go func() {
					for i := int64(0); i < repetitions; i++ {
						view.Client.SendPaymentToPeer(amount)
						time.Sleep(App.MicroPaymentDelay)
					}
				}()
			}).
			AddButton("Settle", func() {
				go view.Client.Settle()
			})
	}
	sendForm.SetFocusFunc(func() {
		if view.Client.HasOpenChannel() {
			setForm()
		}
	})

	return content, func(s string) {
		channelView.SetText(s)
		if view.Client.HasOpenChannel() {
			setForm()
		}
	}
}

func newOpenChannelPage(view *View) tview.Primitive {
	content := tview.NewFlex().SetDirection(tview.FlexRow)
	content.AddItem(view.partyAndBalance, 2, 0, false)
	content.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Open Channel"), 2, 0, false)
	form := tview.NewForm()
	content.AddItem(form, 0, 1, false)
	content.SetFocusFunc(func() {
		clientSelection := make(map[int]client.DemoClient)
		var clientNames []string
		i := 0
		for _, c := range App.Clients {
			if c == view.Client {
				continue
			}
			clientSelection[i] = c
			str := fmt.Sprintf("%s (%s)", c.DisplayName(), c.DisplayAddress())
			clientNames = append(clientNames, str)
			i++
		}
		peerField := tview.NewDropDown().SetLabel("Party").SetOptions(clientNames, nil).SetCurrentOption(0)
		depositField := tview.NewInputField().SetLabel("Deposit").SetFieldWidth(20).SetText("")
		*form = *tview.NewForm().AddFormItem(peerField).AddFormItem(depositField).
			AddButton("Open Channel", func() {
				deposit, err := strconv.ParseFloat(depositField.GetText(), 64)
				if err != nil {
					return
				}
				peerIndex, _ := peerField.GetCurrentOption()
				peer := clientSelection[peerIndex]
				go view.Client.OpenChannel(peer.WireAddress(), deposit)
				view.Pages.SwitchToPage(DisplayChannelPage)
			}).
			AddButton("Cancel", func() {
				view.Pages.SwitchToPage(PartyMenuPage)
			})
		App.TUI.SetFocus(form)
	})
	content.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'r':
			view.Pages.SwitchToPage(PartyMenuPage)
		}
		return event
	})
	var clientNames []string
	for _, c := range App.Clients {
		if c == view.Client {
			continue
		}
		str := fmt.Sprintf("%s (%s)", c.DisplayName(), c.DisplayAddress())
		clientNames = append(clientNames, str)
	}
	return content
}
