package template

import (
	"time"

	"github.com/blang/semver/v4"
	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/dwidget"
	"github.com/dReam-dApps/dReams/gnomes"
	"github.com/dReam-dApps/dReams/rpc"
)

// Your dApp name goes here
const app_name = "Template"

// This connection box is global so it can be used in our
// layout and StartApp(), it is not used when importing Template
var connect_box *dwidget.DeroRpcEntries

// A variable to control Gnomon sync status local to Templates requirements
var template_synced bool

var version = semver.MustParse("0.2.0-dev")

// Global Gnomon
var gnomon = gnomes.NewGnomes()

// We can use this func for any initialization required by Template
func initValues() {
	template_synced = false
	logger.Println("Template Initialized")
}

// // Process loop examples for packages

// All loop structures here are the same and in each further
// fetch() example we will add more functionality to the process

// fetch1() is a basic main process routine for Template dApp
func fetch1(d *dreams.AppObject) {
	// Set any initialization values here
	initValues()

	// Wait a moment before we start the loop
	time.Sleep(3 * time.Second)

	// This loop will receive a signal from dReams when it is ready for Template to do some work
	for {
		select {
		case <-d.Receive():
			// If signal is received we should check that wallet and daemon are connected and continue if not connected
			if !rpc.Wallet.IsConnected() || !rpc.Daemon.IsConnected() {
				// If wallet or daemon is not connected we can disable
				// any actions or reset required values here
				Disconnected()

				// Once reset we can signal back to dReams that Template is done working and continue
				d.WorkDone()
				continue
			}

			// If wallet and daemon are connected we can do our main processing here

			// someFuncs()

			// Once done we can signal back to dReams that Template is done working
			d.WorkDone()

		case <-d.CloseDapp():
			// This is the close signal from dReams, do any required close funcs here and return
			logger.Println("[Template] Done")
			return
		}
	}
}

// Building off fetch1(), in this function three additions
// have been made, we preform a initial Gnomon scan, check if
// dReams is viewing Template and send a notification to dReams
func fetch2(d *dreams.AppObject) {
	initValues()
	time.Sleep(3 * time.Second)
	for {
		select {
		case <-d.Receive():
			if !rpc.Wallet.IsConnected() || !rpc.Daemon.IsConnected() {
				Disconnected()
				d.WorkDone()
				continue
			}

			// Here we can run a one time Gnomon scan to sync
			// to certain SCIDS when we are first connected

			// Template will control template_synced, we do not want
			// to preform this scan while dReams is configuring
			if !template_synced && gnomes.GnomonScan(d.IsConfiguring()) {
				// Preform required funcs and set local synced var to true
				logger.Println("[Template] Syncing")

				//  someGnomonFuncs()

				template_synced = true
			}

			// Here we will see if dReams is currently viewing
			// the Template tab and can process as we need accordingly
			if d.OnTab(app_name) {
				logger.Println("[Template] dReams is looking at Template")
			}

			// someFuncs()

			d.WorkDone()
		case <-d.CloseDapp():
			// Here we will send a notification that dReams will display
			d.Notification("From Template", "Hello dReams")
			logger.Println("[Template] Done")
			return
		}
	}
}

// // Functions below should be exported so other applications can use them

// DreamsMenuIntro() can be used to create menu tree items, it can be a single entry
// or you can build a complex tree of information and instructions about your Template
func DreamsMenuIntro() (entries map[string][]string) {
	entries = map[string][]string{
		// This is a single entry
		"Template": {
			"This is a dReams Template",
			"What is a Template?",
			"Have fun"},

		// This makes a branch of "What is dReams Template"
		"What is dReams Template": {
			"Template is....",
			"Template does....",
			"Template FAQs....",
		},
	}

	return
}

// Function for when Template tab is selected in dReams
func OnTabSelected(d *dreams.AppObject) {
	logger.Println("[Template] OnTabSelected()")
}

// Function for when Template tab is first connected
func OnConnected() {
	// soneFuncs()
}

// Function for when Template tab is disconnected
func Disconnected() {
	template_synced = false
}
