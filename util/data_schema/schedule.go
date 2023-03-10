package data_schema

import (
	"github.com/robfig/cron"
	"sync"
	"time"
)

// Schedule 调度整体信息
type Schedule struct {
	Entries Entries
	Create  chan TaskInfo
	Update  chan TaskInfo
	Delete  chan TaskInfo
	Init    chan int
	Lock    sync.Mutex
}

// Entry 调度任务
type Entry struct {
	Task     TaskInfo
	Schedule cron.Schedule
	Next     time.Time
	Prev     time.Time
}

type Entries []Entry

func (e Entries) Len() int {
	return len(e)
}

func (e Entries) Less(i, j int) bool {
	return e[i].Next.Before(e[j].Next)
}

func (e Entries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
