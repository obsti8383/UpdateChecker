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
	"errors"

	"github.com/sqweek/dialog"
	"golang.org/x/sys/windows/registry"
)

// struct for the registry keys needed to read out installed software
type registryKeys struct {
	rootKey registry.Key
	path    string
	flags   uint32
}

// WindowsVersion can hold the relevant Major, Minor and so on version numbers of Windows (10)
type WindowsVersion struct {
	CurrentMajorVersionNumber, UBR, CurrentMinorVersionNumber uint64
	CurrentBuild, ReleaseID, ProductName                      string
}

// gets Windows version numbers (Major, Minor and CurrentBuild)
func getWindowsVersion() (windowsVersion WindowsVersion, err error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion", registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
	if err != nil {
		return WindowsVersion{0, 0, 0, "", "", ""}, errors.New("Could not get version information from registry")
	}
	defer k.Close()

	pn, _, err := k.GetStringValue("ProductName")
	if err != nil {
		return WindowsVersion{0, 0, 0, "", "", ""}, errors.New("Could not get version information from registry - ProductName")
	}

	cb, _, err := k.GetStringValue("CurrentBuild")
	if err != nil {
		return WindowsVersion{0, 0, 0, "", "", pn}, errors.New("Could not get version information from registry - CurrentBuild")
	}

	maj, _, err := k.GetIntegerValue("CurrentMajorVersionNumber")
	if err != nil {
		return WindowsVersion{0, 0, 0, cb, "", pn}, errors.New("Could not get version information from registry - CurrentMajorVersionNumber")
	}

	// since 20H2 MS has messed up version numbering - working around this here
	// DisplayVersion seems to be new. We take this field, if available, otherwise ReleaseID
	relID := ""
	relDisplayVersion, _, _ := k.GetStringValue("DisplayVersion")
	if relDisplayVersion != "" {
		relID = relDisplayVersion
	} else {
		relID, _, err = k.GetStringValue("ReleaseId")
		if err != nil {
			return WindowsVersion{maj, 0, 0, cb, "", pn}, errors.New("Could not get version information from registry - ReleaseId")
		}
	}

	ubr, _, err := k.GetIntegerValue("UBR")
	if err != nil {
		return WindowsVersion{maj, 0, 0, cb, relID, pn}, errors.New("Could not get version information from registry - UBR")
	}

	min, _, err := k.GetIntegerValue("CurrentMinorVersionNumber")
	if err != nil {
		return WindowsVersion{maj, ubr, 0, cb, relID, pn}, errors.New("Could not get version information from registry - CurrentMinorVersionNumber")
	}

	return WindowsVersion{maj, ubr, min, cb, relID, pn}, nil
}

func checkWindowsVersionError(windowsVersion WindowsVersion, err error) {
	if err == nil {
		Info.Printf("Windows Product Name: %s", windowsVersion.ProductName)
		Info.Printf("Windows Version: %d.%d.%s.%d",
			windowsVersion.CurrentMajorVersionNumber,
			windowsVersion.CurrentMinorVersionNumber,
			windowsVersion.CurrentBuild, windowsVersion.UBR)
		Info.Printf("Windows Release ID: %s", windowsVersion.ReleaseID)
	} else {
		if windowsVersion.ProductName != "" {
			Info.Printf("Windows Product Name: %s", windowsVersion.ProductName)
			if windowsVersion.CurrentBuild != "" {
				Info.Printf("Windows Current Build: %s", windowsVersion.CurrentBuild)
			}
		} else {
			Info.Printf("Error getting Windows Version: %s", err)
			dialog.Message("%s %s", "Error getting Windows Version:", err).Title("Error").Error()
		}
	}
}
