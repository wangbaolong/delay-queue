package timingwheel

import "time"

type TaskEntry struct {
	runnable     func()
	expirationMs int64
	isCancel     bool
	next         *TaskEntry
	prev         *TaskEntry
	createAt int64
}

func NewTaskEntry(run func(), delayTime int64) *TaskEntry {
	delayTime = time.Now().Unix() + delayTime
	return &TaskEntry{runnable: run, expirationMs: delayTime, createAt:time.Now().Unix()}
}

func (t *TaskEntry) Cancel() {
	t.isCancel = true
}

func (t *TaskEntry) Cancelled() bool {
	return t.isCancel
}

func (t *TaskEntry) run() {
	if t.runnable != nil {
		Log("Current Time:", time.Now().Unix(), " Delay Time:", t.expirationMs)
		t.runnable()
	}
}
