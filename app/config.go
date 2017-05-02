package app

import (
	"encoding/json"
)

type ContribType int

const (
	ACTION     ContribType = 1 + iota
	TRIGGER
	ACTIVITY
	FLOW_MODEL
)

var ctStr = [...]string{
	"all",
	"action",
	"trigger",
	"activity",
	"flow-model",
}

func (m ContribType) String() string { return ctStr[m] }

func ToContribType(name string) ContribType {
	switch name {
	case "action":
		return ACTION
	case "trigger":
		return TRIGGER
	case "activity":
		return ACTIVITY
	case "flow-model":
		return FLOW_MODEL
	case "all":
		return 0
	}

	return -1
}

// FlogoAppDescriptor is the descriptor for a Flogo application
type FlogoAppDescriptor struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Version     string `json:"version"`
	Description string `json:"description"`

	Actions  []*ActionDescriptor `json:"actions"`
	Triggers []*TriggerDescriptor `json:"triggers"`
}

// TriggerDescriptor is the config descriptor for a Trigger
type TriggerDescriptor struct {
	ID  string `json:"id"`
	Ref string `json:"ref"`
}

// todo make make ActionDescriptor generic
// ActionDescriptor is the config descriptor for an Action
type ActionDescriptor struct {
	ID  string `json:"id"`
	Ref string `json:"ref"`
	Data *struct {
		Flow *struct {
			RootTask         *Task `json:"rootTask"`
			ErrorHandlerTask *Task `json:"errorHandlerTask"`
		} `json:"flow"`
	} `json:"data"`
}

// Task is part of the flow structure
type Task struct {
	Ref   string `json:"activityRef"`
	Tasks []*Task `json:"tasks"`
}

//FlogoPaletteDescriptor a package: just change to a list of references
type FlogoPaletteDescriptor struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`

	Extensions []Dependency `json:"extensions"`
}

type Descriptor struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type Dependency struct {
	ContribType ContribType
	Ref         string
}

func (d *Dependency) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ContribType string `json:"type"`
		Ref         string   `json:"ref"`
	}{
		ContribType: d.ContribType.String(),
		Ref:         d.Ref,
	})
}

func (d *Dependency) UnmarshalJSON(data []byte) error {
	ser := &struct {
		ContribType string `json:"type"`
		Ref         string   `json:"ref"`
	}{}

	if err := json.Unmarshal(data, ser); err != nil {
		return err
	}

	d.Ref = ser.Ref
	d.ContribType = ToContribType(ser.ContribType)

	return nil
}

type depHolder struct {
	deps []*Dependency
}

// ExtractDependencies extracts dependencies from from application descriptor
func ExtractDependencies(descriptor *FlogoAppDescriptor) []*Dependency {

	dh := &depHolder{}

	for _, action := range descriptor.Actions {
		dh.deps = append(dh.deps, &Dependency{ContribType:ACTION, Ref:action.Ref})

		if action.Data != nil && action.Data.Flow != nil {
			extractDepsFromTask(action.Data.Flow.RootTask, dh)
			//Error handle flow
			if action.Data.Flow.ErrorHandlerTask != nil {
				extractDepsFromTask(action.Data.Flow.ErrorHandlerTask, dh)
			}
		}
	}

	for _, trigger := range descriptor.Triggers {
		dh.deps = append(dh.deps,&Dependency{ContribType:TRIGGER, Ref:trigger.Ref})
	}

	return dh.deps
}

// extractDepsFromTask extract dependencies from a task and is children
func extractDepsFromTask(task *Task, dh *depHolder) {

	if task.Ref != "" {
		dh.deps = append(dh.deps, &Dependency{ContribType:ACTIVITY, Ref:task.Ref})
	}

	for _, childTask := range task.Tasks {
		extractDepsFromTask(childTask, dh)
	}
}
