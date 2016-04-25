package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/TIBCOSoftware/flogo/util"
	"path/filepath"
)

const (
	itActivity = "activity"
	itTrigger = "trigger"
	itModel = "model"
	itFlow = "flow"
)

// ContainsItemPath determines if the path exists in  list of ItemConfigs
func ContainsItemPath(list []*ItemDescriptor, path string) bool {
	for _, v := range list {
		if v.Path == path {
			return true
		}
	}
	return false
}

// ContainsItemName determines if the path exists in  list of ItemConfigs
func ContainsItemName(list []*ItemDescriptor, name string) bool {
	for _, v := range list {
		if v.Name == name {
			return true
		}
	}
	return false
}

// GetItemConfig gets the item config for the specified path or name
func GetItemConfig(list []*ItemDescriptor, itemNameOrPath string) (int, *ItemDescriptor) {

	isPath := strings.Contains(itemNameOrPath, "/")

	for i, v := range list {
		if (isPath && v.Path == itemNameOrPath) || (!isPath && v.Name == itemNameOrPath) {
			return i, v
		}
	}
	return -1, nil
}

func getItemInfo(itemFile *os.File, itemType string) (string, string) {

	itemConfig := &struct {
		Name       string `json:"name"`
		ImportPath string `json:"importpath"`
	}{}

	jsonParser := json.NewDecoder(itemFile)

	if err := jsonParser.Decode(itemConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to parse %s.json, file may be corrupted.\n - %s\n", itemType, err.Error())
		os.Exit(2)
	}

	return itemConfig.Name, itemConfig.ImportPath
}

// AddFlogoItem adds an item(activity, model or trigger) to the flogo project
func AddFlogoItem(gb *fgutil.Gb, itemType string, itemPath string, items []*ItemDescriptor, addToSrc bool) (itemConfig *ItemDescriptor, itemConfigPath string) {

	itemPath = strings.Replace(itemPath, "local://", fgutil.FileURIPrefix, 1)

	if ContainsItemPath(items, itemPath) {
		fmt.Fprintf(os.Stderr, "Error: %s '%s' is already in the project.\n\n", fgutil.Capitalize(itemType), itemPath)
		os.Exit(2)
	}

	pathInfo, err := fgutil.GetPathInfo(itemPath)

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Invalid path '%s'\n", itemPath)
		os.Exit(2)
	}

	var itemName string
	var isLocal bool

	if pathInfo.IsURL {

		if !pathInfo.IsLocal {
			fmt.Fprint(os.Stderr, "Error: Add %s, does not support URL scheme '%s'\n", itemType, pathInfo.FileURL.Scheme)
			os.Exit(2)
		}

		usesGb := false

		itemConfigPath = filepath.Join(pathInfo.FilePath, itemType + ".json")
		itemFile, err := os.Open(itemConfigPath)

		if err != nil {
			itemFile.Close()
			itemConfigPath = path(pathInfo.FilePath, "src", itemType + ".json")
			itemFile, err = os.Open(itemConfigPath)

			usesGb = true
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Path '%s' is not a flogo-based %s.\n\n", itemPath, itemType)
				itemFile.Close()
				os.Exit(2)
			}
		}

		var itemImportPath string
		itemName, itemImportPath = getItemInfo(itemFile, itemType)
		itemFile.Close()

		var toDir string
		if len(itemImportPath) > 0 {

			if itemType != itModel && filepath.Base(itemImportPath) == "rt" {
				itemImportPath = itemImportPath[:len(itemImportPath) - 3]
			}
		} else {
			itemImportPath = filepath.Join(itemType, itemName)
		}

		if addToSrc {
			toDir = filepath.Join(gb.SourcePath, itemImportPath)
		} else {
			toDir = filepath.Join(gb.VendorPath, itemImportPath)
		}

		fromDir := pathInfo.FilePath

		if usesGb {
			fromDir = filepath.Join(pathInfo.FilePath, "src")
		}

		fgutil.CopyDir(fromDir, toDir)

		isLocal = true

	} else {

		//todo handle item already fetched - external or bad cleanup

		//gb vendor delete for now, need proper cleanup on error
		gb.VendorDeleteSilent(itemPath)

		err := gb.VendorFetch(itemPath)
		if err != nil {
			os.Exit(2)
		}

		itemConfigPath = filepath.Join(gb.VendorPath, itemPath, itemType + ".json")
		itemFile, err := os.Open(itemConfigPath)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Path '%s' is not a flogo-based %s.\n\n", itemConfigPath, itemType)
			itemFile.Close()
			os.Exit(2)
		}

		itemName, _ = getItemInfo(itemFile, itemType)
		itemFile.Close()
	}

	return &ItemDescriptor{Name: itemName, Path: itemPath, Version: "latest", Local: isLocal}, itemConfigPath
}

// DelFlogoItem deletes an item(activity, model or trigger) from the flogo project
func DelFlogoItem(gb *fgutil.Gb, itemType string, itemNameOrPath string, items []*ItemDescriptor, useSrc bool) []*ItemDescriptor {

	itemNameOrPath = strings.Replace(itemNameOrPath, "local://", fgutil.FileURIPrefix, 1)

	toRemove, itemConfig := GetItemConfig(items, itemNameOrPath)

	if toRemove == -1 {
		fmt.Fprintf(os.Stderr, "Error: %s '%s' is not in the project.\n\n", fgutil.Capitalize(itemType), itemNameOrPath)
		os.Exit(2)
	}

	itemPath := itemConfig.Path

	pathInfo, err := fgutil.GetPathInfo(itemPath)

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Invalid path '%s'\n", itemPath)
		os.Exit(2)
	}

	if pathInfo.IsURL {

		if !pathInfo.IsLocal {
			fmt.Fprint(os.Stderr, "Error: Add %s, does not support URL scheme '%s'\n", itemType, pathInfo.FileURL.Scheme)
			os.Exit(2)
		}

		// delete it from source and vendor

		toVendorDir := path(gb.VendorPath, itemType, itemConfig.Name)
		toSourceDir := path(gb.SourcePath, itemType, itemConfig.Name)

		os.RemoveAll(toVendorDir)
		os.RemoveAll(toSourceDir)

	} else {

		err := gb.VendorDelete(itemPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(2)
		}
	}

	return append(items[:toRemove], items[toRemove + 1:]...)
}

func CleanupItem(dir string, itemType string, itemName string) int {

	configFile := itemType + ".json"

	var toRemove []string

	delFunc := func(path string, f os.FileInfo, err error) (e error) {

		if !f.IsDir() && f.Name() == configFile {

			currentItemFile, err := os.Open(path)
			if err == nil {
				currentItemName, _ := getItemInfo(currentItemFile, itemType)
				currentItemFile.Close()

				if currentItemName == itemName {
					toRemove = append(toRemove, filepath.Dir(path))
				}
			}
		}

		return nil
	}

	filepath.Walk(dir, delFunc)

	deleted := 0

	for _, d := range toRemove {

		if strings.HasPrefix(d, dir) {
			os.RemoveAll(d)
			deleted++
		}
	}

	return deleted
}