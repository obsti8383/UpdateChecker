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

//go:generate goversioninfo -icon=icon.ico

package main

import (
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

const version = "0.2.2"

const logpath = "UpdateChecker.log"

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

//StatusOutdated means that the software is not up-to-date
const StatusOutdated = 0

//StatusUpToDate means that the software is up-to-date
const StatusUpToDate = 1

//StatusUnknown means that the software status or the software itself is unknown
const StatusUnknown = 2

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

	doSelfUpdate()

	var installedSoftwareMappings []installedSoftwareMapping

	// cleanup from last run
	deleteResultHTML()

	// fetch Windows version
	windowsVersion, err := getWindowsVersion()
	checkWindowsVersionError(windowsVersion, err)

	// fetch installed software
	foundSoftware, err := getInstalledSoftware()
	if err == nil {
		for key, soft := range foundSoftware {
			Trace.Printf("%s: %s %s (%s)", key, soft.DisplayName, soft.DisplayVersion, soft.Publisher)
		}
	} else {
		return
	}

	// fetch software current release information from Vergrabber
	softwareReleaseStatii := getSoftwareVersionsFromVergrabber()

	// get mappings between installed software and currentReleases
	installedSoftwareMappings = verifyInstalledSoftwareVersions(foundSoftware, softwareReleaseStatii)

	// sort installed software mappings
	sort.Slice(installedSoftwareMappings, func(i, j int) bool {
		if installedSoftwareMappings[i].Status < installedSoftwareMappings[j].Status {
			return true
		} else if installedSoftwareMappings[i].Status > installedSoftwareMappings[j].Status {
			return false
		} else {
			return strings.ToUpper(installedSoftwareMappings[i].Name) < strings.ToUpper(installedSoftwareMappings[j].Name)
		}
	})

	// verify OS patch level against Vergrabber
	windowsMapping := verifyOSPatchlevel(windowsVersion, softwareReleaseStatii)

	// create Combined Mapping for Windows itself and installed software
	newMappings := make([]installedSoftwareMapping, 0)
	newMappings = append(newMappings, windowsMapping)
	installedSoftwareMappings = append(newMappings, installedSoftwareMappings...)

	// write results to HTML file and open in browser
	outputResultsInBrowser(installedSoftwareMappings)
}

func doSelfUpdate() {
	v := semver.MustParse(version)
	latest, err := selfupdate.UpdateSelf(v, "obsti8383/UpdateChecker")
	if err != nil {
		Info.Println("Binary update failed:", err)
		return
	}
	if latest.Version.Equals(v) {
		// latest version is the same as current version. It means current binary is up to date.
		Info.Println("Current binary is the latest version", version)
	} else {
		Info.Println("Successfully updated to version", latest.Version)
		Info.Println("Release note:\n", latest.ReleaseNotes)
	}
}
