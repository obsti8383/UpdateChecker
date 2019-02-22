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
