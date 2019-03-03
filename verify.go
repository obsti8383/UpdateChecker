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
	"regexp"
	"strconv"
	"strings"
)

// tries to find matches between installed software components and
// software release statii.
// works at least for Firefox, Chrome, OpenVPN, Adobe Flash, Acrobat Reader, 7-Zip and Teamviewer (in current versions)
func verifyInstalledSoftwareVersions(installedSoftware map[string]installedSoftwareComponent, softwareReleaseStatii map[string]softwareReleaseStatus) []installedSoftwareMapping {
	var returnMapping []installedSoftwareMapping

	for regKey, installedComponent := range installedSoftware {
		var upToDate = false
		var found = false
		var mappedStatValue softwareReleaseStatus

		// regexp for just getting [a-zA-Z ] from beginning of string
		reg := regexp.MustCompile("^[-0-9a-zA-ZäöüÄÖÜß\t\n\v\f\r{}_()® ]+")
		reg2 := regexp.MustCompile(" [0-9]+\\.[0-9]+")

		searchNameWithMajorVersion := strings.Split(installedComponent.DisplayName, ".")[0]
		searchNameWithoutVersion := reg.FindString(installedComponent.DisplayName)
		index := reg2.FindStringIndex(installedComponent.DisplayName)
		if len(index) > 0 {
			searchNameWithoutVersion = installedComponent.DisplayName[:index[0]]
		}
		Trace.Printf("searchNameWithMajorVersion: " + searchNameWithMajorVersion)
		Trace.Printf("searchNameWithoutVersion:   " + searchNameWithoutVersion)
		if searchNameWithMajorVersion != "" {
			for statName, statValue := range softwareReleaseStatii {
				statNameArray := strings.Split(statName, " ")
				if len(statNameArray) > 2 {
					statName = statNameArray[0] + " " + statNameArray[1]
				}
				//fmt.Printf("statName: %s\n", statName)
				if strings.Contains(searchNameWithMajorVersion, statName) || strings.Contains(statName, searchNameWithMajorVersion) {
					//fmt.Printf("Possible match found: Installed software \"%s\" (%s) might match \"%s\" (%s)\n", installedComponent.DisplayName, installedComponent.DisplayVersion, statName, statValue.Version)
					Trace.Printf("Possible match found: Installed software \"%s\" (%s) might match \"%s\" (%s)", installedComponent.DisplayName, installedComponent.DisplayVersion, statName, statValue.Version)
					// special case for Firefox: ESR and standard releases are
					// both contained in vergrabber.json. Therefore if we find two
					// matches, we have to decide which one to match
					if found {
						// we found the item before, with an outdated/not equal software version
						if installedComponent.DisplayVersion == statValue.Version {
							mappedStatValue = statValue
							upToDate = true
							break
						} else {
							if mappedStatValue.Version < statValue.Version {
								// always map the highest version
								mappedStatValue = statValue
							}
						}
					} else {
						// standard case: we (first) found an matching item
						found = true
						mappedStatValue = statValue
						if installedComponent.DisplayVersion == statValue.Version {
							upToDate = true
							break
						}
					}

				} else if searchNameWithoutVersion != "" && strings.Contains(searchNameWithoutVersion, statName) || strings.Contains(statName, searchNameWithoutVersion) {
					Trace.Printf("Possible match found: Installed software \"%s\" (%s) might match \"%s\" (%s)", installedComponent.DisplayName, installedComponent.DisplayVersion, statName, statValue.Version)
					if found {
						// we found the item before, with an outdated/not equal software version
						if installedComponent.DisplayVersion == statValue.Version {
							mappedStatValue = statValue
							upToDate = true
							break
						} else {
							if mappedStatValue.Version < statValue.Version {
								// always map the highest version
								mappedStatValue = statValue
							}
						}
					} else {
						// standard case: we (first) found an matching item
						found = true
						mappedStatValue = statValue
						if installedComponent.DisplayVersion == statValue.Version {
							upToDate = true
							break
						}
					}
				}
			}
		}

		// build mapping object
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
func verifyOSPatchlevel(windowsVersion WindowsVersion, softwareReleaseStatii map[string]softwareReleaseStatus) installedSoftwareMapping {
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
			Info.Printf("Windows seems up to date")
			status = STATUS_UPTODATE
		} else {
			Info.Printf("Windows seems outdated!!")
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
		}
	} else {
		Info.Printf("Windows Version <= Windows 8: Not supported by Update Checker")
		return installedSoftwareMapping{
			Name:   windowsVersion.ProductName,
			Status: STATUS_UNKNOWN,
			InstalledSoftware: installedSoftwareComponent{
				DisplayName:    windowsVersion.ProductName,
				DisplayVersion: windowsVersion.CurrentBuild,
				Publisher:      "Microsoft",
			},
			MappedStatus: softwareReleaseStatus{},
		}
	}
}
