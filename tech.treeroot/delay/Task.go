package delay

import "time"

type Task struct {
	delayTime int64
	runnable func()
}

func (p *Task) GetCompareValue() int64 {
	return int64(p.delayTime)
}

func (p *Task) GetDelay() int64 {
	return int64(p.delayTime) - time.Now().Unix()
}

func (p *Task) Run() {
	(p.runnable)()
}

func NewDelayTask(delayTime int64, runnable func()) *Task {
	delayTime = time.Now().Unix() + delayTime
	return &Task{delayTime:delayTime, runnable:runnable}
}