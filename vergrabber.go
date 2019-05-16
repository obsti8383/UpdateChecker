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
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

const vergrabberURL = "http://vergrabber.kingu.pl/vergrabber.json"
const vergrabberFile = "vergrabber.json"

// fetches current versions of common software from
// http://vergrabber.kingu.pl/vergrabber.json
func getSoftwareVersionsFromVergrabber() map[string]softwareReleaseStatus {
	var jsonFromVergrabber []byte
	softwareReleaseStatii := map[string]softwareReleaseStatus{}

	/// first try to read cached file from filesystem
	jsonFromVergrabber, err := ioutil.ReadFile(vergrabberFile)
	if err != nil {
		Trace.Println("no cached vergrabber.json file available, catching online version")
		jsonFromVergrabber = downloadAndCacheVergrabberJSON()
		if !isVergrabberJSONUptodate(jsonFromVergrabber) {
			Info.Println("Downloaded vergrabber.json is not up-to-date")
			panic("Downloaded vergrabber.json is not up-to-date")
		}
	} else {
		// checking if cached file is still up-to-date
		if !isVergrabberJSONUptodate(jsonFromVergrabber) {
			Trace.Println("cached vergrabber.json file is outdated, catching online version")
			jsonFromVergrabber = downloadAndCacheVergrabberJSON()
			jsonUptodate := isVergrabberJSONUptodate(jsonFromVergrabber)
			if !jsonUptodate {
				Info.Println("Downloaded vergrabber.json is not up-to-date")
				panic("Downloaded vergrabber.json is not up-to-date")
			}
		}
	}

	// parse JSON
	var f map[string]map[string]map[string]softwareReleaseStatus
	err = json.Unmarshal(jsonFromVergrabber, &f)

	for softwareType, valueSoftwareType := range f {
		//fmt.Println("Typ:", softwareType)
		if softwareType == "client" || softwareType == "server" {
			for softwareName, softwareDetails := range valueSoftwareType {
				//fmt.Println("Name:", softwareName)
				for softwareVersion, softwareVersionDetails := range softwareDetails {
					softwareVersionDetails.Name = softwareName
					softwareVersionDetails.MajorRelease = softwareVersion
					softwareReleaseStatii[softwareName+" "+softwareVersion] = softwareVersionDetails
				}
			}
		}
	}

	//Trace.Printf(fmt.Sprintf("Software Releases from Vergrabber: %#v\n", softwareReleaseStatii))

	return softwareReleaseStatii
}

func isVergrabberJSONUptodate(jsonFromVergrabber []byte) bool {
	updatedDate, err := getVergrabberUpdateDate(jsonFromVergrabber)
	if err != nil {
		Info.Println("Downloaded vergrabber.json is not up-to-date")
		panic("Downloaded vergrabber.json is not up-to-date")
	}

	// second: verify date
	if !updatedDate.After(time.Now().Add(-time.Hour * 24 * 2)) {
		// JSON is outdated
		return false
	}
	// JSON is not outdated
	return true
}

func downloadAndCacheVergrabberJSON() []byte {
	// get JSON
	Info.Println("Downloading vergrabber.json")
	resp, err := http.Get(vergrabberURL)
	if err != nil {
		Info.Println("Could not catch vergrabber json from " + vergrabberURL)
		panic(err)
	}
	defer resp.Body.Close()

	// reads json as a slice of bytes
	jsonFromVergrabber, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Info.Println("Error reading vergrabber json: " + err.Error())
		panic(err)
	}

	// Write file to disk
	outputFile, err := os.Create(vergrabberFile)
	if err != nil {
		Info.Println(err)
	}
	defer outputFile.Close()
	_, err = outputFile.Write(jsonFromVergrabber)
	if err != nil {
		Info.Println(err)
	}

	return jsonFromVergrabber
}

func getVergrabberUpdateDate(jsonFromVergrabber []byte) (time.Time, error) {
	// regexp to search for "updated:" entry and its date
	r, _ := regexp.Compile(".*updated\": \"(.*)\".*")

	scanner := bufio.NewScanner(bytes.NewReader(jsonFromVergrabber))
	for scanner.Scan() {
		line := scanner.Text()
		if r.MatchString(line) {
			dateString := r.FindStringSubmatch(line)[1]
			date, err := time.Parse("2006-01-02", dateString)
			if err != nil {
				Info.Println("Error getting update date from vergrabber json: Error in date conversion")
				return time.Now(), errors.New("Could not get update date from vergrabber.json")
			}
			Trace.Printf("vergrabber.json updated: " + date.String())
			return date, nil
		}
	}

	if err := scanner.Err(); err != nil {
		Info.Println("Error getting update date from vergrabber json: " + err.Error())
		return time.Now(), errors.New("Could not get update date from vergrabber.json")
	}

	return time.Now(), errors.New("Could not get update date from vergrabber.json")
}
