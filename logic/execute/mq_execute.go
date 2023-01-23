package execute

import (
	"context"
	"schedule/util/data_schema"
	"schedule/util/log"
)

func MqExecute(ctx context.Context, task data_schema.TaskInfo) {
	log.InfoLogger.Printf("execute succ, task name:%s", task.TaskName)
}
