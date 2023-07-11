package template

import (
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	dreams "github.com/dReam-dApps/dReams"
	"github.com/x-motemen/gore"
)

// Using gore package to test importing and running a packages StartApp() at run time

// Import a go package and call its pkg.StartApp()
func importPackage(path string, d *dreams.AppObject) {
	var stderr strings.Builder
	s, err := gore.NewSession(os.Stdin, &stderr)
	if err != nil {
		logger.Errorln("[Template]", err)
		return
	}

	logger.Println("[Template] Importing:", path)
	if err = s.Eval(fmt.Sprintf(":import %s", path)); err != nil {
		dialog.NewError(err, d.Window).Show()
		logger.Errorln("[Template]", err)
		s.Clear()
		return
	}

	split := strings.Split(path, "/")
	l := len(split)

	if l < 4 || split[l-1] == "" {
		dialog.NewError(fmt.Errorf("invalid package path %s", path), d.Window).Show()
		logger.Errorf("[Template] Invalid package path %s", path)
		s.Clear()
		return
	}

	start_cmd := fmt.Sprintf("%s.StartApp()", split[l-1])
	if err = s.Eval(start_cmd); err != nil {
		dialog.NewError(fmt.Errorf("command %s failed", start_cmd), d.Window).Show()
		logger.Errorf("[Template] Import start command %s failed", start_cmd)
	}

	s.Clear()
}

// Widgets for importPackage()
func importWidget(d *dreams.AppObject) fyne.CanvasObject {
	label := widget.NewLabel("Import and run a go packages StartApp()")
	label.Alignment = fyne.TextAlignCenter

	path_entry := widget.NewEntry()
	path_entry.SetPlaceHolder("Go import path:")

	loading := widget.NewProgressBarInfinite()
	loading.Hide()

	import_button := widget.NewButton("Import", nil)
	import_button.OnTapped = func() {
		go func() {
			import_button.Hide()
			loading.Start()
			loading.Show()
			importPackage(path_entry.Text, d)
			loading.Stop()
			loading.Hide()
			import_button.Show()
		}()
	}

	return container.NewVBox(
		layout.NewSpacer(),
		label,
		path_entry,
		loading,
		import_button,
		layout.NewSpacer())
}
