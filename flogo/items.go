package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/TIBCOSoftware/flogo/util"
)

const (
	itActivity = "activity"
	itTrigger  = "trigger"
	itModel    = "model"
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

func getItemName(itemFile *os.File, itemType string) string {

	itemConfig := &struct {
		Name string `json:"name"`
	}{}

	jsonParser := json.NewDecoder(itemFile)

	if err := jsonParser.Decode(itemConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to parse %s.json, file may be corrupted.\n\n", itemType)
		os.Exit(2)
	}

	return itemConfig.Name
}

// AddFlogoItem adds an item(activity, model or trigger) to the flogo project
func AddFlogoItem(gb *fgutil.Gb, itemType string, itemPath string, items []*ItemDescriptor, addToSrc bool) (itemConfig *ItemDescriptor, itemConfigPath string) {

	itemPath = strings.Replace(itemPath, "file://", "local://", 1)

	if ContainsItemPath(items, itemPath) {
		fmt.Fprintf(os.Stderr, "Error: %s '%s' is already in the project.\n\n", fgutil.Capitalize(itemType), itemPath)
		os.Exit(2)
	}

	localPath, local := extractPathFromLocalURI(itemPath)

	var itemName string
	var isLocal bool

	if local {

		usesGb := false

		itemConfigPath = path(localPath, itemType+".json")
		itemFile, err := os.Open(itemConfigPath)

		if err != nil {
			itemFile.Close()
			itemConfigPath = path(localPath, "src", itemType+".json")
			itemFile, err = os.Open(itemConfigPath)

			usesGb = true
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Path '%s' is not a flogo-based %s.\n\n", itemPath, itemType)
				itemFile.Close()
				os.Exit(2)
			}
		}

		itemName = getItemName(itemFile, itemType)
		itemFile.Close()

		toDir := path(gb.VendorPath, itemType, itemName)

		if addToSrc {
			toDir = path(gb.SourcePath, itemType, itemName)
		}

		fromDir := localPath

		if usesGb {
			fromDir = path(localPath, "src")
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

		itemConfigPath = path(gb.VendorPath, itemPath, itemType+".json")
		itemFile, err := os.Open(itemConfigPath)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Path '%s' is not a flogo-based %s.\n\n", itemConfigPath, itemType)
			itemFile.Close()
			os.Exit(2)
		}

		itemName = getItemName(itemFile, itemType)
		itemFile.Close()
	}

	return &ItemDescriptor{Name: itemName, Path: itemPath, Version: "latest", Local: isLocal}, itemConfigPath
}

// DelFlogoItem deletes an item(activity, model or trigger) from the flogo project
func DelFlogoItem(gb *fgutil.Gb, itemType string, itemNameOrPath string, items []*ItemDescriptor, useSrc bool) []*ItemDescriptor {

	itemNameOrPath = strings.Replace(itemNameOrPath, "file://", "local://", 1)

	toRemove, itemConfig := GetItemConfig(items, itemNameOrPath)

	if toRemove == -1 {
		fmt.Fprintf(os.Stderr, "Error: %s '%s' is not in the project.\n\n", fgutil.Capitalize(itemType), itemNameOrPath)
		os.Exit(2)
	}

	itemPath := itemConfig.Path

	_, local := extractPathFromLocalURI(itemPath)

	if local {

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

	return append(items[:toRemove], items[toRemove+1:]...)
}
