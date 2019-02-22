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
	"errors"

	"golang.org/x/sys/windows/registry"
)

// struct for the registry keys needed to read out installed software
type registryKeys struct {
	rootKey registry.Key
	path    string
	flags   uint32
}

type WindowsVersion struct {
	CurrentMajorVersionNumber, UBR, CurrentMinorVersionNumber uint64
	CurrentBuild, ReleaseId                                   string
}

// gets Windows version numbers (Major, Minor and CurrentBuild)
func getWindowsVersion() (windowsVersion WindowsVersion, err error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion", registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
	if err != nil {
		return WindowsVersion{0, 0, 0, "", ""}, errors.New("Could not get version information from registry")
	}
	defer k.Close()

	// BUG: This does not work with Windows 8.1! There is only CurrentBuild and CurrentVersion (which include Major and Minor, e.g. "6.3")

	maj, _, err := k.GetIntegerValue("CurrentMajorVersionNumber")
	if err != nil {
		return WindowsVersion{0, 0, 0, "", ""}, errors.New("Could not get version information from registry - CurrentMajorVersionNumber")
	}

	relId, _, err := k.GetStringValue("ReleaseId")
	if err != nil {
		return WindowsVersion{maj, 0, 0, "", ""}, errors.New("Could not get version information from registry - ReleaseId")
	}

	ubr, _, err := k.GetIntegerValue("UBR")
	if err != nil {
		return WindowsVersion{maj, 0, 0, "", relId}, errors.New("Could not get version information from registry - UBR")
	}

	min, _, err := k.GetIntegerValue("CurrentMinorVersionNumber")
	if err != nil {
		return WindowsVersion{maj, ubr, 0, "", relId}, errors.New("Could not get version information from registry - CurrentMinorVersionNumber")
	}

	cb, _, err := k.GetStringValue("CurrentBuild")
	if err != nil {
		return WindowsVersion{maj, ubr, min, "", relId}, errors.New("Could not get version information from registry - CurrentBuild")
	}

	return WindowsVersion{maj, ubr, min, cb, relId}, nil
}
