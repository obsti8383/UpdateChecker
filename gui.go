// Update Checker
// Copyright (C) 2020  Florian Probst
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

var InstalledColumn, StatusColumn, InstalledVersionColumn, RecentVersionColumn, RecentVersionReleaseDateColumn *widget.Box
var a fyne.App
var OtherSoftwareText string

func outputResults(installedSoftwareMappings []installedSoftwareMapping) {
	OtherSoftwareText = ""
	for _, entry := range installedSoftwareMappings {
		if entry.Status > 1 {
			// write out unknown software only to second windowInfo
			OtherSoftwareText += entry.Name
			if entry.InstalledSoftware.DisplayVersion != "" {
				OtherSoftwareText += " (Version " + entry.InstalledSoftware.DisplayVersion + ")"
			}
			OtherSoftwareText += "\n"
		} else {
			// show on main window
			InstalledColumn.Append(widget.NewLabel(entry.Name))
			if entry.Status == 0 {
				StatusColumn.Append(widget.NewLabelWithStyle(
					"Outdated",
					fyne.TextAlignLeading,
					fyne.TextStyle{Bold: true}))
			} else if entry.Status == 1 {
				StatusColumn.Append(widget.NewLabelWithStyle("Up-to-date",
					fyne.TextAlignLeading,
					fyne.TextStyle{Bold: false}))
			}
			if entry.InstalledSoftware.DisplayVersion != "" {
				InstalledVersionColumn.Append(widget.NewLabel(entry.InstalledSoftware.DisplayVersion))
			} else {
				InstalledVersionColumn.Append(widget.NewLabel("no version found"))
			}
			if entry.MappedStatus.Version != "" {
				RecentVersionColumn.Append(widget.NewLabel(
					entry.MappedStatus.Version))
			} else {
				RecentVersionColumn.Append(widget.NewLabel("not available"))
			}

			if entry.MappedStatus.Released != "" {
				RecentVersionReleaseDateColumn.Append(widget.NewLabel(
					entry.MappedStatus.Released))
			} else {
				RecentVersionColumn.Append(widget.NewLabel("not available"))
			}
		}
	}
}

func createFyneAppWindow() fyne.Window {
	a = app.New()
	a.Settings().SetTheme(theme.LightTheme())

	InstalledColumn = widget.NewVBox(widget.NewLabelWithStyle(
		"Installed Software                     ",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true}))
	StatusColumn = widget.NewVBox(widget.NewLabelWithStyle(
		"Status                       ",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true}))
	InstalledVersionColumn = widget.NewVBox(widget.NewLabelWithStyle(
		"Installed Version    ",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true}))
	RecentVersionColumn = widget.NewVBox(widget.NewLabelWithStyle(
		"Recent Version       ",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true}))
	RecentVersionReleaseDateColumn = widget.NewVBox(widget.NewLabelWithStyle(
		"Release Date",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true}))

	appList := widget.NewHBox(
		InstalledColumn,
		StatusColumn,
		InstalledVersionColumn,
		RecentVersionColumn,
		RecentVersionReleaseDateColumn,
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
	grid := widget.NewTextGridFromString(OtherSoftwareText)
	grid.ShowLineNumbers = true
	otherWin.SetContent(widget.NewScrollContainer(fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, nil, nil), grid)))
	otherWin.Resize(fyne.NewSize(1024, 600))

	otherWin.Show()
}
