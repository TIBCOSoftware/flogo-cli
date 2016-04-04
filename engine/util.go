package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/TIBCOSoftware/flogo-tools/fg"
	"github.com/TIBCOSoftware/flogo-tools/fgutil"
)

type getItems func(cfg *EngineProjectConfig) []*ItemConfig

// AddEngineItem adds an item(activity, model or trigger) to the engine
func AddEngineItem(c flogo.Command, projectConfig *EngineProjectConfig, itemType string, args []string, gi getItems, useSrc bool) (itemConfig *ItemConfig, itemConfigPath string) {

	itemPath := args[0]

	if ContainsItemPath(gi(projectConfig), itemPath) {
		fmt.Fprintf(os.Stderr, "Error: %s '%s' is already in engine project.\n\n", fgutil.Capitalize(itemType), itemPath)
		os.Exit(2)
	}

	var altPath string
	var localPath string

	if strings.HasPrefix(itemPath, "local://") {
		localPath = itemPath[8:]
		altPath = "file://" + localPath
	} else if strings.HasPrefix(itemPath, "file://") {
		localPath = itemPath[6:]
		altPath = "local://" + localPath
	}

	//todo: handle paths that end in "rt"

	var itemName string

	if len(localPath) > 0 {

		if ContainsItemPath(gi(projectConfig), altPath) {
			fmt.Fprintf(os.Stderr, "Error: %s '%s' is already in engine project.\n\n", fgutil.Capitalize(itemType), itemPath)
			os.Exit(2)
		}

		usesGb := false

		itemConfigPath = path(localPath, itemType + ".json")
		fmt.Print("itemConfigPath: " + itemConfigPath)
		itemFile, err := os.Open(itemConfigPath)

		if err != nil {
			itemFile.Close()
			itemFile, err = os.Open(path(localPath, "src", itemType + ".json"))

			usesGb = true
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Path '%s' is not a flogo-based %s.\n\n", itemPath, itemType)
				itemFile.Close()
				os.Exit(2)
			}
		}

		itemName = getItemName(itemFile, itemType)
		itemFile.Close()

		sourcePath := path("src")
		vendorPath := path("vendor", "src")

		toDir := path(vendorPath, itemType, itemConfig.Name)

		if useSrc {
			toDir = path(sourcePath, itemType, itemConfig.Name)
		}

		fromDir := localPath

		if usesGb {
			fromDir = path(localPath, "src")
		}

		fgutil.CopyDir(fromDir, toDir)

	} else {

		//todo handle item already fetched - external or bad cleanup

		//gb vendor delete for now, need proper cleanup on error
		cmd := exec.Command("gb", "vendor", "delete", itemPath)
		err := cmd.Run()

		cmd = exec.Command("gb", "vendor", "fetch", itemPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			os.Exit(2)
		}

		vendorPath := path("vendor", "src")

		itemConfigPath = path(vendorPath, itemPath, itemType + ".json")
		itemFile, err := os.Open(itemConfigPath)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Path '%s' is not a flogo-based %s.\n\n", itemConfigPath, itemType)
			itemFile.Close()
			os.Exit(2)
		}

		itemName = getItemName(itemFile, itemType)
		itemFile.Close()
	}

	return &ItemConfig{Name:itemName, Path: itemPath, Version: "latest"}, itemConfigPath
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

// AddEngineItem adds an item(activity, model or trigger) to the engine
func DelEngineItem(c flogo.Command, itemType string, args []string, gi getItems, useSrc bool) (idx int, engineConfig *EngineProjectConfig) {

	configFile, err := os.Open(fileProjectConfig)

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Current working directory is not a flogo-based engine project.\n\n")
		os.Exit(2)
	}

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: %s name or path not specified\n\n", fgutil.Capitalize(itemType))
		Tool().CmdUsage(c)
	}

	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Error: Too many arguments given\n\n")
		Tool().CmdUsage(c)
	}

	engineConfig = &EngineProjectConfig{}
	jsonParser := json.NewDecoder(configFile)

	if err = jsonParser.Decode(engineConfig); err != nil {
		fmt.Fprint(os.Stderr, "Error: Unable to parse engine.json, file may be corrupted.\n\n")
		os.Exit(2)
	}

	configFile.Close()

	itemNameOrPath := args[0]

	i, itemConfig := GetItemConfig(gi(engineConfig), itemNameOrPath)

	if i == -1 {
		fmt.Fprintf(os.Stderr, "Error: %s '%s' is not in engine project.\n\n", fgutil.Capitalize(itemType), itemNameOrPath)
		os.Exit(2)
	}

	itemPath := itemConfig.Path

	var localPath string

	if strings.HasPrefix(itemPath, "local://") {
		localPath = itemPath[8:]
	} else if strings.HasPrefix(itemPath, "file://") {
		localPath = itemPath[6:]
	}

	if len(localPath) > 0 {

		// delete it from source and vendor

		sourcePath := path("src")
		vendorPath := path("vendor", "src")

		toVendorDir := path(vendorPath, itemType, itemConfig.Name)
		toSourceDir := path(sourcePath, itemType, itemConfig.Name)

		os.RemoveAll(toVendorDir)
		os.RemoveAll(toSourceDir)

	} else {

		cmd := exec.Command("gb", "vendor", "delete", itemPath)
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(2)
		}
	}

	return i, engineConfig
}

func updateProjectConfigFiles(engineConfig *EngineProjectConfig) {
	fgutil.WriteJSONtoFile(fileProjectConfig, engineConfig)

	sourcePath := path("src", engineConfig.Name)

	// create imports test Go file
	f, _ := os.Create(path(sourcePath, fileImportsGo))
	fgutil.RenderTemplate(f, tplImportsGoFile, engineConfig)
	f.Close()
}