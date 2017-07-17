package device

import (
	"math"
	"strings"
)

func init() {

	//todo possibly convert actions to plugins
	RegisterActionContrib("github.com/TIBCOSoftware/flogo-contrib/device/action/activity", tplActionActivity)
	RegisterActionContrib("github.com/TIBCOSoftware/flogo-contrib/device/action/flow", tplActionDeviceFlow)
}


//todo hardcoded for now, should be generated from action-ref
type ActivityActionConfig struct {
	UseTriggerVal bool          `json:"useTriggerVal"`
	Activity      *ActivityConfig  `json:"activity"`
}

//todo hardcoded for now, should be generated from action-ref
type FlowActionConfig struct {
	Flow map[string]interface{}  `json:"flow"`
}

type ActivityConfig struct {
	Id         string                `json:"id"`
	Ref        string                `json:"ref"`
	Attributes map[string]string `json:"attributes"`
}

func (ac *ActivityConfig) GetSetting(key string) string {
	return ac.Attributes[key]
}

//todo move to contrib contributions
func RegisterActionContrib(ref string, tpl string) *ActionContrib {

	action := &ActionContrib{Ref: ref, Template: tpl}
	actionContribs[action.Ref] = action

	return action
}

var tplActionActivity = `

void a_{{.Id}}_init() {
 	{{ template "activity_init" .Data }}
}

void a_{{.Id}}(int value) {
 	{{ template "activity_code" .Data }}
}
`

var tplActionDeviceFlow = `

void a_{{.Id}}_init() {
	{{range $task := .AllTasks -}}
	ac_{{$task.FlowId}}_{{$task.Id}}_init();
	{{ end }}
}

void a_{{.Id}}(int value) {
 	{{ template "T" .FirstTask }}
}
`

var tplActionDeviceFlowEval = `{{ define "T" -}}
	{{if .Precondition }}
	if ({{.Precondition}}) {
	    ac_{{.FlowId}}_{{.Id}}(value);
	    {{range $task := .NextTasks}}{{with $task}}{{template "T" .}}{{end}}
	    {{end}}
	}
	{{- else -}}
	ac_{{.FlowId}}_{{.Id}}(value);
	{{range $task := .NextTasks}}{{with $task}}{{template "T" .}}{{end}}
	{{end}}{{end}}{{ end }}
`

///////////////////////////////////////////////////////////////////////////////

//todo should move to future "flow action" plugin

type FlowTree struct {
	Id        string
	FirstTask *Task
	AllTasks  []*Task
}

type Task struct {
	Id           int
	FlowId       string
	ActivityRef  string
	Attributes   map[string]string
	Precondition string
	NextTasks    []*Task
	isFirst      bool
}

func (t *Task) GetSetting(key string) string {
	return t.Attributes[key]
}

func toFlowTree(Id string, flow map[string]interface{}) *FlowTree {

	flowTree := &FlowTree{Id: Id}

	tasks := make(map[int]*Task)

	taskReps := flow["tasks"].([]interface{})
	for _, taskRep := range taskReps {
		taskData := taskRep.(map[string]interface{})
		task := &Task{isFirst: true, FlowId:Id}
		task.Id = toInt(taskData["id"])
		task.ActivityRef = taskData["activityRef"].(string)
		task.Attributes = toStringMap(taskData["attributes"].(map[string]interface{}))
		tasks[task.Id] = task

		flowTree.AllTasks = append(flowTree.AllTasks, task)

	}

	if val, ok := flow["links"]; ok {

		linkReps := val.([]interface{})
		for _, linkRep := range linkReps {
			linkData := linkRep.(map[string]interface{})
			fromId := toInt(linkData["from"])
			toId := toInt(linkData["to"])
			linkType := toInt(linkData["type"])
			condition := ""
			if linkType == 1 {
				condition = linkData["value"].(string)
				condition = strings.Replace(condition, "${value}","value",-1)
			}

			tasks[toId].isFirst = false
			tasks[toId].Precondition = condition
			tasks[fromId].NextTasks = append(tasks[fromId].NextTasks, tasks[toId])
		}
	}

	for _, task := range tasks {
		if task.isFirst {
			flowTree.FirstTask = task
			break
		}
	}

	return flowTree
}

func toStringMap(inMap map[string]interface{}) map[string]string {

	strMap := make(map[string]string)
	for key, value := range inMap {
		strMap[key] = value.(string)
	}

	return strMap
}

func toInt(val interface{}) int {
	return int(math.Ceil(val.(float64)))
}