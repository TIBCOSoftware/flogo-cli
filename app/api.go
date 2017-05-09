package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/TIBCOSoftware/flogo-cli/util"
	"github.com/TIBCOSoftware/flogo-cli/env"
)

// BuildPreProcessor interface for build pre-processors
type BuildPreProcessor interface {
	PrepareForBuild(env env.Project) error
}

// CreateApp creates an application from the specified json application descriptor
func CreateApp(env env.Project, appJson string, appDir string, appName string, vendorDir string) error {

	descriptor, err := ParseAppDescriptor(appJson)
	if err != nil {
		return err
	}

	if appName != "" {
		// override the application name

		altJson := strings.Replace(appJson, `"`+descriptor.Name+`"`, `"`+appName+`"`, 1)
		altDescriptor, err := ParseAppDescriptor(altJson)

		//see if we can get away with simple replace so we don't reorder the existing json
		if err == nil && altDescriptor.Name == appName {
			appJson = altJson
		} else {
			//simple replace didn't work so we have to unmarshal & re-marshal the supplied json
			var appObj map[string]interface{}
			err := json.Unmarshal([]byte(appJson), &appObj)
			if err != nil {
				return err
			}

			appObj["name"] = appName

			updApp, err := json.MarshalIndent(appObj, "", "  ")
			if err != nil {
				return err
			}
			appJson = string(updApp)
		}

		descriptor.Name = appName
	}

	env.Init(appDir)
	err = env.Create(false, vendorDir)
	if err != nil {
		return err
	}

	err = fgutil.CreateFileFromString(path.Join(appDir, "flogo.json"), appJson)
	if err != nil {
		return err
	}

	//todo allow ability to specify flogo-lib version
	env.InstallDependency(pathFlogoLib, "")

	deps := ExtractDependencies(descriptor)

	for _, dep := range deps {
		path, version := splitVersion(dep.Ref)
		err = env.InstallDependency(path, version)
		if err != nil {
			return err
		}
	}

	// create source files
	cmdPath := path.Join(env.GetSourceDir(), strings.ToLower(descriptor.Name))
	os.MkdirAll(cmdPath, 0777)

	createMainGoFile(cmdPath,"")
	createImportsGoFile(cmdPath, deps)

	return nil
}

type PrepareOptions struct {
	PreProcessor    BuildPreProcessor
	OptimizeImports bool
	EmbedConfig     bool
}

// PrepareApp do all pre-build setup and pre-processing
func PrepareApp(env env.Project, options *PrepareOptions) (err error) {

	if options == nil {
		options = &PrepareOptions{}
	}

	if options.PreProcessor != nil {
		err = options.PreProcessor.PrepareForBuild(env)
		if err != nil {
			return err
		}
	}

	//generate metadatas
	err = generateGoMetadatas(env)
	if err != nil {
		return err
	}

	//load descriptor
	appJson, err := fgutil.LoadLocalFile(path.Join(env.GetRootDir(),"flogo.json"))

	if err != nil {
		return err
	}
	descriptor, err := ParseAppDescriptor(appJson)
	if err != nil {
		return err
	}

	//generate imports file
	var deps []*Dependency

	if options.OptimizeImports {

		deps = ExtractDependencies(descriptor)

	} else {
		deps, err = ListDependencies(env, 0)
	}

	cmdPath := path.Join(env.GetSourceDir(), strings.ToLower(descriptor.Name))
	createImportsGoFile(cmdPath, deps)

	if options.EmbedConfig {
		createEmbeddedAppGoFile(cmdPath, appJson)
	} else {
		removeEmbeddedAppGoFile(cmdPath)
	}

	return
}

type BuildOptions struct {
	*PrepareOptions

	SkipPrepare bool
}

// BuildApp build the flogo application
func BuildApp(env env.Project, options *BuildOptions) (err error) {

	if options == nil {
		options = &BuildOptions{}
	}

	if !options.SkipPrepare {
		err = PrepareApp(env, options.PrepareOptions)

		if err != nil {
			return err
		}
	}

	err = env.Build()
	if err != nil {
		return err
	}

	if !options.EmbedConfig {
		fgutil.CopyFile(path.Join(env.GetRootDir(), fileDescriptor), path.Join(env.GetBinDir(), fileDescriptor))
		if err != nil {
			return err
		}
	} else {
		os.Remove(path.Join(env.GetBinDir(), fileDescriptor))
	}

	return
}

// InstallPalette install a palette
func InstallPalette(env env.Project, path string) error {

	var file []byte

	file, _ = ioutil.ReadFile(path)

	var paletteDescriptor *FlogoPaletteDescriptor
	err := json.Unmarshal(file, &paletteDescriptor)

	var deps []Dependency

	if err != nil {
		err = json.Unmarshal(file, &deps)
	} else {
		deps = paletteDescriptor.Extensions
	}

	if err != nil {
		return err
		//fmt.Fprint(os.Stderr, "Error: Unable to parse palette descriptor, file may be corrupted.\n\n")
		//os.Exit(2)
	}

	for _, dep := range deps {
		err = env.InstallDependency(dep.Ref, "")
		if err != nil {
			return err
		}
	}

	//fmt.Fprintf(os.Stdout, "Adding Palette: %s [%s]\n\n", paletteDescriptor.Name, paletteDescriptor.Description)

	return nil
}

