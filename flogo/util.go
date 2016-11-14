package main

import (
	"github.com/TIBCOSoftware/flogo-cli/util"
	"github.com/TIBCOSoftware/flogo-cli/config"
)

func updateProjectFiles(gb *fgutil.Gb, projectDescriptor *config.FlogoProjectDescriptor) {
	fgutil.WriteJSONtoFile(fileDescriptor, projectDescriptor)
	createImportsGoFile(gb.CodeSourcePath, projectDescriptor)
}
