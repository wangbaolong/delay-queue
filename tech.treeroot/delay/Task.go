package delay

import (
	"fmt"
	"time"
)

type Task struct {
	delayTime int64
	runnable  func()
}

func (p *Task) GetCompareValue() int64 {
	return int64(p.delayTime)
}

func (p *Task) GetDelay() int64 {
	return int64(p.delayTime) - time.Now().Unix()
}

func (p *Task) Run() {
	fmt.Println("exec time:", time.Now().Unix(), "delay time:", p.delayTime)
	(p.runnable)()
}

// delayTime 已秒为单位
func NewDelayTask(delayTime int64, runnable func()) *Task {
	delayTime = time.Now().Unix() + delayTime
	return &Task{delayTime: delayTime, runnable: runnable}
}
