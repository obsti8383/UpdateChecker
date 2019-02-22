// Update Checker
// Copyright (C) 2019  Florian Probst
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
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
)

const logpath = "UpdateChecker.log"

var installedSoftwareMappings []installedSoftwareMapping

// Loggers for log output (we only need info and trace, errors have to be
// displayed in the GUI)
var (
	Trace *log.Logger
	Info  *log.Logger
)

// initLogging inits loggers
func initLogging(traceHandle io.Writer, infoHandle io.Writer) {
	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

// this struct is used for filling in software attributes for most current stable
// version of software available for download (attributes must match the
// corresponding JSON elements for Vergrabber.json)
type softwareReleaseStatus struct {
	Name         string // filled in manually from Vergrabber.json
	MajorRelease string // filled in manually from Vergrabber.json
	Stable       bool   // automatically unmarshalled from Vergrabber.json
	Version      string // automatically unmarshalled from Vergrabber.json
	Latest       bool   // automatically unmarshalled from Vergrabber.json
	Ends         string // automatically unmarshalled from Vergrabber.json
	Edition      string // automatically unmarshalled from Vergrabber.json
	Product      string // automatically unmarshalled from Vergrabber.json
	Released     string // automatically unmarshalled from Vergrabber.json
}

// this struct is used for filling in software attributes for software
// that is actually installed on the system
type installedSoftwareComponent struct {
	DisplayName    string
	DisplayVersion string
	Publisher      string
}

const STATUS_OUTDATED = 0
const STATUS_UPTODATE = 1
const STATUS_UNKNOWN = 2

type installedSoftwareMapping struct {
	Name              string
	Status            int
	InstalledSoftware installedSoftwareComponent
	MappedStatus      softwareReleaseStatus
}

func main() {
	// init logging
	var logfile, err = os.Create(logpath)
	if err != nil {
		panic(err)
	}
	initLogging(logfile, logfile)

	// fetch Windows version
	windowsVersion, err := getWindowsVersion()
	if err == nil {
		Info.Printf("Windows Version: %d.%d.%s.%d",
			windowsVersion.CurrentMajorVersionNumber,
			windowsVersion.CurrentMinorVersionNumber,
			windowsVersion.CurrentBuild, windowsVersion.UBR)
		Info.Printf("Windows Release ID: %s", windowsVersion.ReleaseId)
	} else {
		Info.Printf("Error getting Windows Version: %s", err)
	}

	// fetch installed software
	foundSoftware, err := getInstalledSoftware()
	if err == nil {
		for key, soft := range foundSoftware {
			Info.Printf("%s: %s %s (%s)", key, soft.DisplayName, soft.DisplayVersion, soft.Publisher)
		}
	} else {
		return
	}

	// fetch software current release information from Vergrabber
	softwareReleaseStatii := getSoftwareVersionsFromVergrabber()
	Trace.Printf(fmt.Sprintf("Software Releases from Vergrabber: %#v\n", softwareReleaseStatii))
	//fmt.Println("Software Releases from Vergrabber:\n", softwareReleaseStatii)

	// get mappings between installed software and currentReleases
	installedSoftwareMappings = verifyInstalledSoftwareVersions(foundSoftware, softwareReleaseStatii)

	// sort installed software mappings
	sort.Slice(installedSoftwareMappings, func(i, j int) bool {
		return installedSoftwareMappings[i].Name < installedSoftwareMappings[j].Name
	})
	sort.Slice(installedSoftwareMappings, func(i, j int) bool {
		return installedSoftwareMappings[i].Status < installedSoftwareMappings[j].Status
	})

	// verify OS patch level against Vergrabber
	windowsMapping, err := verifyOSPatchlevel(windowsVersion, softwareReleaseStatii)
	if err != nil {
		Info.Printf(err.Error())
	}
	Trace.Printf("WindowsMapping: %#v\n", windowsMapping)

	// create Combined Mapping for Windows itself and installed software
	newMappings := make([]installedSoftwareMapping, 0)
	newMappings = append(newMappings, windowsMapping)
	installedSoftwareMappings = append(newMappings, installedSoftwareMappings...)

	//t, _ := template.ParseFiles("main.html")
	//t.Execute(os.Stdout, installedSoftwareMappings)

	// present with Webserver
	http.HandleFunc("/", mainHttpHandler) // setting router rule
	listener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}

	// open browser
	err = openBrowser("http://localhost:3000/")
	if err != nil {
		log.Println(err)
	}

	// Start the blocking server loop.
	log.Fatal(http.Serve(listener, nil))

	return
}
