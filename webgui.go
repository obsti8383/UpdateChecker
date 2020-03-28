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
	"bufio"
	"html/template"
	"os"
	"os/exec"
	"runtime"

	"github.com/sqweek/dialog"
)

const resultFileName = "updatechecker_result.html"

// open opens the specified URL in the default browser of the user.
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func outputResultsInBrowser(installedSoftwareMappings []installedSoftwareMapping) {
	Trace.Println("Generating results html page...")
	// Write HTML output
	outputFile, err := os.Create(resultFileName)
	if err != nil {
		Info.Println(err)
		dialog.Message("%s", "Error writing output file. Aborting!").Title("Error").Error()
		os.Exit(1)
	}
	defer outputFile.Close()
	outputWriter := bufio.NewWriter(outputFile)

	Trace.Println("Executing Template...")
	t, _ := template.ParseFiles("main.html")
	t.Execute(outputWriter, installedSoftwareMappings)
	outputWriter.Flush()
	outputFile.Close()

	// open browser
	Trace.Println("Opening Browser...")
	err = openBrowser(resultFileName)
	if err != nil {
		Info.Println(err)
	}
}

func deleteResultHTML() {
	// if file exists...
	if _, err := os.Stat(resultFileName); err == nil {
		// delete it
		err := os.Remove(resultFileName)
		if err != nil {
			Info.Println(err)
			dialog.Message("%s", "Error deleting old output file. Aborting!").Title("Error").Error()
			os.Exit(1)
		}

	}
}
