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

	"golang.org/x/sys/windows/registry"
)

// reads installed software from Microsoft Windows official registry keys
func getInstalledSoftware() (map[string]installedSoftwareComponent, error) {
	// Software from Uninstall registry keys
	regKeysUninstall := []registryKeys{
		{registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall", registry.ENUMERATE_SUB_KEYS | registry.QUERY_VALUE | registry.WOW64_64KEY},
		{registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall", registry.ENUMERATE_SUB_KEYS | registry.QUERY_VALUE | registry.WOW64_32KEY},
		{registry.CURRENT_USER, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall", registry.ENUMERATE_SUB_KEYS | registry.QUERY_VALUE},
	}

	foundSoftware := make(map[string]installedSoftwareComponent)

	for i := 0; i < len(regKeysUninstall); i++ {
		Info.Printf("%d:%s", regKeysUninstall[i].rootKey, regKeysUninstall[i].path)
		key, err := registry.OpenKey(regKeysUninstall[i].rootKey, regKeysUninstall[i].path, regKeysUninstall[i].flags)
		if err != nil {
			Info.Printf("Could not open registry key %s due to error %s", regKeysUninstall[i].path, err.Error())
			return nil, errors.New(fmt.Sprintf("Could not open registry key %s due to error %s", regKeysUninstall[i].path, err.Error()))
		}
		defer key.Close()

		//keyInfo, _ := key.Stat()
		//Info.Printf("Number of subkeys: %i", int(keyInfo.SubKeyCount))
		subKeys, err := key.ReadSubKeyNames(0)
		if err != nil {
			Info.Printf("Could not read sub keys of registry key %s due to error %s", regKeysUninstall[i].path, err.Error())
			return nil, errors.New(fmt.Sprintf("Could not read sub keys of registry key %s due to error %s", regKeysUninstall[i].path, err.Error()))
		}

		for j := 0; j < len(subKeys); j++ {
			subKey, err := registry.OpenKey(regKeysUninstall[i].rootKey, regKeysUninstall[i].path+"\\"+subKeys[j], regKeysUninstall[i].flags)
			if err != nil {
				Info.Printf("Could not open registry key %s due to error %s", subKeys[j], err.Error())
				return nil, errors.New(fmt.Sprintf("Could not open registry key %s due to error %s", subKeys[j], err.Error()))
			}
			defer subKey.Close()

			displayName, _, _ := subKey.GetStringValue("DisplayName")
			if displayName == "" {
				displayName = subKeys[j]
			}
			displayVersion, _, _ := subKey.GetStringValue("DisplayVersion")
			publisher, _, _ := subKey.GetStringValue("Publisher")
			Trace.Printf("getInstalledSoftware: %s: %s %s (%s)", subKeys[j], displayName, displayVersion, publisher)

			newSoftwareFound := installedSoftwareComponent{displayName, displayVersion, publisher}
			foundSoftware[subKeys[j]] = newSoftwareFound
		}
	}

	return foundSoftware, nil
}
