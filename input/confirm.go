package input

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func Confirm() {
	// Create a new application
	myApp := app.New()
	myWindow := myApp.NewWindow("Confirmation")

	// // Load the background image
	// img := canvas.NewImageFromFile("path_to_your_image.jpg")
	// img.FillMode = canvas.ImageFillStretch

	// Create Yes and No buttons
	yesButton := widget.NewButton("Yes", func() {
		dialog.ShowInformation("Response", "You clicked Yes!", myWindow)
		myWindow.Close()
	})

	noButton := widget.NewButton("No", func() {
		dialog.ShowInformation("Response", "You clicked No!", myWindow)
		myWindow.Close()
	})

	// Create a horizontal box to hold the buttons
	buttonBox := container.NewHBox(yesButton, noButton)

	// Create a vertical box that contains the image and the buttons
	content := container.NewVBox(buttonBox)

	// Set the content of the window
	myWindow.SetContent(content)

	// Set the window size and focus
	myWindow.Resize(fyne.NewSize(400, 300))
	myWindow.CenterOnScreen()
	myWindow.Show()

	// Make sure the window is always on top
	myWindow.RequestFocus()

	// Run the application
	myApp.Run()
}
