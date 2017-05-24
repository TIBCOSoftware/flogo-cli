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

	//cfg := &SettingsConfig{DeviceName: deviceName, Settings:descriptor.Settings}
	//err = generateSourceFiles(env.GetSourceDir(), details, cfg, false)

	err = generateSourceFilesNew(env.GetSourceDir(), descriptor)
	if err != nil {
		return err
	}

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

func generateSourceFilesNew(srcDir string, descriptor *FlogoDeviceDescriptor) error {

	details := GetDevice(descriptor.DeviceType)


	generateMainCode(srcDir, descriptor, details)

	if descriptor.MqttEnabled {

		generateMqttCode(srcDir, descriptor, details)
	}

	//generate triggers
	for _, trgCfg := range descriptor.Triggers {

		//settings := &SettingsConfig{ID: trgCfg.Id, ActionId: Settings:trgCfg.Settings}

		f, _ := os.Create(path.Join(srcDir, trgCfg.Id + ".ino"))
		tpl := triggerTpls[trgCfg.Ref]
		err := RenderTemplate(f, tpl, trgCfg)
		f.Close()
		if err != nil {
			return err
		}
	}

	// generate actions
	for _, actCfg := range descriptor.Actions {

		// Action should define how to do this

		template.ParseFiles()

		t := template.New("action")
		t.Funcs(DeviceFuncMap)
		_, err := t.Parse(actionTpls[actCfg.Ref])
		if err != nil {
			return err
		}

		tmpl := t.New("activity")
		_, err = tmpl.Parse(activityTpls[actCfg.Data.Activity.Ref])
		if err != nil {
			return err
		}

		f, _ := os.Create(path.Join(srcDir, actCfg.Id + ".ino"))

		err = t.Execute(f, actCfg)
		f.Close()
		if err != nil {
			return err
		}
	}

	//template.ParseFiles()

	return nil
}

func generateMainCode(srcDir string, descriptor *FlogoDeviceDescriptor, details *DeviceDetails) error {
	var actionIds []string

	for _, value := range descriptor.Actions {
		actionIds = append(actionIds, value.Id)
	}

	var triggerIds []string
	var mqttTriggerIds []string

	for _, value := range descriptor.Triggers {

		if !strings.Contains(value.Ref, "mqtt") {
			triggerIds = append(triggerIds, value.Id)
		} else {
			mqttTriggerIds = append(mqttTriggerIds, value.Id)
		}
	}

	data := struct {
		MqttEnabled bool
		Actions []string
		Triggers []string
		MqttTriggers []string
	}{
		descriptor.MqttEnabled,
		actionIds,
		triggerIds,
		mqttTriggerIds,
	}

	tpl := details.MainFile
	//todo fix name resolution
	f, _ := os.Create(path.Join(srcDir, "main.ino"))
	err := RenderTemplate(f, tpl, data)
	f.Close()
	if err != nil {
		return err
	}

	return nil
}

func generateMqttCode(srcDir string, descriptor *FlogoDeviceDescriptor, details *DeviceDetails) error {

	settings := &SettingsConfig{DeviceName: descriptor.Name, Settings:descriptor.Settings}

	for name, tpl := range details.MqttFiles {

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
