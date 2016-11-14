package config

import (
	"os"
	"fmt"
	"encoding/json"
)

func LoadProjectDescriptor() *FlogoProjectDescriptor {

	projectDescriptorFile, err := os.Open(FileProjectDescriptor)

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Current working directory is not a flogo-based engine project.\n\n")
		os.Exit(2)
	}

	projectDescriptor := &FlogoProjectDescriptor{}
	jsonParser := json.NewDecoder(projectDescriptorFile)

	if err = jsonParser.Decode(projectDescriptor); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to parse flogo.json, file may be corrupted.\n - %s\n", err.Error())
		os.Exit(2)
	}

	projectDescriptorFile.Close()

	return projectDescriptor
}

func LoadTriggersConfig() *TriggersConfig {

	triggersConfigFile, err := os.Open("bin/" + FileTriggersConfig)

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: bin/triggers.json not found.\n\n")
		os.Exit(2)
	}

	triggersConfig := &TriggersConfig{}
	jsonParser := json.NewDecoder(triggersConfigFile)

	if err = jsonParser.Decode(triggersConfig); err != nil {
		fmt.Fprint(os.Stderr, "Error: Unable to parse application triggers.json, file may be corrupted.\n\n")
		os.Exit(2)
	}

	triggersConfigFile.Close()

	return triggersConfig;
}