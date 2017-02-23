package app

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

// Extract references from from application descriptor
func ExtractRefs(descriptor *FlogoAppDescriptor) []string {

	var refs []string

	for _, action := range descriptor.Actions {
		refs = append(refs, action.Ref)

		if action.Data != nil && action.Data.Flow != nil {
			extractRefsFromTask(action.Data.Flow.RootTask, refs)
		}
	}

	for _, trigger := range descriptor.Triggers {
		refs = append(refs, trigger.Ref)
	}

	return refs
}

// extractRefsFromTask extract references from a task and is children
func extractRefsFromTask(task *Task, refs []string) {

	refs = append(refs, task.Ref)

	for _, childTask := range task.Tasks {
		extractRefsFromTask(childTask, refs)
	}
}
