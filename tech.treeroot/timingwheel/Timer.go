package timingwheel

import (
	"delay-queue/tech.treeroot/queue"
	"sync"
	"time"
)

type Timer struct {
	timingWheel *TimingWheel
	queue       *queue.DelayQueue
	taskCounter *AtomicInteger
	rwLock      sync.RWMutex
}

func NewTimer(tickMs int64, wheelSize int32, pollInterval int64) *Timer {
	queue := queue.NewDelayQueue()
	taskCounter := NewAtomicInteger()
	timingWheel := newTimingWheel(tickMs, wheelSize, time.Now().Unix(), taskCounter, queue)
	timer := &Timer{timingWheel: timingWheel, queue: queue, taskCounter: taskCounter}
	timer.init(pollInterval)
	return timer
}

func (t *Timer) init(pollInterval int64) {
	go func(int64) {
		for {
			t.advanceClock(1000)
		}
	}(pollInterval)
}

func (t *Timer) Add(task *TaskEntry) {
	t.rwLock.RLock()
	Log("Timer Add RLock")
	defer func() {
		t.rwLock.RUnlock()
		Log("Timer Add RUnlock")
	}()
	t.addTimerTaskEntry(task)

}

func (t *Timer) addTimerTaskEntry(entry *TaskEntry) {
	if !t.timingWheel.add(entry) {
		if !entry.isCancel {
			go func() {
				entry.run()
			}()
		}
	}
}

func (t *Timer) advanceClock(timeout int64) {
	temp := t.queue.TakeWithTimeout(timeout)
	var bucket *TaskList
	if temp != nil {
		bucket = temp.(*TaskList)
		t.rwLock.Lock()
		Log("Timer advanceClock Lock")
		t.timingWheel.advanceClock(bucket.expiration.Get())
		bucket.flush(t.addTimerTaskEntry)
		t.rwLock.Unlock()
		Log("Timer advanceClock Unlock")
	}
}

func (t *Timer) size() int32 {
	return int32(t.taskCounter.Get())
}

func (t *Timer) shutdown() {

}
