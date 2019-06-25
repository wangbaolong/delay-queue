package main

import (
	"delay-queue/tech.treeroot/queue"
	"delay-queue/tech.treeroot/timingwheel"
	"fmt"
	"time"
)

type Demo struct {
	name string
	next *Demo
}

func main() {
	done := make(chan bool)
	//delayQueue()
	//demo := Demo{}
	//fmt.Println(demo.Count.IncrementAndGet())
	//fmt.Println(demo)
	timerTest()
	<- done
}

func timerTest() {
	timerWheel := timingwheel.NewTimer(1, 12, 10)
	go func() {
		count := 0
		for {
			timerWheel.Add(timingwheel.NewTaskEntry(func() {
				fmt.Println("10 延迟任务执行")
			}, 10))
			timerWheel.Add(timingwheel.NewTaskEntry(func() {
				fmt.Println("15 延迟任务执行")
			}, 15))
			timerWheel.Add(timingwheel.NewTaskEntry(func() {
				fmt.Println("20 延迟任务执行")
			}, 20))
			timerWheel.Add(timingwheel.NewTaskEntry(func() {
				fmt.Println("30 延迟任务执行")
			}, 30))
			timerWheel.Add(timingwheel.NewTaskEntry(func() {
				fmt.Println("50 延迟任务执行")
			}, 50))
			time.Sleep(1 * time.Second)
			count++
			if count >= 10 {
				break
			}
		}
	}()
}

func delayQueue() {
	delayQueue := queue.NewDelayQueue()
	go func() {
		count := 0
		for {
			delayQueue.Offer(queue.NewDelayTask(5, func() {
				fmt.Println("5s 任务执行了")
			}))
			delayQueue.Offer(queue.NewDelayTask(10, func() {
				fmt.Println("10s 任务执行了")
			}))
			delayQueue.Offer(queue.NewDelayTask(20, func() {
				fmt.Println("20s 任务执行了")
			}))
			delayQueue.Offer(queue.NewDelayTask(30, func() {
				fmt.Println("30s 任务执行了")
			}))
			time.Sleep(10 * time.Millisecond)
			count++
			if count >= 10 {
				break
			}
		}
	}()

	// 启动 两个协程读取延迟队列的任务
	go func() {
		for {
			fmt.Println("take delayQueue len:", delayQueue.Len())
			// 没有任务则10s超时返回
			task := delayQueue.TakeWithTimeout(10)
			fmt.Println("take task:", task)
			if task != nil {
				go func() {
					task.Run()
				}()
			}
		}
	}()
	go func() {
		for {
			fmt.Println("take delayQueue len:", delayQueue.Len())
			// 没有任务则一直阻塞
			task := delayQueue.Take()
			fmt.Println("take task:", task)
			if task != nil {
				go func() {
					task.Run()
				}()
			}
		}
	}()

}
