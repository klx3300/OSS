package main

import (
	"logger"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	logger.Log.Succln("GUI Initializing...")
	myapp := app.New()
	mywnd := myapp.NewWindow("Hello World!")
	mywnd.SetContent(
		widget.NewVBox(
			widget.NewLabel("Hello!"),
			widget.NewButton("Quit", func() {
				myapp.Quit()
			})))
	mywnd.ShowAndRun()
}
