package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/TIBCOSoftware/flogo-cli/util"
	"path/filepath"
)

const (
	itActivity = "activity"
	itTrigger = "trigger"
	itModel = "model"
	itFlow = "flow"
	itPalette = "palette"
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
		if (isPath && (v.Path == itemNameOrPath || v.LocalPath == itemNameOrPath)) || (!isPath && v.Name == itemNameOrPath) {
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
func AddFlogoItem(gb *fgutil.Gb, itemType string, itemPath string, version string, items []*ItemDescriptor, addToSrc bool, ignoreDup bool) (itemConfig *ItemDescriptor, itemConfigPath string) {

	itemPath = strings.Replace(itemPath, "local://", fgutil.FileURIPrefix, 1)

	if ContainsItemPath(items, itemPath) {

		if (ignoreDup) {
			fmt.Fprintf(os.Stdout, "Warning: %s '%s' is already in the project.\n\n", fgutil.Capitalize(itemType), itemPath)
			return nil, ""

		} else {
			fmt.Fprintf(os.Stderr, "Error: %s '%s' is already in the project.\n\n", fgutil.Capitalize(itemType), itemPath)
			os.Exit(2)
		}
	}

	pathInfo, err := fgutil.GetPathInfo(itemPath)

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Invalid path '%s'\n", itemPath)
		os.Exit(2)
	}

	var itemName string

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

			if itemType != itModel && filepath.Base(itemImportPath) == "runtime" {
				itemImportPath = itemImportPath[:len(itemImportPath) - 8]
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

		return &ItemDescriptor{Name: itemName, Path: itemImportPath, Version: "latest", LocalPath: itemPath}, itemConfigPath

	} else {

		//todo handle item already fetched - external or bad cleanup

		//gb vendor delete for now, need proper cleanup on error
		gb.VendorDeleteSilent(itemPath)

		err := gb.VendorFetch(itemPath, version)
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

		return &ItemDescriptor{Name: itemName, Path: itemPath, Version: "latest"}, itemConfigPath
	}
}

// DelFlogoItem deletes an item(activity, model or trigger) from the flogo project
func DelFlogoItem(gb *fgutil.Gb, itemType string, itemNameOrPath string, items []*ItemDescriptor, useSrc bool) []*ItemDescriptor {

	itemNameOrPath = strings.Replace(itemNameOrPath, "local://", fgutil.FileURIPrefix, 1)

	toRemove, itemConfig := GetItemConfig(items, itemNameOrPath)

	if toRemove == -1 {
		fmt.Fprintf(os.Stderr, "Error: %s '%s' is not in the project.\n\n", fgutil.Capitalize(itemType), itemNameOrPath)
		os.Exit(2)
	}

	if itemConfig.Local() {

		// delete it from source and vendor
		toVendorDir :=filepath.Join(gb.VendorPath, itemConfig.Path)
		toSourceDir := filepath.Join(gb.SourcePath, itemConfig.Path)

		os.RemoveAll(toVendorDir)
		os.RemoveAll(toSourceDir)

		//todo clean up empty directory hierarchy

	} else {

		err := gb.VendorDelete(itemConfig.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(2)
		}
	}

	return append(items[:toRemove], items[toRemove + 1:]...)
}