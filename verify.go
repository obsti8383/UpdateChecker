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
	"fmt"
	"strconv"
	"strings"
)

// tries to find matches between installed software components and
// software release statii.
// works at least for Firefox, Chrome, OpenVPN, Adobe Flash and Teamviewer (in current versions)
func verifyInstalledSoftwareVersions(installedSoftware map[string]installedSoftwareComponent, softwareReleaseStatii map[string]softwareReleaseStatus) []installedSoftwareMapping {
	var returnMapping []installedSoftwareMapping

	for regKey, installedComponent := range installedSoftware {
		var upToDate = false
		var found = false
		var mappedStatValue softwareReleaseStatus
		searchName := strings.Split(installedComponent.DisplayName, ".")[0]
		if searchName != "" {
			for statName, statValue := range softwareReleaseStatii {
				statNameArray := strings.Split(statName, " ")
				if len(statNameArray) > 2 {
					statName = statNameArray[0] + " " + statNameArray[1]
				}
				//fmt.Printf("statName: %s\n", statName)
				if strings.Contains(searchName, statName) || strings.Contains(statName, searchName) {
					fmt.Printf("Possible match found: Installed software \"%s\" (%s) might match \"%s\" (%s)\n", installedComponent.DisplayName, installedComponent.DisplayVersion, statName, statValue.Version)
					Trace.Printf("Possible match found: Installed software \"%s\" (%s) might match \"%s\" (%s)", installedComponent.DisplayName, installedComponent.DisplayVersion, statName, statValue.Version)
					found = true
					mappedStatValue = statValue
					if installedComponent.DisplayVersion == statValue.Version {
						upToDate = true
						break
					}
				}
			}
		}
		if upToDate {
			returnMapping = append(returnMapping, installedSoftwareMapping{
				Name:              installedComponent.DisplayName,
				Status:            STATUS_UPTODATE,
				InstalledSoftware: installedComponent,
				MappedStatus:      mappedStatValue,
			})
			Info.Printf("%s seems up to date (%s)", installedComponent.DisplayName, installedComponent.DisplayVersion)
		} else if found {
			returnMapping = append(returnMapping, installedSoftwareMapping{
				Name:              installedComponent.DisplayName,
				Status:            STATUS_OUTDATED,
				InstalledSoftware: installedComponent,
				MappedStatus:      mappedStatValue,
			})
			Info.Printf("%s seems outdated!! (%s)", installedComponent.DisplayName, installedComponent.DisplayVersion)
		} else {
			returnMapping = append(returnMapping, installedSoftwareMapping{
				Name:              installedComponent.DisplayName,
				Status:            STATUS_UNKNOWN,
				InstalledSoftware: installedComponent,
				MappedStatus:      mappedStatValue,
			})
			Trace.Printf("No Information for %s (%s)", installedComponent.DisplayName, regKey)
		}
	}

	return returnMapping
}

// verify OS patchlevel
func verifyOSPatchlevel(windowsVersion WindowsVersion, softwareReleaseStatii map[string]softwareReleaseStatus) (installedSoftwareMapping, error) {
	var status int

	if windowsVersion.CurrentMajorVersionNumber == 10 {
		// Windows 10
		windowsReleaseName := "Microsoft Windows 10 " + string(windowsVersion.ReleaseId)
		Trace.Println("windowsReleaseName: ", windowsReleaseName)
		Trace.Println("windowsVersion.UBR: ", windowsVersion.UBR)
		Trace.Println("string(windowsVersion.UBR): ", strconv.FormatUint(windowsVersion.UBR, 10))

		windowsVersionString := windowsVersion.CurrentBuild + "." + strconv.FormatUint(windowsVersion.UBR, 10)
		Trace.Println("windowsVersionString: ", windowsVersionString)
		uptodateRelease := softwareReleaseStatii[windowsReleaseName]
		Trace.Println("uptodateRelease: ", uptodateRelease)

		if uptodateRelease.Version == windowsVersionString {
			status = STATUS_UPTODATE
		} else {
			status = STATUS_OUTDATED
		}

		return installedSoftwareMapping{
			Name:   "Microsoft Windows 10",
			Status: status,
			InstalledSoftware: installedSoftwareComponent{
				DisplayName:    windowsReleaseName,
				DisplayVersion: windowsVersionString,
				Publisher:      "Microsoft",
			},
			MappedStatus: uptodateRelease,
		}, nil
	}
	return installedSoftwareMapping{}, errors.New("Could not find corresponding windows version")
}
