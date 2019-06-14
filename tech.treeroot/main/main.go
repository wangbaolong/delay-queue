package main

import (
	"delay-queue/tech.treeroot/delay"
	"fmt"
	"time"
)

func main() {
	done := make(chan bool)
	delayQueue()
	<- done
}

func delayQueue() {
	delayQueue := delay.New()
	go func() {
		for {
			delayQueue.Offer(delay.NewDelayTask(5, func() {
				fmt.Println("5s 任务执行了")
			}))
			delayQueue.Offer(delay.NewDelayTask(10, func() {
				fmt.Println("10s 任务执行了")
			}))
			delayQueue.Offer(delay.NewDelayTask(20, func() {
				fmt.Println("20s 任务执行了")
			}))
			delayQueue.Offer(delay.NewDelayTask(30, func() {
				fmt.Println("30s 任务执行了")
			}))
			time.Sleep(1 * time.Second)
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
