// Update Checker
// Copyright (C) 2020-22  Florian Probst
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

var installedColumn, statusColumn, installedVersionColumn, recentVersionColumn, recentVersionReleaseDateColumn *widget.Box
var a fyne.App
var otherSoftwareText string

func outputResults(installedSoftwareMappings []installedSoftwareMapping) {
	otherSoftwareText = ""
	for _, entry := range installedSoftwareMappings {
		if entry.Status > 1 {
			// write out unknown software only to second windowInfo
			otherSoftwareText += entry.Name
			if entry.InstalledSoftware.DisplayVersion != "" {
				otherSoftwareText += " (Version " + entry.InstalledSoftware.DisplayVersion + ")"
			}
			otherSoftwareText += "\n"
		} else {
			// show on main window
			installedColumn.Append(widget.NewLabel(entry.Name))
			if entry.Status == 0 {
				statusColumn.Append(widget.NewLabelWithStyle(
					"Outdated",
					fyne.TextAlignLeading,
					fyne.TextStyle{Bold: true}))
			} else if entry.Status == 1 {
				statusColumn.Append(widget.NewLabelWithStyle("Up-to-date",
					fyne.TextAlignLeading,
					fyne.TextStyle{Bold: false}))
			}
			if entry.InstalledSoftware.DisplayVersion != "" {
				installedVersionColumn.Append(widget.NewLabel(entry.InstalledSoftware.DisplayVersion))
			} else {
				installedVersionColumn.Append(widget.NewLabel("no version found"))
			}
			if entry.MappedStatus.Version != "" {
				recentVersionColumn.Append(widget.NewLabel(
					entry.MappedStatus.Version))
			} else {
				recentVersionColumn.Append(widget.NewLabel("not available"))
			}

			if entry.MappedStatus.Released != "" {
				recentVersionReleaseDateColumn.Append(widget.NewLabel(
					entry.MappedStatus.Released))
			} else {
				recentVersionReleaseDateColumn.Append(widget.NewLabel("not available"))
			}
		}
	}
}

func createFyneAppWindow() fyne.Window {
	a = app.New()
	a.Settings().SetTheme(theme.LightTheme())

	installedColumn = widget.NewVBox(widget.NewLabelWithStyle(
		"Installed Software                     ",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true}))
	statusColumn = widget.NewVBox(widget.NewLabelWithStyle(
		"Status                       ",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true}))
	installedVersionColumn = widget.NewVBox(widget.NewLabelWithStyle(
		"Installed Version    ",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true}))
	recentVersionColumn = widget.NewVBox(widget.NewLabelWithStyle(
		"Recent Version       ",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true}))
	recentVersionReleaseDateColumn = widget.NewVBox(widget.NewLabelWithStyle(
		"Release Date",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true}))

	appList := widget.NewHBox(
		installedColumn,
		statusColumn,
		installedVersionColumn,
		recentVersionColumn,
		recentVersionReleaseDateColumn,
	)

	top := widget.NewVBox(widget.NewButton("Quit", func() {
		a.Quit()
	}), widget.NewButton("Show Other Software", func() {
		showOtherSoftware()
	}))
	mainContent := widget.NewScrollContainer(fyne.NewContainerWithLayout(layout.NewBorderLayout(top, nil, nil, nil),
		top, widget.NewGroup("Results", appList)))

	mainWindow := a.NewWindow("Update Checker")
	mainWindow.SetContent(mainContent)
	//mainWindow.SetIcon(fyne)
	mainWindow.Resize(fyne.NewSize(1024, 600))

	return mainWindow
}

func showOtherSoftware() {
	// other software window
	otherWin := a.NewWindow("Other installed software")
	grid := widget.NewTextGridFromString(otherSoftwareText)
	//grid.ShowLineNumbers = true
	otherWin.SetContent(widget.NewScrollContainer(fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, nil, nil), grid)))
	otherWin.Resize(fyne.NewSize(1200, 600))
	//otherWin.FullScreen()

	otherWin.Show()
}
