package timetask

import (
	"fmt"
)

// TaskFactory ...
type TaskFactory interface {
	TaskImmediately(interface{}) error
	TaskGenerate(interface{}) error
	TaskSetting(interface{}) error
	TaskExcute(interface{}) error
	ModifyTimeTask(interface{})
	DeleteTimeTask(interface{})
	TasksHooked([]interface{})
}

// TaskTypeMap ...
var TaskTypeMap = make(map[string]TaskFactory)

// Register ...
func Register(name string, task TaskFactory) {
	if task == nil {
		panic("task: Register failed, taskType is nil")
	}
	if _, ok := TaskTypeMap[name]; ok {
		panic("auth: authType has called twice" + name)
	}
	TaskTypeMap[name] = task
}

// GetTaskIns ...
func GetTaskIns(name string) (task TaskFactory, err error) {
	if name == "" {
		panic("task: taskType is null " + name)
	}
	if task, ok := TaskTypeMap[name]; ok {
		return task, nil
	}
	return nil, fmt.Errorf("GetTaskIns failed, type is %s ", name)
}
