# Bug 是什么
context deadline 和取消没有沿 middleware/store/webhook 调用链传播。
# 如何触发
使用短超时或已取消 context 调用相关方法。
# 错误信息
webhook 在取消后仍发送请求，store 执行不感知取消。
