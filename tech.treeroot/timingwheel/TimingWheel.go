package timingwheel

import (
	"delay-queue/tech.treeroot/queue"
	"fmt"
	"sync"
	"sync/atomic"
)

type TimingWheel struct {
	tickMs            int64
	wheelSize         int32
	interval          int64
	startMs           int64
	taskCounter       *AtomicInteger
	queue             *queue.DelayQueue
	currentTime       int64
	buckets           []*TaskList
	overflowWheel     *TimingWheel
	overflowWheelLock sync.Mutex
}

func newTimingWheel(tickMs int64, wheelSize int32, startMs int64, taskCounter *AtomicInteger, queue *queue.DelayQueue) *TimingWheel {
	interval := tickMs * int64(wheelSize)
	currentTime := startMs - (startMs % tickMs)
	var buckets []*TaskList
	for i := 0; i < int(wheelSize); i++ {
		buckets = append(buckets, newTaskList(taskCounter))
	}
	return &TimingWheel{tickMs: tickMs, wheelSize: wheelSize, interval: interval, startMs: startMs, taskCounter: taskCounter, queue: queue, currentTime: currentTime, buckets: buckets}
}

func (tw *TimingWheel) addOverflowWheel() {
	tw.overflowWheelLock.Lock()
	defer tw.overflowWheelLock.Unlock()
	if tw.overflowWheel == nil {
		tw.overflowWheel = newTimingWheel(tw.interval, tw.wheelSize, tw.startMs, tw.taskCounter, tw.queue)
	}
}

func (tw *TimingWheel) add(entry *TaskEntry) bool {
	eMs := entry.expirationMs
	if entry.Cancelled() {
		return false
	} else if eMs < tw.currentTime+tw.tickMs {
		return false
	} else if eMs < tw.currentTime+tw.interval {
		virtualId := eMs / tw.tickMs
		index := virtualId % int64(tw.wheelSize)
		bucket := tw.buckets[index]
		bucket.add(entry)
		if bucket.setExpiration(virtualId * tw.tickMs) {
			tw.queue.Offer(bucket)
		}
		return true
	} else {
		if tw.overflowWheel == nil {
			tw.addOverflowWheel()
		}
		tw.overflowWheel.add(entry)
		return true
	}
}

func (tw *TimingWheel) advanceClock(timeMs int64) {
	if timeMs >= tw.currentTime+tw.tickMs {
		tw.currentTime = timeMs - (timeMs % tw.tickMs)
		if tw.overflowWheel != nil {
			tw.overflowWheel.advanceClock(tw.currentTime)
		}
	}
}

func (tw *TimingWheel) getOverflowWheel(overflowWheel *TimingWheel) *TimingWheel {
	tw.overflowWheelLock.Lock()
	defer tw.overflowWheelLock.Unlock()
	return tw.overflowWheel
}

type AtomicInteger struct {
	num int64
}

func NewAtomicInteger() *AtomicInteger {
	return &AtomicInteger{0}
}

func (t *AtomicInteger) IncrementAndGet() int64 {
	return atomic.AddInt64(&t.num, 1)
}

func (t *AtomicInteger) DecrementAndGet() int64 {
	return atomic.AddInt64(&t.num, -1)
}

func (t *AtomicInteger) GetAndSet(num int64) int64 {
	return atomic.SwapInt64(&t.num, num)
}

func (t *AtomicInteger) Get() int64 {
	return t.num
}

func Log(msg ...interface{}) {
	if false {
		fmt.Println(msg)
	}
}
