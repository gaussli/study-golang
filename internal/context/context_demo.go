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

// 说明：返回一个空的、不会被取消、没有截止时间的 Context, 通常用于整个请求处理的最顶层，作为其他 Context 的基础。这个
// Context 不包含任何值，也不能被取消。
// 场景：在构建一个 Web 服务器应用时，在处理 HTTP 请求的最开始，可以使用 Context.Background() 作为整个请求处理流程的起点，
// 后续基于它创建带有请求特定信息（如请求ID、用户信息等）或带有取消、超时等功能的 Context。
func createBackground() {
	ctx := context.Background()
	fmt.Println(ctx)
	// 可以将这个基础 Context 传递给其他函数
	doSomething(ctx)
}

// 说明：返回一个空的、不会被取消、没有截止时间的 Context, 用于内部库中不确定应该使用什么 Context 时，或者在还没有确定具体的
// Context 使用方式时暂时使用这个 Context。它与 Context.Background() 类似，但更倾向于用于尚未确定具体使用场景的情况。
// 场景：在开发一个通用的工具库时，如果库的某些功能暂时不确定应该使用什么样的 Context （如是基于用户的请求 Context 还是全局的后台任务
// Context），可以先使用 Context.TODO()，后续根据实际业务需求再进行调整。
func createTODO() {
	ctx := context.TODO()
	fmt.Println(ctx)
	// 暂时不清楚如何使用这个 Context 时可以先用它
	doSomething(ctx)
}

func doSomething(ctx context.Context) {
	// 在函数内部可以基于这个基础 Context 创建带有新特性的 Context
	fmt.Println(ctx, "in doSomething method")
}

// 说明：这个函数创建一个带有取消功能的 Context。当调用返回的 cancel 函数时，它会取消这个新的 Context，并且通知所有监听这个
// Context 的 goroutine。
// 场景：在一个长时间运行的任务（如文件处理任务、网络请求任务等）中，当用户要求取消任务或者任务需要提前终止时，可以使用
// Context.WithCancel()。例如，在一个文件下载应用中，如果用户点击取消下载按钮，就可以调用 cancel 函数来取消下载任务。
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

// 说明：它创建一个带有超时时间的 Context。如果在这个超时时间之前任务没有完成，Context 会被自动取消。
// 场景：在调用外部服务（如 API 调用）时，为了避免长时间等待导致系统资源被占用，可以使用 Context.WithTimeout()。例如，在一个电商
// 应用中，当向支付网关发送支付请求时，可以设置一个超时时间，如果支付网关在超时时间内没有返回结果，就取消支付请求。
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

// 说明：它创建一个带有截止时间的 Context。如果当前时间超过了这个截止时间，Context 会被自动取消。
// 场景：在一个任务调度系统中，有些任务必须在特定的时间点之前完成。比如，在一个实时数据处理系统中，处理某个时间段的数据必须在这个时间段结束后的一定时间内完
// 成，否则数据就会失去价值了，这时可以使用 Context.WithDeadline() 来确保任务在规定的时间内完成。
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

// 说明：用于向 Context 中添加键值对数据。这可以在整个请求链中传递一些与请求相关的上下文信息，比如用户身份信息、请求 ID 等。
// 场景：在一个分布式系统中，当处理一个用户请求时，可以在请求进入系统的入口处（如 API 网关）使用 Context.WithValue() 向
// Context 添加用户身份信息，然后这个 Context 会随着请求在系统内部的各个服务之间传递，后续的服务就可以从 Context 中获取用户身份信
// 息来执行相应的操作，如权限验证等。
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
