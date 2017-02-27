package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/TIBCOSoftware/flogo-cli/util"
	"github.com/TIBCOSoftware/flogo-cli/env"
)

// BuildPreProcessor interface for build pre-processors
type BuildPreProcessor interface {
	PrepareForBuild(env env.Project)
}

// CreateApp creates an application from the specified json application descriptor
func CreateApp(env env.Project, appJson string, appName string) error {

	descriptor, err := ParseAppDescriptor(appJson)
	if err != nil {
		return err
	}

	if appName != "" {
		var appObj map[string]interface{}

		err := json.Unmarshal([]byte(appJson), &appObj)
		if err != nil {
			return err
		}

		appObj["name"] = appName

		updApp, err := json.Marshal(appObj)
		if err != nil {
			return err
		}

		appJson = string(updApp)
	}

	currentDir, err := os.Getwd()

	if err != nil {
		return err
	}

	appDir := fgutil.Path(currentDir, descriptor.Name)
	env.Init(appDir)
	err = env.Create(false)
	if err != nil {
		return err
	}

	err = fgutil.CreateFileFromString(fgutil.Path(appDir, "flogo.json"), appJson)
	if err != nil {
		return err
	}

	//todo allow ability to specify flogo-lib version
	env.InstallDependency(pathFlogoLib, "")

	refs := ExtractRefs(descriptor)

	for _, ref := range refs {
		path, version := splitVersion(ref)
		err = env.InstallDependency(path, version)
		if err != nil {
			return err
		}
	}

	// create source files
	cmdPath := fgutil.Path(env.GetSourceDir(), strings.ToLower(descriptor.Name))
	os.MkdirAll(cmdPath, 0777)

	createMainGoFile(cmdPath)
	createImportsGoFile(cmdPath, refs)

	return nil
}

// BuildApp build the flogo application
func BuildApp(env env.Project, customPreProcessor BuildPreProcessor) (err error) {

	if customPreProcessor != nil {
		customPreProcessor.PrepareForBuild(env)
	}

	//todo do standard pre-processing
	// regenerate imports?

	err = env.Build()
	if err != nil {
		return err
	}

	err = fgutil.MoveFiles(env.GetBinDir(), env.GetRootDir())
	if err != nil {
		return err
	}

	return
}

// InstallDependency install a dependency
func InstallDependency(env env.Project, path string, version string) error {

	return env.InstallDependency(path, version)
}

// ListDependencies lists all installed dependencies
func ListDependencies(env env.Project, cType ContribType) ([]*Dependency, error) {

	vendorSrc := env.GetVendorDir()
	var deps []*Dependency

	err := filepath.Walk(vendorSrc, func(path string, info os.FileInfo, _ error) error {

		if !info.IsDir() {

			ref := refPath(vendorSrc, path)

			switch info.Name() {
			case "action.json":
				if cType == 0 || cType == ACTION {
					desc, err := readDescriptor(path, info)
					if err == nil && desc.Type == "flogo:action" {
						deps = append(deps, &Dependency{ContribType: ACTION, Ref: ref})
					}
				}
			case "trigger.json":
				if cType == 0 || cType == TRIGGER {
					desc, err := readDescriptor(path, info)
					if err == nil && desc.Type == "flogo:trigger" {
						deps = append(deps, &Dependency{ContribType: TRIGGER, Ref: ref})
					}
				}
			case "activity.json":
				if cType == 0 || cType == ACTIVITY {
					desc, err := readDescriptor(path, info)
					if err == nil && desc.Type == "flogo:activity" {
						deps = append(deps, &Dependency{ContribType: ACTIVITY, Ref: ref})
					}
				}
			case "flow-model.json":
				if cType == 0 || cType == FLOW_MODEL {
					desc, err := readDescriptor(path, info)
					if err == nil && desc.Type == "flogo:flow-model" {
						deps = append(deps, &Dependency{ContribType: FLOW_MODEL, Ref: ref})
					}
				}
			}

		}

		return nil
	})

	return deps, err
}

func refPath(vendorSrc string, filePath string) string {

	startIdx := len(vendorSrc) + 1
	endIdx := strings.LastIndex(filePath, "/")

	return filePath[startIdx:endIdx]
}

func readDescriptor(path string, info os.FileInfo) (*Descriptor, error) {

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("error: " + err.Error())
		return nil, err
	}

	return ParseDescriptor(string(raw))
}

// ParseAppDescriptor parse the application descriptor
func ParseDescriptor(descJson string) (*Descriptor, error) {
	descriptor := &Descriptor{}

	err := json.Unmarshal([]byte(descJson), descriptor)

	if err != nil {
		return nil, err
	}

	return descriptor, nil
}

// ParseAppDescriptor parse the application descriptor
func ParseAppDescriptor(appJson string) (*FlogoAppDescriptor, error) {
	descriptor := &FlogoAppDescriptor{}

	err := json.Unmarshal([]byte(appJson), descriptor)

	if err != nil {
		return nil, err
	}

	return descriptor, nil
}
