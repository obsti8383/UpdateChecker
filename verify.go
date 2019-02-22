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
	"strconv"
	"strings"
)

// tries to find matches between installed software components and
// software release statii.
// works at least for Firefox, Chrome, OpenVPN and Teamviewer (in current versions)
// TODO: Does not do anything right now beneath logging
func verifyInstalledSoftwareVersions(installedSoftware map[string]installedSoftwareComponent, softwareReleaseStatii map[string]softwareReleaseStatus) []installedSoftwareMapping {
	var returnMapping []installedSoftwareMapping

	for regKey, installedComponent := range installedSoftware {
		var upToDate = false
		var found = false
		var mappedStatValue softwareReleaseStatus
		searchName := strings.Split(installedComponent.DisplayName, ".")[0]
		if searchName != "" {
			for _, statValue := range softwareReleaseStatii {
				searchStatiiName := strings.Split(statValue.Product, ".")[0]

				//fmt.Println("checking if", searchName, " contains ", searchStatKey)
				if strings.Contains(searchName, searchStatiiName) || strings.Contains(searchStatiiName, searchName) {
					//fmt.Printf("Possible match found: Installed software \"%s\" (%s) might match \"%s\" (%s)\n", installedComponent.displayName, installedComponent.displayVersion, statKey, statValue.Version)
					Trace.Printf("Possible match found: Installed software \"%s\" (%s) might match \"%s\" (%s)", installedComponent.DisplayName, installedComponent.DisplayVersion, statValue.Product, statValue.Version)
					found = true
					mappedStatValue = statValue
					if strings.HasPrefix(installedComponent.DisplayVersion, statValue.Version) {
						upToDate = true
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

			/*const STATUS_OUTDATED = 0
			const STATUS_UPTODATE = 1
			const STATUS_UNKNOWN = 2
			}*/
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
			Info.Printf("No Information for %s (%s)", installedComponent.DisplayName, regKey)
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
		// Name:"Microsoft Windows 10", MajorRelease:"1809", Stable:true,
		// Version:"17763.316", Latest:true, Ends:"2020-05-12", Edition:"1809",
		// Product:"Microsoft Windows 10", Released:"2019-02-12"},
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
