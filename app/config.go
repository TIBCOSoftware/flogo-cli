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
			RootTask *Task `json:"rootTask"`
		}`json:"flow"`
	} `json:"data"`
}

// Task is part of the flow structure
type Task struct {
	Ref   string `json:"activityRef"`
	Tasks []*Task `json:"tasks"`
}

// FlogoPaletteDescriptor is the flogo palette descriptor object
type FlogoExtension struct {
	Type string `json:"type"`
	Ref  string `json:"ref"`
}

//FlogoPaletteDescriptor a package: just change to a list of references
type FlogoPaletteDescriptor struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`

	Extensions []FlogoExtension `json:"extensions"`
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

type refHolder struct {
	refs []string
}

// Extract references from from application descriptor
func ExtractRefs(descriptor *FlogoAppDescriptor) []string {

	rh := &refHolder{}

	for _, action := range descriptor.Actions {
		rh.refs = append(rh.refs, action.Ref)

		if action.Data != nil && action.Data.Flow != nil {
			extractRefsFromTask(action.Data.Flow.RootTask, rh)
		}
	}

	for _, trigger := range descriptor.Triggers {
		rh.refs = append(rh.refs, trigger.Ref)
	}

	return rh.refs
}

// extractRefsFromTask extract references from a task and is children
func extractRefsFromTask(task *Task, rh *refHolder) {

	if task.Ref != "" {
		rh.refs = append(rh.refs, task.Ref)
	}

	for _, childTask := range task.Tasks {
		extractRefsFromTask(childTask, rh)
	}
}
