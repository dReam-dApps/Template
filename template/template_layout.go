package template

import (
	"image/color"

	"github.com/dReam-dApps/dImports/dimport"
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
// Fyne stack (max) container and starts any required routines. Using this
// Template package you can create local dApps that will run independently
// while also being able to be imported for use in other Go/Fyne applications.
// dReams has structs like dream.ContainerStack that can help in organizing content
// although you are free to set your layout up in any manner you please
// see developer.fyne.io for more info on Fyne layouts

// // Simple LayoutAllItems() example for package dApps

// Global display labels
var gnomonLabel = widget.NewLabel("Gnomon Height:")
var daemonLabel = widget.NewLabel("Daemon Height:")
var indexLabel = widget.NewLabel("Indexed SCIDs:")

// Entire Template dApp layout is in this func, returned as fyne.CanvasObject
// it can be handled differently depending on if it is imported or not,
// 'd' is the dReams app object which Templates routines will get signals and checks from
func LayoutAllItems(imported bool, d *dreams.AppObject) fyne.CanvasObject {
	// Place the global labels into containers
	label := container.NewCenter(container.NewHBox(gnomonLabel, indexLabel, daemonLabel))

	// Radio widget that will change Templates skin and connect_box Balance label color
	radio := widget.NewRadioGroup([]string{"Dark", "Light"}, nil)
	if bundle.AppColor == color.White {
		radio.SetSelected("Light")
	} else {
		radio.SetSelected("Dark")
	}
	radio.Horizontal = true
	radio.OnChanged = func(s string) {
		switch s {
		case "Dark":
			// We can tie into the current AppColor of dReams with bundle.AppColor
			// when Templates closes it will save the AppColor
			bundle.AppColor = color.Black
		case "Light":
			bundle.AppColor = color.White
		default:

		}

		go func() {
			// bundle.DeroTheme() has a light and dark theme
			fyne.CurrentApp().Settings().SetTheme(bundle.DeroTheme(bundle.AppColor))

			// Reset all widgets with new skin, this will disconnect wallet
			d.Window.Content().Refresh()

			// We can tie into the current text color of dReams with bundle.TextColor
			connect_box.Balance.Color = bundle.TextColor
			connect_box.Balance.Refresh()
		}()
	}

	// A container for our widgets
	tab1_cont := container.NewBorder(label, container.NewCenter(radio), nil, nil)

	// Another container for widgets on a different tab,
	// ImportWidget() can import and run Go packages
	// see "github.com/dReam-dApps/dImports/dimport" for details
	tab2_cont := container.NewStack(container.NewCenter(container.NewAdaptiveGrid(3, layout.NewSpacer(), dimport.ImportWidget(d))))

	//// Tab 3 start here
	rpc_entries := dwidget.NewHorizontalEntries("", 1)
	rpc_entries.Button.OnTapped = func() {
		rpc.GetAddress(app_name)
		rpc.Ping()
	}

	rpc_entries.AddIndicator(menu.StartIndicators())

	message_button := widget.NewButton("Message", func() {
		menu.SendMessageMenu("", nil)
	})

	dest_entry := widget.NewEntry()
	dest_entry.SetPlaceHolder("Address:")

	asset_button := widget.NewButton("Asset", func() {
		rpc.SendAsset(rpc.DreamsSCID, dest_entry.Text, true)
	})

	tab3_cont := container.NewBorder(rpc_entries.Container, dest_entry, asset_button, message_button)

	// These are the tabs we want in our Template
	// First tab is labels and radio widget with a dynamic alpha layer behind it
	// Second is a empty tab
	// Third is a UI log which can be used to record session TXs and info
	tabs := container.NewAppTabs(
		container.NewTabItem("Tab1", container.NewStack(bundle.NewAlpha120(), tab1_cont)),
		container.NewTabItem("Tab2", tab2_cont),
		container.NewTabItem("Tab3", tab3_cont),
		container.NewTabItem("Log", rpc.SessionLog()))

	//// Workshop address here
	///  dero1qyr725edhmd5lqrg75y56guj58cldv2fsau49ee2n7f0cdkhy2fkgqq4s06km

	// What will happen when tabs are selected locally
	tabs.OnSelected = func(ti *container.TabItem) {
		switch ti.Text {
		case "Tab1":
			logger.Println("[Template] Tab1 Selected")
		case "Tab2":
			logger.Println("[Template] Tab2 Selected")
		default:

		}
	}

	// Local tabs should be placed at bottom of Template
	tabs.SetTabLocation(container.TabLocationBottom)

	// If we are importing this Template we can start our
	// packages routine passing in 'd'
	if imported {
		go fetch1(d)
	} else {
		// If Template is running independently, rpc connection or the layout
		// could be altered here if required. For this Template the
		// required rpc objects are placed inside StartApp()
	}

	return container.NewStack(tabs)
}

// This is a construction of dwidget.DeroRpcEntries which is used to
// run Template independently, connect_box has been declared globally
// to the package for ease of use. It will connect to Dero wallet and
// daemon RPC, and start Gnomon if connected, indexing all SCIDS
func connectBox() *fyne.Container {
	// Initialize connect_box on trailing edge to fit our tabs in LayoutAllItems()
	connect_box = dwidget.NewHorizontalEntries(app_name, 1)

	// Set what we'd like to occur when button is pressed
	connect_box.Button.OnTapped = func() {
		// Get Dero wallet address
		rpc.GetAddress(app_name)

		// Ping daemon
		rpc.Ping()

		// Here we are starting Gnomon without a search filter to index all SCIDs
		if rpc.Daemon.IsConnected() && !menu.Gnomes.IsInitialized() && !menu.Gnomes.Start {
			go menu.StartGnomon(app_name, menu.Gnomes.DBType, []string{}, 0, 0, nil)
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
	connect_box.AddIndicator(menu.StartIndicators())

	return connect_box.Container
}
