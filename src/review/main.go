//##
// DTrack Package: GUI Review
//
// Splits audio (+video optional) files into 2-second clips and provides a
// GUI tool that moves tagged audio clips into special directories, for traning.
//##
package review

import (
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2"
    "dtrack/state" // Assuming this contains state.Runtime.Record_Inspect_Models
    "sync"
)

// Primary post-bootstrap entry point
func Start() {
	main()
}


// Global status label and mutex for safe cross-thread updates
var (
	statusBarLabel *widget.Label
	statusMutex    sync.Mutex
)

// UpdateStatus sets the text of the global status bar label.
// It's safe to call from any goroutine.
func UpdateStatus(text string) {
	statusMutex.Lock()
	defer statusMutex.Unlock()
	if statusBarLabel != nil {
		statusBarLabel.SetText("Status: " + text)
	}
}

func main() {
	// 1. Setup Application and Window
	a := app.New()
	w := a.NewWindow("dtrack Video Analyzer")

	// 2. Status Bar Setup
	// Initialize the global status label and set an initial message.
	statusBarLabel = widget.NewLabel("Status: Ready")
	statusBarLabel.Alignment = fyne.TextAlignLeading // Ensure text starts from the left

	// Use a container for the status bar to ensure it's visually distinct/full-width
	statusBar := container.NewHBox(
		layout.NewSpacer(), // Push the label to the left
		statusBarLabel,
		layout.NewSpacer(), // Keep it centered if desired, or remove to left-align fully
	)
	// Make sure the label itself is wide enough to show the text

	// 3. Menu Bar (First Row)
	selectVideoButton := widget.NewButton("Select Video", func() {
		UpdateStatus("Video selection dialog opened.")
		// Placeholder for video selection logic
	})

	menuBarButtons := []fyne.CanvasObject{selectVideoButton}

	// Add buttons based on the external state
	for _, button_name := range state.Runtime.Record_Inspect_Models {
		button_name_copy := button_name // Capture range variable
		btn := widget.NewButton(button_name_copy, func() {
			UpdateStatus("Clicked model button: " + button_name_copy)
			// Placeholder for model action
		})
		menuBarButtons = append(menuBarButtons, btn)
	}

	// Use an HBox for the menu_bar, which arranges items horizontally.
	menuBar := container.NewHBox(menuBarButtons...)

	// 4. Clip View (Second Row)

	// Left: Listbox (List of Strings)
	data := []string{"Clip 1: Start 00:00:00", "Clip 2: Start 00:00:15", "Clip 3: Start 00:00:30"}
	clipList := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template") // Template object
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i]) // Data binding
		},
	)

	clipList.OnSelected = func(id widget.ListItemID) {
		UpdateStatus("Selected clip: " + data[id])
	}

	// Right: Image Space
	// The space for the image. We use a placeholder rectangle for visual representation,
	// but an *image.Image would be used here in a real app.
	// We use container.NewMax to make it fill the available space.
	imagePlaceholder := widget.NewLabel("Video Frame Image Will Go Here")
	imagePlaceholder.Alignment = fyne.TextAlignCenter
	imagePlaceholder.TextStyle.Italic = true

	// Wrap the list in a container to give it a fixed width or size it explicitly.
	// A good starting width might be a fraction of the screen or a fixed pixel value.
	// For this example, we'll let the HSplit handle the relative sizing.
	listContainer := container.NewVScroll(clipList)
	listContainer.SetMinSize(fyne.NewSize(250, 0)) // Example: minimum width of 250 pixels

	// The HSplit is key to achieving the desired "weight" distribution (like tkinter grid weights).
	// The list on the left and the image space on the right.
	// HSplit allows the user to resize the split, and the second object (image) will take up
	// the maximum available space by default.
	clipView := container.NewHSplit(
		listContainer,
		imagePlaceholder, // The image space, which will grow to fill space
	)
	// Set the initial ratio to favor the image space (e.g., 0.2 split means 20% for the list)
	clipView.SetOffset(0.2)

	// 5. Main Content Layout (All Rows)
	// The main layout uses a Border layout to place the status bar at the bottom (SOUTH)
	// and the rest of the content (menuBar and clipView) in the center.

	// The center content combines the menuBar (top) and clipView (middle) using a VBox.
	// VBox arranges items vertically and allocates proportional space.
	centerContent := container.NewVBox(
		menuBar,
		// The clipView is wrapped in a container.NewMax to ensure it uses the maximum
		// available vertical space remaining after the menuBar is placed.
		container.NewMax(clipView),
	)

	// Border layout for the main window content
	content := container.NewBorder(
		nil,           // TOP: Nothing
		statusBar,     // BOTTOM: Status bar
		nil,           // LEFT: Nothing
		nil,           // RIGHT: Nothing
		centerContent, // CENTER: Menu bar and Clip view
	)

	// 6. Finalize and Run
	w.SetContent(content)
	w.Resize(fyne.NewSize(1024, 768)) // Set a reasonable default size
	w.ShowAndRun()
}

// NOTE: You would need to ensure "dtrack/state" is correctly implemented
// and contains the structure used: `state.Runtime.Record_Inspect_Models`
// (which must be a slice of strings or similar iterable).