// InstallDependency install a dependency
func InstallDependency(env env.Project, path string, version string) error {

	return env.InstallDependency(path, version)
}

// UninstallDependency uninstall a dependency
func UninstallDependency(env env.Project, path string) error {

	return env.UninstallDependency(path)
}

// ListDependencies lists all installed dependencies
func ListDependencies(env env.Project, cType ContribType) ([]*Dependency, error) {

	vendorSrc := env.GetVendorSrcDir()
	var deps []*Dependency

	err := filepath.Walk(vendorSrc, func(filePath string, info os.FileInfo, _ error) error {

		if !info.IsDir() {

			switch info.Name() {
			case "action.json":
				if cType == 0 || cType == ACTION {
					ref := refPath(vendorSrc, filePath)
					desc, err := readDescriptor(filePath, info)
					if err == nil && desc.Type == "flogo:action" {
						deps = append(deps, &Dependency{ContribType: ACTION, Ref: ref})
					}
				}
			case "trigger.json":
				//temporary hack to handle old contrib dir layout
				dir := filePath[0:len(filePath)-12]
				if _, err := os.Stat(fmt.Sprintf("%s/../trigger.json", dir)); err == nil {
					//old trigger.json, ignore
					return nil
				}
				if cType == 0 || cType == TRIGGER {
					ref := refPath(vendorSrc, filePath)
					desc, err := readDescriptor(filePath, info)
					if err == nil && desc.Type == "flogo:trigger" {
						deps = append(deps, &Dependency{ContribType: TRIGGER, Ref: ref})
					}
				}
			case "activity.json":
				//temporary hack to handle old contrib dir layout
				dir := filePath[0:len(filePath)-13]
				if _, err := os.Stat(fmt.Sprintf("%s/../activity.json", dir)); err == nil {
					//old activity.json, ignore
					return nil
				}
				if cType == 0 || cType == ACTIVITY {
					ref := refPath(vendorSrc, filePath)
					desc, err := readDescriptor(filePath, info)
					if err == nil && desc.Type == "flogo:activity" {
						deps = append(deps, &Dependency{ContribType: ACTIVITY, Ref: ref})
					}
				}
			case "flow-model.json":
				if cType == 0 || cType == FLOW_MODEL {
					ref := refPath(vendorSrc, filePath)
					desc, err := readDescriptor(filePath, info)
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
	endIdx := strings.LastIndex(filePath, string(os.PathSeparator))

	return strings.Replace(filePath[startIdx:endIdx], string(os.PathSeparator), "/", -1)
}

func readDescriptor(path string, info os.FileInfo) (*Descriptor, error) {

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("error: " + err.Error())
		return nil, err
	}

	return ParseDescriptor(string(raw))
}

func generateGoMetadatas(env env.Project) error {
	//todo optimize metadata recreation to minimize compile times
	dependencies, err := ListDependencies(env, 0)

	if err != nil {
		return err
	}

	for _, dependency := range dependencies {
		createMetadata(env, dependency)
	}

	return nil
}

func createMetadata(env env.Project, dependency *Dependency) error {

	vendorSrc := env.GetVendorSrcDir()
	mdFilePath := path.Join(vendorSrc, dependency.Ref)
	mdGoFilePath := path.Join(vendorSrc, dependency.Ref)
	pkg := path.Base(mdFilePath)

	tplMetadata := tplMetadataGoFile

	switch dependency.ContribType {
	case ACTION:
		mdFilePath = path.Join(mdFilePath, "action.json")
		mdGoFilePath = path.Join(mdGoFilePath, "action_metadata.go")
	case TRIGGER:
		mdFilePath = path.Join(mdFilePath, "trigger.json")
		mdGoFilePath = path.Join(mdGoFilePath, "trigger_metadata.go")
		tplMetadata = tplTriggerMetadataGoFile
	case ACTIVITY:
		mdFilePath = path.Join(mdFilePath, "activity.json")
		mdGoFilePath = path.Join(mdGoFilePath, "activity_metadata.go")
		tplMetadata = tplActivityMetadataGoFile
	default:
		return nil
	}

	raw, err := ioutil.ReadFile(mdFilePath)
	if err != nil {
		return err
	}

	info := &struct {
		Package      string
		MetadataJSON string
	}{
		Package:      pkg,
		MetadataJSON: string(raw),
	}

	f, _ := os.Create(mdGoFilePath)
	fgutil.RenderTemplate(f, tplMetadata, info)
	f.Close()

	return nil
}

var tplMetadataGoFile = `package {{.Package}}

var jsonMetadata = ` + "`{{.MetadataJSON}}`" + `

func getJsonMetadata() string {
	return jsonMetadata
}
`

var tplActivityMetadataGoFile = `package {{.Package}}

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

var jsonMetadata = ` + "`{{.MetadataJSON}}`" + `

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(NewActivity(md))
}
`

var tplTriggerMetadataGoFile = `package {{.Package}}

import (
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

var jsonMetadata = ` + "`{{.MetadataJSON}}`" + `

// init create & register trigger factory
func init() {
	md := trigger.NewMetadata(jsonMetadata)
	trigger.RegisterFactory(md.ID, NewFactory(md))
}
`

// ParseDescriptor parse a descriptor
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
