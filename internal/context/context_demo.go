package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	createBackground()

	createTODO()

	createWithCancel()

	createWithTimeout()

	createWithDeadline()

	createWithValue()
}

func createBackground() {
	ctx := context.Background()
	fmt.Println(ctx)
	// 可以将这个基础 Context 传递给其他函数
	doSomething(ctx)
}

func createTODO() {
	ctx := context.TODO()
	fmt.Println(ctx)
	// 暂时不清楚如何使用这个 Context 时可以先用它
	doSomething(ctx)
}

func doSomething(ctx context.Context) {
	// 在函数内部可以基于这个基础 Context 创建带有新特性的 Context
	fmt.Println(ctx)
}

func createWithCancel() {
	ctx, cancel := context.WithCancel(context.Background())
	// 启动一个 goroutine 执行任务
	go doLongTaskWithCancel(ctx)
	// 等待一段时间后取消任务
	time.Sleep(5 * time.Second)
	cancel()
	// 确保 cancel 通知已经生效
	time.Sleep(2 * time.Second)
}

func doLongTaskWithCancel(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[WithCancel] Task is cancelled")
			return
		default:
			fmt.Println("[WithCancel] Task is running")
			time.Sleep(1 * time.Second)
		}
	}
}

func createWithTimeout() {
	// 设置超时时间为 2 秒
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // 确保释放资源
	go doLongTaskWithTimeout(ctx)
	// 等待一段时间
	time.Sleep(3 * time.Second)
}

func doLongTaskWithTimeout(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[WithTimeout] Task is timed out")
			return
		default:
			fmt.Println("[WithTimeout] Task is running")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func createWithDeadline() {
	deadline := time.Now().Add(2 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel() // 确保释放资源
	go doLongTaskWithDeadline(ctx)
	// 等待一段时间
	time.Sleep(3 * time.Second)
}

func doLongTaskWithDeadline(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[WithDeadline] Task is beyond deadline")
			return
		default:
			fmt.Println("[WithDeadline] Task is running")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func createWithValue() {
	// 创建一个带有值的 Context，添加用户 ID 信息
	ctx := context.WithValue(context.Background(), "userID", "12345")
	// 在函数中使用这个 Context
	doSomethingWithValue(ctx)
}

func doSomethingWithValue(ctx context.Context) {
	// 从 Context 中获取值
	if userID, ok := ctx.Value("userID").(string); ok {
		fmt.Println("[WithValue] User ID:", userID)
	} else {
		fmt.Println("[WithValue] No User ID found")
	}
}
