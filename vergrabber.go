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
	"encoding/json"

	"io/ioutil"
)

// fetches current versions of common software from
// http://vergrabber.kingu.pl/vergrabber.json
func getSoftwareVersionsFromVergrabber() map[string]softwareReleaseStatus {
	softwareReleaseStatii := map[string]softwareReleaseStatus{}

	// get JSON
	// TODO: cache vergrabber.json
	//url := "http://vergrabber.kingu.pl/vergrabber.json"
	//resp, err := http.Get(url)
	/// read from file system for now
	jsonFromVergrabber, err := ioutil.ReadFile("vergrabber.json")
	// handle the error if there is one
	if err != nil {
		panic(err)
	}
	// do this now so it won't be forgotten
	//defer resp.Body.Close()

	// reads json as a slice of bytes
	//jsonFromVergrabber, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	panic(err)
	//}

	Info.Printf("%s\n", jsonFromVergrabber)

	// parse JSON
	var f map[string]map[string]map[string]softwareReleaseStatus
	err = json.Unmarshal(jsonFromVergrabber, &f)

	for _, valueSoftwareType := range f {
		//fmt.Println("Typ:", softwareType)
		for softwareName, softwareDetails := range valueSoftwareType {
			//fmt.Println("Name:", softwareName)
			for softwareVersion, softwareVersionDetails := range softwareDetails {
				softwareVersionDetails.Name = softwareName
				softwareVersionDetails.MajorRelease = softwareVersion
				softwareReleaseStatii[softwareName+" "+softwareVersion] = softwareVersionDetails
			}
		}
	}

	return softwareReleaseStatii
}
