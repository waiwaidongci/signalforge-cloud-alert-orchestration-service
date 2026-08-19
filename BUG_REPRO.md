# Bug 是什么
通知 Dispatcher 对未知通道和缺失 destination 返回的错误没有使用 sentinel 包装，调用方无法稳定识别错误类型。
# 如何触发
调用 `Dispatcher.Send` 传入未注册的 channel；或调用 `Service.Notify` 使用未注册 channel。
# 错误信息
`unknown notification channel "missing"` 或 `notification destination is required`，但 `errors.Is` 无法识别 sentinel。
