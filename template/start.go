package template

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/civilware/Gnomon/structures"
	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/gnomes"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"
	"github.com/sirupsen/logrus"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

// // dReams Template dApp

var logger = structures.Logger.WithFields(logrus.Fields{})

// Run Template as a single dApp, this func is used in cmd/main.go
func StartApp() {
	// Set max cpu
	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)

	// Initialize logger to Stdout
	gnomes.InitLogrusLog(logrus.InfoLevel)

	// Read dReams config file
	config := menu.ReadDreamsConfig(app_name)

	// Initialize Fyne app and set Templates base theme from 'config.Skin'
	a := app.New()
	a.Settings().SetTheme(bundle.DeroTheme(config.Skin))

	// Initialize Fyne window with icon from dReams
	// bundle package and size our window
	w := a.NewWindow("dReams Template")
	w.SetIcon(bundle.ResourceBlueBadgePng)
	w.Resize(fyne.NewSize(1400, 800))
	w.SetMaster()

	// Channel used for closing
	done := make(chan struct{})

	// Here we make a dreams.AppObject for our Template,
	// initializing Background that is a stack (max) container with a canvas.Image,
	// Window as our Fyne window 'w' and App as Fyne App 'a'
	menu.Theme.Img = *canvas.NewImageFromResource(nil)
	d := dreams.AppObject{
		App:        a,
		Window:     w,
		Background: container.NewStack(&menu.Theme.Img),
	}

	// Set what we'd like to happen on close.
	// Here we are saving dReams config file with Skin,
	// closing Gnomon, our main process and then Fyne window
	closeFunc := func() {
		menu.WriteDreamsConfig(
			dreams.SaveData{
				Skin:   bundle.AppColor,
				Daemon: []string{rpc.Daemon.Rpc},
				DBtype: gnomon.DBStorageType(),
			})
		gnomon.Stop(app_name)
		d.StopProcess()
		w.Close()
	}

	// Assign closeFunc to window close
	w.SetCloseIntercept(closeFunc)

	// This channel will receive CTRL+C signal to close app
	// and do the same as the above SetCloseIntercept()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println()
		closeFunc()
	}()

	// Initialize Gnomon fast sync to true so we can quickly create a DB to use
	gnomon.SetFastsync(true)

	// Here we start the main process, if Template was being imported for use this would be
	// the main process of the parent app. It will lead then call Templates fetch() and close it when exiting.
	// As you build your Template, additions can be made either here or in fetch(), keeping parent app
	// functions in mind if used as a import so to avoid redundant calls.
	go func() {
		// Delay the routine to give time for app to start
		time.Sleep(3 * time.Second)

		// dReams runs on 3 second base tick
		ticker := time.NewTicker(3 * time.Second)

		// Initialize dReams token balance maps
		rpc.InitBalances()

		for {
			select {
			// do this on our 3 second interval
			case <-ticker.C:
				// Ping our daemon and wallet for connection
				rpc.Ping()
				rpc.EchoWallet(app_name)

				// Get any balances required
				rpc.Wallet.GetBalance()

				// // To get a single tokens balance
				// rpc.Wallet.GetTokenBalance("token_name", TOKEN_SCID)

				// // To get all supported dReams token balances
				// rpc.GetDreamsBalances(rpc.SCIDs)

				// Refresh Dero balance in UI display
				connect_box.RefreshBalance()

				// Update Gnomon end point if daemon has changed
				gnomes.GnomonEndPoint()

				// If daemon is connected and Gnomon is initialized we will set the hidden Disconnect
				// object in the connect_box to true, this Disconnect will control Gnomon shut down
				if rpc.Daemon.IsConnected() && gnomon.IsInitialized() {
					connect_box.Disconnect.SetChecked(true)

					// If Gnomon is running we can start to do some checks
					if gnomon.IsRunning() {
						// This will populate the Gnomes.SCID count var
						contracts := gnomon.IndexContains()

						// Refresh Gnomon and daemon height displays
						gnomonLabel.SetText(fmt.Sprintf("Gnomon Height: %d", gnomon.GetLastHeight()))
						indexLabel.SetText(fmt.Sprintf("Indexed SCIDs: %d", len(contracts)))
						daemonLabel.SetText(fmt.Sprintf("Daemon Height: %d", gnomon.GetChainHeight()))

						// This will set the Gnomes.Check value to
						// true once 100 or more contracts have been indexed
						if gnomon.HasIndex(100) {
							gnomon.Checked(true)
						}
					}

					// Here we can use some of the Indexer vars to set our Synced status
					if gnomon.GetLastHeight() >= gnomon.GetChainHeight()-3 {
						gnomon.Synced(true)
					} else {
						// If Template is not synced we should handle that here
						gnomon.Synced(false)
						gnomon.Checked(false)
					}
				} else {
					// If daemon is not connected we will set Disconnect object to false
					connect_box.Disconnect.SetChecked(false)
				}

				// Signal to LayoutAllItems() that we are ready for it to do some work
				d.SignalChannel()

				// Exit Templates main process
			case <-d.Closing():
				logger.Println("[Template] Closing...")

				// Stop Gnomon indicator if it exists
				if gnomes.Indicator.Icon != nil {
					gnomes.Indicator.Icon.Stop()
				}

				// Stop Templates ticker
				ticker.Stop()

				// Stop running dApp routines
				d.CloseAllDapps()
				time.Sleep(time.Second)

				// Send done signal once all routines are closed before returning
				done <- struct{}{}
				return
			}
		}
	}()

	// LayoutAllItems() has one routine of fetch1(), we can set a channel for it in dreams.AppObject
	d.SetChannels(1)

	// Set Templates content as a stack (max) container with our 'd.Background'
	// first, followed by Templates LayoutAllItems(), finally a VBox
	// that will contain the required components for Dero RPC
	// connection and status indicators, see connectBox() for more info
	go func() {
		// using a delayed routine here to allow our window to run for a moment before placing layout
		time.Sleep(450 * time.Millisecond)
		w.SetContent(container.NewStack(d.Background, LayoutAllItems(true, &d), container.NewVBox(layout.NewSpacer(), connectBox())))
	}()

	// Start Template dApp
	w.ShowAndRun()

	// We can use this channel for ensuring any closures after main window has closed.
	<-done
	logger.Printf("[%s] Closed\n", app_name)
}
