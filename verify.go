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
	"strconv"
	"strings"
)

func compareVersionStrings(version1, version2 string) int {
	// replace "-" with "." to remediate version strings like 2.4.6-602
	version1 = strings.Replace(version1, "-", ".", -1)
	version2 = strings.Replace(version2, "-", ".", -1)

	v1Split := strings.Split(version1, ".")
	v2Split := strings.Split(version2, ".")
	lenV1 := len(v1Split)
	lenV2 := len(v2Split)

	for i, v1 := range v1Split {
		if i >= lenV2 {
			return -1
		}

		v1Int, err1 := strconv.ParseUint(v1, 10, 64)
		v2Int, err2 := strconv.ParseUint(v2Split[i], 10, 64)
		if err1 != nil || err2 != nil {
			// compare as string
			return strings.Compare(v1, v2Split[i])
		}

		if v1Int > v2Int {
			return 1
		} else if v2Int == v1Int {
			if i == lenV1-1 {
				if lenV1 == lenV2 {
					return 0
				}
				return -1
			}
			// else: go on
		} else {
			return -1
		}
	}

	// TODO: we should return an error here instead
	return -1
}

// verifies installed software versions
func verifyInstalledSoftwareVersions(installedSoftware map[string]installedSoftwareComponent, softwareReleaseStatii map[string]softwareReleaseStatus) []installedSoftwareMapping {
	var returnMapping []installedSoftwareMapping

	for regKey, installedComponent := range installedSoftware {
		var upToDate = false
		var found = false
		var mappedStatValue softwareReleaseStatus

		// ignore list
		if strings.HasPrefix(installedComponent.DisplayName, "Java Auto Updater") {
			continue
		}

		// Firefox (special due to long term support releases)
		if strings.HasPrefix(installedComponent.DisplayName, "Mozilla Firefox") {

			version := installedComponent.DisplayVersion
			versionSplit := strings.Split(version, ".")
			minorVersion := versionSplit[0] + "." + versionSplit[1]
			currentRelease, inStatii := softwareReleaseStatii["Mozilla Firefox "+minorVersion]
			if inStatii {
				mappedStatValue = currentRelease
				found = true
				if compareVersionStrings(currentRelease.Version, version) == 0 {
					upToDate = true
				}
			} else {
				// go through all version and select newest
				for statName, statValue := range softwareReleaseStatii {
					if strings.HasPrefix(statName, "Mozilla Firefox") {
						if mappedStatValue.Version != "" {
							if compareVersionStrings(mappedStatValue.Version, statValue.Version) > 0 {
								// ignore, we already found a newer release
							} else {
								found = true
								mappedStatValue = statValue
							}
						} else {
							found = true
							mappedStatValue = statValue
						}
					}
				}
			}
		}

		// LibreOffice (special due to fresh and still releases)
		if strings.HasPrefix(installedComponent.DisplayName, "LibreOffice") {

			version := installedComponent.DisplayVersion
			versionSplit := strings.Split(version, ".")
			minorVersion := versionSplit[0] + "." + versionSplit[1]
			subMinorVersion := versionSplit[0] + "." + versionSplit[1] + "." + versionSplit[2]
			currentRelease, inStatii := softwareReleaseStatii["LibreOffice "+minorVersion]
			if inStatii {
				mappedStatValue = currentRelease
				found = true
				if compareVersionStrings(currentRelease.Version, subMinorVersion) == 0 {
					upToDate = true
				}
			} else {
				// go through all versions and select newest
				for statName, statValue := range softwareReleaseStatii {
					if strings.HasPrefix(statName, "LibreOffice") {
						if mappedStatValue.Version != "" {
							if compareVersionStrings(mappedStatValue.Version, statValue.Version) > 0 {
								// ignore, we already found a newer release
							} else {
								found = true
								mappedStatValue = statValue
							}
						} else {
							found = true
							mappedStatValue = statValue
						}
					}
				}
			}
		}

		// "Adobe Flash Player"
		if strings.HasPrefix(installedComponent.DisplayName, "Adobe Flash Player") {
			upToDate = false
			found = true
			mappedStatValue = softwareReleaseStatus{
				Name:     "Adobe Flash Player",
				Version:  "out of service - please uninstall",
				Released: "2020-12-31",
				Ends:     "2020-12-31",
			}

		}

		// "Adobe Acrobat Reader"
		if strings.HasPrefix(installedComponent.DisplayName, "Adobe Acrobat Reader") {
			version := installedComponent.DisplayVersion
			adobeName := "Adobe Acrobat Reader"
			if strings.HasPrefix(installedComponent.DisplayName, "Adobe Acrobat Reader DC") {
				// DC version
				adobeName += " DC"
			} else {
				adobeName = installedComponent.DisplayName[0:25]
			}
			Trace.Printf("Adobe name: " + adobeName)

			for statName, statValue := range softwareReleaseStatii {
				if strings.HasPrefix(statName, adobeName) {
					Trace.Printf("Adobe Reader version mapping found for %s", statName)

					mappedStatValue = statValue
					found = true
					if strings.HasPrefix(version, statValue.Version) {
						upToDate = true
					}
				}
			}
		}

		// other software
		softwares := []string{"Google Chrome", "OpenVPN",
			"7-Zip", "TeamViewer",
			"Mozilla Thunderbird", "VeraCrypt", "Java"}
		for _, name := range softwares {
			if strings.HasPrefix(installedComponent.DisplayName, name) {
				version := installedComponent.DisplayVersion
				versionSplit := strings.Split(version, ".")
				minorVersion := versionSplit[0] + "." + versionSplit[1]
				majorVersion := versionSplit[0]
				for statName, statValue := range softwareReleaseStatii {
					if strings.HasPrefix(statName, name+" "+minorVersion) {
						Trace.Printf("Minor version mapping found for %s", statName)

						mappedStatValue = statValue
						found = true
						if strings.HasPrefix(version, statValue.Version) { //compareVersionStrings(, ) == 0 {
							upToDate = true
						}
					} else if strings.HasPrefix(statName, name+" "+majorVersion) {
						Trace.Printf("Major version mapping found for %s", statName)

						mappedStatValue = statValue
						found = true
						if strings.HasPrefix(version, statValue.Version) {
							upToDate = true
						}
					} else if strings.HasPrefix(statName, name) {
						if mappedStatValue.Version != "" {
							// ignore, we already found a correct release
						} else {
							Trace.Printf("Name only mapping found for %s", statName)

							compareResult := compareVersionStrings(statValue.Version, version)
							if compareResult == 0 {
								found = true
								mappedStatValue = statValue
								upToDate = true
							} else if compareResult < 0 {
								// perhaps newer version as we know of (e.g. for Java) -> leave on unknown
							} else {
								found = true
								mappedStatValue = statValue
							}
						}
					}
				}
			}
		}

		// build mapping object
		if upToDate {
			returnMapping = append(returnMapping, installedSoftwareMapping{
				Name:              installedComponent.DisplayName,
				Status:            StatusUpToDate,
				InstalledSoftware: installedComponent,
				MappedStatus:      mappedStatValue,
			})
			Info.Printf("%s seems up to date (%s)", installedComponent.DisplayName, installedComponent.DisplayVersion)
		} else if found {
			returnMapping = append(returnMapping, installedSoftwareMapping{
				Name:              installedComponent.DisplayName,
				Status:            StatusOutdated,
				InstalledSoftware: installedComponent,
				MappedStatus:      mappedStatValue,
			})
			Info.Printf("%s seems outdated!! (%s)", installedComponent.DisplayName, installedComponent.DisplayVersion)
		} else {
			returnMapping = append(returnMapping, installedSoftwareMapping{
				Name:              installedComponent.DisplayName,
				Status:            StatusUnknown,
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
		windowsReleaseName := "Microsoft Windows 10 " + string(windowsVersion.ReleaseID)
		Trace.Println("windowsReleaseName: ", windowsReleaseName)
		Trace.Println("windowsVersion.UBR: ", windowsVersion.UBR)
		Trace.Println("string(windowsVersion.UBR): ", strconv.FormatUint(windowsVersion.UBR, 10))

		windowsVersionString := windowsVersion.CurrentBuild + "." + strconv.FormatUint(windowsVersion.UBR, 10)
		Trace.Println("windowsVersionString: ", windowsVersionString)
		uptodateRelease := softwareReleaseStatii[windowsReleaseName]
		Trace.Println("uptodateRelease: ", uptodateRelease)

		if uptodateRelease.Version == windowsVersionString {
			Info.Printf("Windows seems up to date")
			status = StatusUpToDate
		} else {
			Info.Printf("Windows seems outdated!!")
			status = StatusOutdated
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
	}

	Info.Printf("Windows Version <= Windows 8: Not supported by Update Checker")
	return installedSoftwareMapping{
		Name:   windowsVersion.ProductName,
		Status: StatusUnknown,
		InstalledSoftware: installedSoftwareComponent{
			DisplayName:    windowsVersion.ProductName,
			DisplayVersion: windowsVersion.CurrentBuild,
			Publisher:      "Microsoft",
		},
		MappedStatus: softwareReleaseStatus{},
	}

}
