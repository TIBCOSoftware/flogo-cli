package device

import (
	"encoding/json"
	"os"
	"path"
	"strings"

	"github.com/TIBCOSoftware/flogo-cli/util"
	"io"
	"text/template"
)

// BuildPreProcessor interface for build pre-processors
type BuildPreProcessor interface {
	PrepareForBuild(env Project) error
}

// CreateDevice creates an device project from the specified json device descriptor
func CreateDevice(env Project, deviceJson string, deviceDir string, deviceName string) error {

	descriptor, err := ParseDeviceDescriptor(deviceJson)
	if err != nil {
		return err
	}

	if deviceName != "" {
		// override the application name

		altJson := strings.Replace(deviceJson, `"`+descriptor.Name+`"`, `"`+deviceName+`"`, 1)
		altDescriptor, err := ParseDeviceDescriptor(altJson)

		//see if we can get away with simple replace so we don't reorder the existing json
		if err == nil && altDescriptor.Name == deviceName {
			deviceJson = altJson
		} else {
			//simple replace didn't work so we have to unmarshal & re-marshal the supplied json
			var appObj map[string]interface{}
			err := json.Unmarshal([]byte(deviceJson), &appObj)
			if err != nil {
				return err
			}

			appObj["name"] = deviceName

			updApp, err := json.MarshalIndent(appObj, "", "  ")
			if err != nil {
				return err
			}
			deviceJson = string(updApp)
		}

		descriptor.Name = deviceName
	}

	env.Init(deviceDir)

	details := GetDevice(descriptor.DeviceType)

	err = env.Create(details.Board)
	if err != nil {
		return err
	}
	err = fgutil.CreateFileFromString(path.Join(deviceDir, "device.json"), deviceJson)
	if err != nil {
		return err
	}

	cfg := &SettingsConfig{DeviceName: deviceName, Settings:descriptor.Settings}
	err = generateSourceFiles(env.GetSourceDir(), details, cfg, false)

	return nil
}

type PrepareOptions struct {
	PreProcessor    BuildPreProcessor
}

// PrepareDevice do all pre-build setup and pre-processing
func PrepareDevice(project Project, options *PrepareOptions) (err error) {

	if options == nil {
		options = &PrepareOptions{}
	}

	if options.PreProcessor != nil {
		err = options.PreProcessor.PrepareForBuild(project)
		if err != nil {
			return err
		}
	}

	//load descriptor
	appJson, err := fgutil.LoadLocalFile(path.Join(project.GetRootDir(),"device.json"))
	if err != nil {
		return err
	}

	descriptor, err := ParseDeviceDescriptor(appJson)
	if err != nil {
		return err
	}

	details := GetDevice(descriptor.DeviceType)

	cfg := &SettingsConfig{DeviceName: descriptor.Name, Settings:descriptor.Settings}
	err = generateSourceFiles(project.GetSourceDir(), details, cfg, true)
	if err != nil {
		return err
	}

	for name, id := range details.Libs {

		err := project.InstallLib(name, id)
		if err != nil {
			return err
		}
	}

	return nil
}

type BuildOptions struct {
	*PrepareOptions

	SkipPrepare bool
}

// BuildDevice build the flogo application
func BuildDevice(env Project, options *BuildOptions) (err error) {

	if options == nil {
		options = &BuildOptions{}
	}

	if !options.SkipPrepare {
		err = PrepareDevice(env, options.PrepareOptions)

		if err != nil {
			return err
		}
	}

	err = env.Build()
	if err != nil {
		return err
	}

	return
}

// UploadDevice upload the device application
func UploadDevice(env Project) error {

	err := env.Upload()

	return err
}

//func readDescriptor(path string, info os.FileInfo) (*Descriptor, error) {
//
//	raw, err := ioutil.ReadFile(path)
//	if err != nil {
//		fmt.Println("error: " + err.Error())
//		return nil, err
//	}
//
//	return ParseDescriptor(string(raw))
//}

func generateSourceFiles(srcDir string, details *DeviceDetails, settings *SettingsConfig, skipMain bool) error {

	for name, tpl := range details.Files {

		if skipMain && strings.HasPrefix(name, "main") {
			continue
		}

		f, _ := os.Create(path.Join(srcDir, name))
		err := RenderTemplate(f, tpl, settings)
		f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

//RenderFileTemplate renders the specified template
func RenderTemplate(w io.Writer, tpl string, data interface{}) error {

	t := template.New("source")
	t.Funcs(DeviceFuncMap)

	t.Parse(tpl)
	err := t.Execute(w, data)

	return err
}
