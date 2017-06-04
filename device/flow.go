package device

import (
	"math"
	"strings"
)

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