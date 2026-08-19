# Bug 是什么
事件 fan-out 的 WaitGroup 计数在 goroutine 内，错误通道未关闭，成功路径会永久阻塞。
# 如何触发
调用 `FanOut` 执行全部成功的事件追加。
# 错误信息
测试等待错误通道超时。
