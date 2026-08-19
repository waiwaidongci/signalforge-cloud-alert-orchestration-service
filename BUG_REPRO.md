# Bug 是什么
规范化告警输入时，若 labels/annotations 或 source.FieldMapping 为零值，会向 nil map 写入并 panic。
# 如何触发
传入不带 labels 的告警，或使用 FieldMapping 为 nil 的 Source。
# 错误信息
`panic: assignment to entry in nil map`
