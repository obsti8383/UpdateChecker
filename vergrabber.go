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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const VERGRABBER_URL = "http://vergrabber.kingu.pl/vergrabber.json"
const VERGRABBER_FILE = "vergrabber.json"

// fetches current versions of common software from
// http://vergrabber.kingu.pl/vergrabber.json
func getSoftwareVersionsFromVergrabber() map[string]softwareReleaseStatus {
	softwareReleaseStatii := map[string]softwareReleaseStatus{}
	var jsonFromVergrabber []byte

	/// first try to read cached file from filesystem
	jsonFromVergrabber, err := ioutil.ReadFile(VERGRABBER_FILE)
	// handle the error if there is one
	if err != nil {
		Trace.Println("no cached vergrabber.json file available, catching online version")
		jsonFromVergrabber = downloadVergrabberJSON()
	} else {
		updatedDate, err := getVergrabberUpdateDate(jsonFromVergrabber)
		if err != nil {
			Info.Println("Could not get update date from vergrabber.json")
		}
		if !updatedDate.After(time.Now().Add(-time.Hour * 24 * 2)) {
			jsonFromVergrabber = downloadVergrabberJSON()
			updatedDate, err = getVergrabberUpdateDate(jsonFromVergrabber)
			if err != nil {
				Info.Println("Downloaded vergrabber.json is not up-to-date")
				panic("Downloaded vergrabber.json is not up-to-date")
			}
		}
	}

	//Trace.Printf("%s\n", jsonFromVergrabber)

	// TODO: check if vergrabber.json is up to date

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

	return softwareReleaseStatii
}

func downloadVergrabberJSON() []byte {
	// get JSON
	Info.Println("Downloading vergrabber.json")
	fmt.Println("Downloading vergrabber.json")
	resp, err := http.Get(VERGRABBER_URL)
	if err != nil {
		Info.Println("Could not catch vergrabber json from " + VERGRABBER_URL)
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
	outputFile, err := os.Create(VERGRABBER_FILE)
	if err != nil {
		Info.Println(err)
		fmt.Println(err)
	}
	defer outputFile.Close()
	_, err = outputFile.Write(jsonFromVergrabber)
	if err != nil {
		Info.Println(err)
		fmt.Println(err)
	}

	return jsonFromVergrabber
}

func getVergrabberUpdateDate(jsonFromVergrabber []byte) (time.Time, error) {

	r, _ := regexp.Compile(".*updated\": \"(.*)\".*")

	fmt.Println()

	scanner := bufio.NewScanner(bytes.NewReader(jsonFromVergrabber))

	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		if r.MatchString(line) {
			dateString := r.FindStringSubmatch(line)[1]
			//fmt.Println("vergrabber.json updated: " + dateString)
			splitDate := strings.Split(dateString, "-")
			year, err1 := strconv.Atoi(splitDate[0])
			month, err2 := strconv.Atoi(splitDate[1])
			day, err3 := strconv.Atoi(splitDate[2])
			if err1 != nil || err2 != nil || err3 != nil {
				Info.Println("Error getting update date from vergrabber json: Error in integer conversion")
				return time.Now(), errors.New("Could not get update date from vergrabber.json")
			}
			date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
			//fmt.Println("vergrabber.json updated: " + date.String())
			return date, nil
		}
	}

	if err := scanner.Err(); err != nil {
		Info.Println("Error getting update date from vergrabber json: " + err.Error())
		return time.Now(), errors.New("Could not get update date from vergrabber.json")
	}

	return time.Now(), errors.New("Could not get update date from vergrabber.json")
}
