package queue

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

type Delayed interface {
	Comparable
	GetDelay() int64
	Run()
}

type DelayQueue struct {
	queue *PriorityQueue
	lock  sync.Mutex
	//queue-queue *time.Timer
	timerMap map[*time.Timer]*time.Timer
}

func NewDelayQueue() *DelayQueue {
	return &DelayQueue{queue: &PriorityQueue{}, timerMap: make(map[*time.Timer]*time.Timer)}
}

// 添加一个延迟任务
func (d *DelayQueue) Offer(ele Delayed) {
	log("DelayQueue Offer Exec")
	d.lock.Lock()
	log("DelayQueue offer lock")
	defer func() {
		d.lock.Unlock()
		log("DelayQueue offer unlock")
	}()
	heap.Push(d.queue, ele)
	d.signal()
}

// 获取一个延迟任务，没有任务或者任务没有到期则阻塞
func (d *DelayQueue) Take() Delayed {
	log("DelayQueue Take Exec")
	d.lock.Lock()
	log("DelayQueue Take lock")
	defer func() {
		d.lock.Unlock()
		log("DelayQueue Take unlock")
	}()
	for {
		var first Delayed
		if d.queue.Len() > 0 {
			first = (*d.queue)[0].(Delayed)
		}
		if first == nil {
			d.wait()
		} else if first.GetDelay() <= 0 {
			return heap.Pop(d.queue).(Delayed)
		} else {
			d.waitWithTimeout(first.GetDelay())
		}
	}
}

// 获取一个延迟任务，没有任务则阻塞timeout时长的时间，如果超时则返回nil
// 有任务并且已经到时间了则立即返回一个需要执行的任务
// timeout 是以秒为单位的时间
func (d *DelayQueue) TakeWithTimeout(timeout int64) Delayed {
	log("DelayQueue TakeWithTimeout Exec")
	d.lock.Lock()
	log("DelayQueue TakeWithTimeout lock")
	defer func() {
		d.lock.Unlock()
		log("DelayQueue TakeWithTimeout unlock")
	}()
	for {
		var first Delayed
		if d.queue.Len() > 0 {
			first = (*d.queue)[0].(Delayed)
		}
		if first == nil {
			if timeout <= 0 {
				return nil
			} else {
				timeout = d.waitWithTimeout(timeout)
			}
		} else {
			delay := first.GetDelay()
			if delay <= 0 {
				return heap.Pop(d.queue).(Delayed)
			} else if timeout <= 0 {
				return nil
			}
			first = nil
			if timeout < delay {
				timeout = d.waitWithTimeout(timeout)
			} else {
				timeout = d.waitWithTimeout(delay)
			}
			log("DelayQueue TakeWithTimeout 剩余超时时间timeout:", timeout)
		}
	}
}

// 无超时等待
func (d *DelayQueue) wait() {
	// 因为没有无限等待，所以这里设置成4字节最大的正数，表示无限等待
	d.waitWithTimeout(0x7fffffff)
}

// 超时等待
func (d *DelayQueue) waitWithTimeout(timeout int64) int64 {
	log("DelayQueue waitWithTimeout:", timeout)
	timer := time.NewTimer(time.Duration(timeout) * time.Second)
	d.timerMap[timer] = timer
	log("DelayQueue waitWithTimeout timerMap len:", len(d.timerMap))
	d.lock.Unlock()
	absTimeout := time.Now().Unix() + timeout
	select {
	case <-timer.C:
		// 为了剩余时间计算更准确不使用defer 获取锁
		d.lock.Lock()
		delete(d.timerMap, timer)
		return absTimeout - time.Now().Unix()
	}
}

// 只遍历出一个timer 直接唤醒
func (d *DelayQueue) signal() {
	for _, val := range d.timerMap {
		val.Reset(0)
		break
	}
}

// 唤醒所有timer
func (d *DelayQueue) signalAll() {
	for _, val := range d.timerMap {
		val.Reset(0)
	}
}

func (d *DelayQueue) Len() int {
	return d.queue.Len()
}

func log(msg ...interface{}) {
	if false {
		fmt.Println(msg)
	}
}
