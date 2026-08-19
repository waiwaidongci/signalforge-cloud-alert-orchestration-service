# Bug 是什么
调度 worker 并发执行时，多个 goroutine 无锁写共享结果 map。
# 如何触发
并发调用 `runWorkers`，并用 race detector 执行。
# 错误信息
`WARNING: DATA RACE`
