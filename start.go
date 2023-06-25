package template

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// // dReams Template dApp

// Run Template as a single dApp, this func is used in cmd/main.go
func StartApp() {
	// Set max cpu
	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)

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

	// Channels used for closing
	quit := make(chan struct{})
	done := make(chan struct{})

	// Set what we'd like to happen on close.
	// Here we are saving dReams config file with Skin,
	// closing Gnomon, our main process and then Fyne window
	w.SetCloseIntercept(func() {
		menu.WriteDreamsConfig(
			dreams.DreamSave{
				Skin:   bundle.AppColor,
				Daemon: []string{rpc.Daemon.Rpc},
			})
		menu.Gnomes.Stop(app_name)
		quit <- struct{}{}
		w.Close()
	})

	// This channel will receive CTRL+C signal to close app
	// and do the same as the above SetCloseIntercept()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println()
		menu.WriteDreamsConfig(
			dreams.DreamSave{
				Skin:   bundle.AppColor,
				Daemon: []string{rpc.Daemon.Rpc},
			})
		menu.Gnomes.Stop(app_name)
		quit <- struct{}{}
		w.Close()
	}()

	// Initialize Gnomon fast sync to true so we can quickly create a DB to use
	menu.Gnomes.Fast = true

	// Here we make a dreams.DreamsObject for our Template,
	// initializing Background that is a max container with a canvas.Image
	// and Window as our Fyne window 'w'
	dreams.Theme.Img = *canvas.NewImageFromResource(nil)
	d := dreams.DreamsObject{
		Window:     w,
		Background: container.NewMax(&dreams.Theme.Img),
	}

	// Here we start Templates main process, as you build your Template
	// additions may be needed to this process to meet to your requirements
	go func() {
		// Print out some info about our Template
		log.Printf("[%s] %s %s %s\n", app_name, rpc.DREAMSv, runtime.GOOS, runtime.GOARCH)

		// Delay the routine to give time for app to start
		time.Sleep(6 * time.Second)

		// dReams runs on 3 second base tick, we can use any offset required
		var offset int
		ticker := time.NewTicker(3 * time.Second)

		// Initialize a token balance map
		rpc.Wallet.TokenBal = make(map[string]uint64)

		for {
			select {
			// do this on our 3 second interval
			case <-ticker.C:
				// Ping our daemon and wallet for connection
				rpc.Ping()
				rpc.EchoWallet(app_name)

				// Get any balances required
				rpc.Wallet.GetBalance()
				//rpc.Wallet.GetTokenBalance("token_name", TOKEN_SCID)

				// Refresh Dero balance in UI display
				connect_box.RefreshBalance()

				// Update Gnomon end point if daemon has changed
				menu.GnomonEndPoint()

				// If daemon is connected and Gnomon is initialized we will set the hidden Disconnect
				// object in the connect_box to true, this Disconnect will control Gnomon shut down
				if rpc.Daemon.IsConnected() && menu.Gnomes.IsInitialized() {
					connect_box.Disconnect.SetChecked(true)

					// If Gnomon is running we can start to do some checks
					if menu.Gnomes.IsRunning() {
						// This will populate the Gnomes.SCID count var
						menu.Gnomes.IndexContains()

						// This will set the Gnomes.Check value to
						// true once 100 or more contracts have been indexed
						if menu.Gnomes.HasIndex(100) {
							menu.Gnomes.Checked(true)
						}
					}

					// Here we can use some of the Indexer vars to set our Synced status
					if menu.Gnomes.Indexer.LastIndexedHeight >= menu.Gnomes.Indexer.ChainHeight-3 {
						menu.Gnomes.Synced(true)
					} else {
						// If Template is not synced we should handle that here
						menu.Gnomes.Synced(false)
						menu.Gnomes.Checked(false)
					}
				} else {
					// If daemon is not connected we will set Disconnect object to false
					connect_box.Disconnect.SetChecked(false)
				}

				// Here we can use our offset to call a function once every 30 seconds
				offset++
				if offset%10 == 0 {
					log.Println("[Template] Offset called here")
					offset = 0
				}

				// Exit Templates main process
			case <-quit:
				log.Println("[Template] Closing...")

				// Stop Gnomon indicator if it exists
				if menu.Gnomes.Icon_ind != nil {
					menu.Gnomes.Icon_ind.Stop()
				}

				// Stop Templates ticker
				ticker.Stop()

				// Close out all other routines and send done signal before returning
				done <- struct{}{}
				return
			}
		}
	}()

	// Set Templates content as a max container with our 'd.Background' first
	// followed by Templates LayoutAllItems(), using a delayed routine here to allow
	// our window to run for a moment before placing layout
	go func() {
		time.Sleep(450 * time.Millisecond)
		w.SetContent(container.NewMax(d.Background, LayoutAllItems(false, d)))
	}()

	// Start Template dApp
	w.ShowAndRun()

	// We can use this channel for ensuring any closures after main window has closed.
	<-done
	log.Printf("[%s] Closed\n", app_name)
}
