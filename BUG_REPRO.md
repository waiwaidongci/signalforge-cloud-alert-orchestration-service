# Bug 是什么
告警确认动作把状态写为 resolved，确认和关闭语义未区分。
# 如何触发
对 firing 告警执行 acknowledge。
# 错误信息
状态直接变为 resolved，而不是 acknowledged。
