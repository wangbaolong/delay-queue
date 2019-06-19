package timerwheel

import (
	"sync"
	"time"
)

type TaskList struct {
	expiration  *AtomicInteger // 原子变量
	taskCounter *AtomicInteger // 原子变量
	root        *TaskEntry
	lock        sync.Mutex
}

func newTaskList(taskCounter *AtomicInteger) *TaskList {
	root := NewTaskEntry(nil, 0)
	root.next = root
	root.prev = root
	return &TaskList{taskCounter: taskCounter, root: root, expiration: NewAtomicInteger()}
}

func (tl *TaskList) add(entry *TaskEntry) {
	tl.lock.Lock()
	defer tl.lock.Unlock()
	tail := tl.root.prev
	entry.next = tl.root
	entry.prev = tail
	tail.next = entry
	tl.root.prev = entry
	tl.taskCounter.IncrementAndGet()
}

func (tl *TaskList) remove(entry *TaskEntry) {
	tl.lock.Lock()
	defer tl.lock.Unlock()
	entry.next.prev = entry.prev
	entry.prev.next = entry.next
	entry.next = nil
	entry.prev = nil
	tl.taskCounter.DecrementAndGet()
}

func (tl *TaskList) flush(run func(entry *TaskEntry)) {
	head := tl.root.next
	for head != tl.root {
		tl.remove(head)
		run(head)
		head = tl.root.next
	}
}

func (tl *TaskList) setExpiration(expiration int64) bool {
	return tl.expiration.GetAndSet(expiration) != expiration
}



// Delayed 接口方法
func (tl *TaskList) GetDelay() int64 {
	return tl.expiration.Get() - time.Now().Unix()
}

// Delayed 接口方法
func (tl *TaskList) GetCompareValue() int64 {
	return tl.expiration.Get()
}

// Delayed 接口方法
func (tl *TaskList) Run() {

}
