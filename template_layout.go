package template

import (
	"image/color"
	"log"

	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/dwidget"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// The premise of how dReams imports packages is for the package
// to be contained within a LayoutAllItems() which returns a
// Fyne max container and starts any required routines. Using this
// Template package you can create local dApps that will run independently
// while also being able to be imported for use in other Go/Fyne applications.
// dReams has structs like dream.DreamsItems that can help in organizing content
// although you are free to set your layout up in any manner you please
// see developer.fyne.io for more info on Fyne layouts

// // Simple LayoutAllItems() example for package dApps

// Entire Template dApp layout is in this func, returned as fyne.CanvasObject
// it can be handled differently depending on if it is imported or not,
// 'd' is the dReams app object which Templates routines will get signals and checks from
func LayoutAllItems(imported bool, d dreams.DreamsObject) (max fyne.CanvasObject) {
	// A couple widgets to get started
	label := container.NewCenter(widget.NewLabel("This is a label"))

	// Text entry
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Text Entry:")

	// Button that will print out the current text of entry when pressed
	button := widget.NewButton("Button", func() {
		log.Printf("[Template] Text entry contains {%s}\n", entry.Text)
	})

	// Radio widget that will change Templates skin and connect_box Balance label color
	radio := widget.NewRadioGroup([]string{"Dark", "Light"}, func(s string) {
		switch s {
		case "Dark":
			// We can tie into the current AppColor of dReams with bundle.AppColor
			// when Templates closes it wil save the AppColor
			bundle.AppColor = color.Black

		case "Light":
			bundle.AppColor = color.White
		default:

			return
		}

		// bundle.DeroTheme() has a light and dark theme
		fyne.CurrentApp().Settings().SetTheme(bundle.DeroTheme(bundle.AppColor))

		// We can tie into the current text color of dReams with bundle.TextColor
		connect_box.Balance.Color = bundle.TextColor
		connect_box.Balance.Refresh()
	})

	// A container for our widgets
	tab1_cont := container.NewBorder(label, button, nil, nil, entry)

	// Another container for widgets on a different tab
	tab2_cont := container.NewMax(container.NewCenter(radio))

	// These are the tabs we want in our Template
	// First tab is some widgets with a dynamic alpha layer behind it
	// Second is the radio widget to change Template skin
	// Third is a UI log which can be used to record session TXs and info
	tabs := container.NewAppTabs(
		container.NewTabItem("Tab1", container.NewMax(bundle.NewAlpha120(), tab1_cont)),
		container.NewTabItem("Tab2", tab2_cont),
		container.NewTabItem("Log", rpc.SessionLog()))

	// What will happen when tabs are selected locally
	tabs.OnSelected = func(ti *container.TabItem) {
		switch ti.Text {
		case "Tab1":
			log.Println("[Template] Tab1 Selected")
		case "Tab2":
			log.Println("[Template] Tab2 Selected")
		case "Log":
			log.Println("[Template] Log Selected")
		default:

		}
	}

	// Local tabs should be placed at bottom of Template
	tabs.SetTabLocation(container.TabLocationBottom)

	// If we are importing this Template we can place our
	// tabs in final max container and start required routines passing in 'd'
	if imported {
		max = container.NewMax(tabs)
		go fetch1(d)
	} else {
		// If Template is running independently, we can add a VBox
		// that will contain the required components for Dero RPC
		// connection and status indicators, see connectBox() below for more info
		max = container.NewMax(tabs, container.NewVBox(layout.NewSpacer(), connectBox()))
	}

	return
}

// This is a construction of dwidget.DeroRpcEntries which is used to
// run Template independently, connect_box has been declared globally
// to the package for easy of use. It will connect to Dero wallet and
// daemon RPC, and start Gnomon if connected, indexing all SCIDS
func connectBox() *fyne.Container {
	// Initialize connect_box on trailing edge to fit our tabs in LayoutAllItems()
	connect_box = dwidget.HorizontalEntries(app_name, 1)

	// Set what we'd like to occur when button is pressed
	connect_box.Button.OnTapped = func() {
		// Get Dero wallet address
		rpc.GetAddress(app_name)

		// Ping daemon
		rpc.Ping()

		// Here we are starting Gnomon without a search filter to index all SCIDs
		if rpc.Daemon.IsConnected() && !menu.Gnomes.IsInitialized() && !menu.Gnomes.Start {
			go menu.StartGnomon(app_name, "boltdb", []string{}, 0, 0, nil)
		}
	}

	// dwidget.DeroRpcEntries have a hidden check which we
	// will use here to shut down Gnomon on daemon disconnection
	connect_box.Disconnect.OnChanged = func(b bool) {
		if !b {
			menu.Gnomes.Stop(app_name)
		}
	}

	// Read dReams config file and set saved daemon option
	config := menu.ReadDreamsConfig(app_name)
	connect_box.AddDaemonOptions(config.Daemon)

	// Adding dReams indicator panel for wallet, daemon and Gnomon
	connect_box.Container.Objects[0].(*fyne.Container).Add(menu.StartIndicators())

	return connect_box.Container
}
