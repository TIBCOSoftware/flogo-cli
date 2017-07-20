package device

import (
	"errors"
	"path"

	"github.com/TIBCOSoftware/flogo-cli/util"
)

var triggerContribs = make(map[string]*TriggerContrib)
var actionContribs = make(map[string]*ActionContrib)
var activityContribs = make(map[string]*ActivityContrib)

type TriggerContrib struct {
	Template   string
	Descriptor *TriggerDescriptor
}

func (tc *TriggerContrib) Libs() []*Lib {
	return tc.Descriptor.Libs
}

type ActionContrib struct {
	Ref      string
	Template string
	libs     []*Lib
}

func (ac *ActionContrib) GetActivities(tree *FlowTree) []*ActivityContrib {

	var activities []*ActivityContrib

	for _, task := range tree.AllTasks {
		contrib := activityContribs[task.ActivityRef]
		activities = append(activities, contrib)
	}

	return activities
}

type ActivityContrib struct {
	Template   string
	Descriptor *ActivityDescriptor
}

func (ac *ActivityContrib) Libs() []*Lib {
	return ac.Descriptor.Libs
}

func LoadTriggerContrib(proj Project, ref string) (*TriggerContrib, error) {
	contrib, exists := triggerContribs[ref]

	if !exists {
		proj.InstallContribution(ref, "")

		descFile := path.Join(proj.GetContributionDir(), ref, "trigger.json")
		descJson, err := fgutil.LoadLocalFile(descFile)

		if err != nil {
			return nil, err
		}

		desc, err := ParseTriggerDescriptor(descJson)
		//validate that device trigger

		//validate that supports device

		//load template that corresponds for the device

		var details *DeviceSupportDetails

		//fix hardcoded framework, should be based on project
		for _, value := range desc.DeviceSupport {
			if value.Framework == "arduino" {
				details = value
				break
			}
		}

		if details != nil {
			tmplFile := path.Join(proj.GetContributionDir(), ref, details.TemplateFile)
			tmpl, err := fgutil.LoadLocalFile(tmplFile)
			if err != nil {
				return nil, err
			}

			contrib = &TriggerContrib{Descriptor: desc, Template: tmpl}
			triggerContribs[ref] = contrib
		}
	}

	return contrib, nil
}

func LoadActionContrib(proj Project, ref string) (*ActionContrib, error) {

	contrib, exists := actionContribs[ref]

	if !exists {
		//todo implement remote contribution for action

		return nil, errors.New("Action Not Supported")
	}

	return contrib, nil
}

func LoadActivityContrib(proj Project, ref string) (*ActivityContrib, error) {
	contrib, exists := activityContribs[ref]

	if !exists {
		proj.InstallContribution(ref, "")

		descFile := path.Join(proj.GetContributionDir(), ref, "activity.json")
		actJson, err := fgutil.LoadLocalFile(descFile)

		if err != nil {
			return nil, err
		}

		desc, err := ParseActivityDescriptor(actJson)
		//validate that device trigger

		//validate that supports device

		//load template that corresponds for the device

		var details *DeviceSupportDetails

		//fix hardcoded framework, should be based on project
		for _, value := range desc.DeviceSupport {
			if value.Framework == "arduino" {
				details = value
				break
			}
		}


		if details != nil {
			tmplFile := path.Join(proj.GetContributionDir(), ref, details.TemplateFile)
			tmpl, err := fgutil.LoadLocalFile(tmplFile)
			if err != nil {
				return nil, err
			}

			contrib = &ActivityContrib{Descriptor: desc, Template: tmpl}
			activityContribs[ref] = contrib
		}
		//else return error

	}

	return contrib, nil
}
